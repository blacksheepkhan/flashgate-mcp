#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd -- "${SCRIPT_DIR}/.." && pwd)"
BINARY_PATH="${REPO_ROOT}/build/flashgate-mcp"
STAMP="$$-$(date +%s%N)"
STDOUT_PATH="${REPO_ROOT}/build/startup-smoke-${STAMP}-stdout.txt"
STDERR_PATH="${REPO_ROOT}/build/startup-smoke-${STAMP}-stderr.txt"
ROOT_FILE="${REPO_ROOT}/build/startup-root-file-${STAMP}.txt"
MISSING_ROOT="${REPO_ROOT}/build/startup-missing-root-${STAMP}"

cleanup() {
  rm -f "${STDOUT_PATH}" "${STDERR_PATH}" "${ROOT_FILE}"
  rm -rf "${MISSING_ROOT}"
}
trap cleanup EXIT

if [[ ! -x "${BINARY_PATH}" ]]; then
  echo "Binary not found or not executable: ${BINARY_PATH}" >&2
  echo "Run: go build -o build/flashgate-mcp ./cmd/server" >&2
  exit 1
fi

mkdir -p "${REPO_ROOT}/build"
printf '%s' 'not a directory' > "${ROOT_FILE}"

assert_no_host_path_leak() {
  local value="$1"
  if grep -Eq '([A-Za-z]:[\\/]|\\\\[^\\]+|/(home|etc|private|tmp)/)' <<<"${value}"; then
    echo "Startup diagnostics leaked a host path: ${value}" >&2
    exit 1
  fi
}

run_case() {
  local name="$1"
  local expected_exit="$2"
  local expected_stderr="$3"
  shift 3

  : > "${STDOUT_PATH}"
  : > "${STDERR_PATH}"

  set +e
  (
    cd "${REPO_ROOT}"
    env -u MCP_ROOT -u MCP_ALLOW_CWD_ROOT -u MCP_READ_ONLY "$@" \
      "${BINARY_PATH}" < /dev/null > "${STDOUT_PATH}" 2> "${STDERR_PATH}"
  )
  local exit_code=$?
  set -e

  if [[ ${exit_code} -ne ${expected_exit} ]]; then
    echo "Case '${name}': expected exit ${expected_exit}, got ${exit_code}" >&2
    cat "${STDERR_PATH}" >&2
    exit 1
  fi
  if [[ -s "${STDOUT_PATH}" ]]; then
    echo "Case '${name}': expected empty stdout" >&2
    cat "${STDOUT_PATH}" >&2
    exit 1
  fi

  local stderr_value
  stderr_value="$(cat "${STDERR_PATH}")"
  if [[ "${stderr_value}" != "${expected_stderr}" ]]; then
    echo "Case '${name}': expected stderr '${expected_stderr}', got '${stderr_value}'" >&2
    exit 1
  fi
  assert_no_host_path_leak "${stderr_value}"
}

run_case "missing root" 3 "flashgate-mcp: startup failed (missing_root)"
run_case "empty root" 3 "flashgate-mcp: startup failed (invalid_root)" MCP_ROOT=
run_case "whitespace root" 3 "flashgate-mcp: startup failed (invalid_root)" "MCP_ROOT=   "
run_case "relative root" 3 "flashgate-mcp: startup failed (invalid_root)" MCP_ROOT=subdir
run_case "parent root" 3 "flashgate-mcp: startup failed (invalid_root)" MCP_ROOT=..
run_case "dot without opt-in" 3 "flashgate-mcp: startup failed (invalid_root)" MCP_ROOT=.
run_case "opt-in without root" 3 "flashgate-mcp: startup failed (missing_root)" MCP_ALLOW_CWD_ROOT=true
run_case "invalid development option" 3 "flashgate-mcp: startup failed (invalid_development_option)" "MCP_ROOT=${REPO_ROOT}" MCP_ALLOW_CWD_ROOT=TRUE
run_case "missing filesystem root" 3 "flashgate-mcp: startup failed (root_not_found)" "MCP_ROOT=${MISSING_ROOT}"
run_case "file root" 3 "flashgate-mcp: startup failed (root_not_directory)" "MCP_ROOT=${ROOT_FILE}"
run_case "invalid read-only profile" 3 "flashgate-mcp: startup failed (invalid_profile)" "MCP_ROOT=${REPO_ROOT}" MCP_READ_ONLY=invalid
run_case "valid absolute root" 0 "" "MCP_ROOT=${REPO_ROOT}" MCP_READ_ONLY=true
run_case "development CWD" 0 "flashgate-mcp: warning: development CWD root enabled" MCP_ROOT=. MCP_ALLOW_CWD_ROOT=true MCP_READ_ONLY=true

echo "Startup negative smoke test passed."
