#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd -- "${SCRIPT_DIR}/.." && pwd)"

BINARY_PATH="${REPO_ROOT}/build/flashgate-mcp"
STAMP="$$-$(date +%s%N)"
REQUEST_PATH="${REPO_ROOT}/build/smoke-jsonrpc-${STAMP}-request.jsonl"
RESPONSE_PATH="${REPO_ROOT}/build/smoke-jsonrpc-${STAMP}-response.jsonl"
MOVE_SOURCE_RELATIVE="build/smoke-move-${STAMP}-source.txt"
MOVE_TARGET_RELATIVE="build/smoke-move-${STAMP}-target.txt"
MOVE_SOURCE_PATH="${REPO_ROOT}/${MOVE_SOURCE_RELATIVE}"
MOVE_TARGET_PATH="${REPO_ROOT}/${MOVE_TARGET_RELATIVE}"

cleanup() {
  rm -f "${REQUEST_PATH}" "${RESPONSE_PATH}" "${MOVE_SOURCE_PATH}" "${MOVE_TARGET_PATH}"
}
trap cleanup EXIT

if [[ ! -x "${BINARY_PATH}" ]]; then
  echo "Binary not found or not executable: ${BINARY_PATH}" >&2
  echo "Run: go build -o build/flashgate-mcp ./cmd/server" >&2
  exit 1
fi

export MCP_ROOT="${REPO_ROOT}"

mkdir -p "${REPO_ROOT}/build"

cat > "${REQUEST_PATH}" <<'JSONRPC'
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"smoke-test","version":"dev"}}}
{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"list_directory","arguments":{}}}
{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"get_path_info","arguments":{"path":"README.md"}}}
{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"get_path_info","arguments":{"path":"smoke-missing-path"}}}
JSONRPC

if [[ "${MCP_READ_ONLY:-}" != "true" ]]; then
  printf '%s' 'move-smoke' > "${MOVE_SOURCE_PATH}"
  printf '%s\n' "{\"jsonrpc\":\"2.0\",\"id\":6,\"method\":\"tools/call\",\"params\":{\"name\":\"move_path\",\"arguments\":{\"source\":\"${MOVE_SOURCE_RELATIVE}\",\"target\":\"${MOVE_TARGET_RELATIVE}\"}}}" >> "${REQUEST_PATH}"
fi

"${BINARY_PATH}" < "${REQUEST_PATH}" > "${RESPONSE_PATH}"

python3 - "${RESPONSE_PATH}" "${MOVE_TARGET_PATH}" <<'PY'
import json
import os
import sys

response_path = sys.argv[1]
move_target_path = sys.argv[2]

with open(response_path, "r", encoding="utf-8") as handle:
    responses = [json.loads(line) for line in handle if line.strip()]

expected_response_count = 5 if os.environ.get("MCP_READ_ONLY") == "true" else 6
if len(responses) != expected_response_count:
    raise SystemExit(f"Expected {expected_response_count} JSON-RPC responses, got {len(responses)}. Response file: {response_path}")

initialize = responses[0]
tools_list = responses[1]
list_directory = responses[2]
existing_path_info = responses[3]
missing_path_info = responses[4]

if initialize.get("id") != 1:
    raise SystemExit(f"Expected initialize response id 1, got {initialize.get('id')}")

protocol_version = initialize.get("result", {}).get("protocolVersion")
if not protocol_version:
    raise SystemExit("Initialize response does not contain protocolVersion")

server_name = initialize.get("result", {}).get("serverInfo", {}).get("name")
if server_name != "flashgate":
    raise SystemExit(f"Expected serverInfo.name flashgate, got {server_name!r}")

if tools_list.get("id") != 2:
    raise SystemExit(f"Expected tools/list response id 2, got {tools_list.get('id')}")

tools = tools_list.get("result", {}).get("tools")
if not tools:
    raise SystemExit("tools/list response does not contain tools")

tool_names = [tool.get("name") for tool in tools]

expected_tools = [
    "list_directory",
    "read_file",
    "get_path_info",
]

if os.environ.get("MCP_READ_ONLY") != "true":
    expected_tools.extend([
        "write_file",
        "create_directory",
        "delete_path",
        "copy_path",
        "move_path",
    ])

for expected_tool in expected_tools:
    if expected_tool not in tool_names:
        raise SystemExit(
            f"Expected tool {expected_tool!r} was not listed. Actual tools: {', '.join(tool_names)}"
        )

for tool_name in tool_names:
    if tool_name not in expected_tools:
        raise SystemExit(
            f"Unexpected tool {tool_name!r} was listed. Expected tools: {', '.join(expected_tools)}"
        )

if list_directory.get("id") != 3 or "entries" not in list_directory.get("result", {}):
    raise SystemExit("list_directory did not return entries")
if existing_path_info.get("id") != 4 or existing_path_info.get("result", {}).get("exists") is not True:
    raise SystemExit("get_path_info did not report README.md as existing")
if missing_path_info.get("id") != 5 or missing_path_info.get("result") != {"path": "smoke-missing-path", "exists": False}:
    raise SystemExit("get_path_info did not report the missing path correctly")
if os.environ.get("MCP_READ_ONLY") != "true":
    move_result = responses[5]
    if move_result.get("id") != 6 or move_result.get("result", {}).get("moved") is not True:
        raise SystemExit("move_path did not return moved=true")
    if not os.path.isfile(move_target_path):
        raise SystemExit("move_path did not perform rename semantics")

print("JSON-RPC smoke test passed.")
print(f"Protocol version: {protocol_version}")
print(f"Tools: {', '.join(tool_names)}")
PY
