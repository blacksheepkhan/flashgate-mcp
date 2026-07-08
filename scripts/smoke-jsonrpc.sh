#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd -- "${SCRIPT_DIR}/.." && pwd)"

BINARY_PATH="${REPO_ROOT}/build/fileserver-mcp"
REQUEST_PATH="${REPO_ROOT}/build/smoke-jsonrpc-request.jsonl"
RESPONSE_PATH="${REPO_ROOT}/build/smoke-jsonrpc-response.jsonl"

if [[ ! -x "${BINARY_PATH}" ]]; then
  echo "Binary not found or not executable: ${BINARY_PATH}" >&2
  echo "Run: go build -o build/fileserver-mcp ./cmd/server" >&2
  exit 1
fi

export MCP_ROOT="${REPO_ROOT}"

mkdir -p "${REPO_ROOT}/build"

cat > "${REQUEST_PATH}" <<'JSONRPC'
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"smoke-test","version":"dev"}}}
{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}
JSONRPC

"${BINARY_PATH}" < "${REQUEST_PATH}" > "${RESPONSE_PATH}"

python3 - "${RESPONSE_PATH}" <<'PY'
import json
import os
import sys

response_path = sys.argv[1]

with open(response_path, "r", encoding="utf-8") as handle:
    responses = [json.loads(line) for line in handle if line.strip()]

if len(responses) != 2:
    raise SystemExit(f"Expected 2 JSON-RPC responses, got {len(responses)}. Response file: {response_path}")

initialize = responses[0]
tools_list = responses[1]

if initialize.get("id") != 1:
    raise SystemExit(f"Expected initialize response id 1, got {initialize.get('id')}")

protocol_version = initialize.get("result", {}).get("protocolVersion")
if not protocol_version:
    raise SystemExit("Initialize response does not contain protocolVersion")

if tools_list.get("id") != 2:
    raise SystemExit(f"Expected tools/list response id 2, got {tools_list.get('id')}")

tools = tools_list.get("result", {}).get("tools")
if not tools:
    raise SystemExit("tools/list response does not contain tools")

tool_names = [tool.get("name") for tool in tools]

expected_tools = [
    "list_files",
    "read_file",
    "stat_path",
    "exists_path",
]

if os.environ.get("MCP_READ_ONLY") != "true":
    expected_tools.extend([
        "write_file",
        "mkdir",
        "delete_path",
        "move_path",
        "copy_path",
        "rename_path",
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

print("JSON-RPC smoke test passed.")
print(f"Protocol version: {protocol_version}")
print(f"Tools: {', '.join(tool_names)}")
PY
