#!/usr/bin/env bash
set -uo pipefail

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
failures=()
check_count=0
temporary_directory=""

check() {
  local condition="$1"
  local name="$2"
  check_count=$((check_count + 1))
  if [[ "${condition}" != true ]]; then
    failures+=("${name}")
  fi
}

cleanup() {
  if [[ -n "${temporary_directory}" && -d "${temporary_directory}" ]]; then
    rm -rf -- "${temporary_directory}"
  fi
}
trap cleanup EXIT

temporary_directory="$(mktemp -d)"
fake_repository="${temporary_directory}/repo"
fake_scripts="${fake_repository}/scripts"
stub_directory="${temporary_directory}/stub"
mkdir -p -- "${fake_scripts}" "${stub_directory}"
wrapper="${fake_scripts}/benchmark.sh"
cp -- "${script_dir}/benchmark.sh" "${wrapper}"
marker="${temporary_directory}/go-invoked.marker"
output="${temporary_directory}/baseline.linux-amd64.json"
stub="${stub_directory}/go"
printf '#!/usr/bin/env bash\n: > %q\nexit 97\n' "${marker}" >"${stub}"
chmod 700 "${stub}"

stdout_path="${temporary_directory}/stdout.log"
stderr_path="${temporary_directory}/stderr.log"
PATH="${stub_directory}:${PATH}" bash "${wrapper}" --record-baseline --output "${output}" >"${stdout_path}" 2>"${stderr_path}"
exit_code=$?
stderr_value="$(<"${stderr_path}")"

check "$([[ ${exit_code} -ne 0 ]] && printf true || printf false)" 'record mode exits nonzero'
check "$([[ "${stderr_value}" == *'Authoritative baseline recording is not supported by scripts/benchmark.sh'* ]] && printf true || printf false)" 'record mode emits policy error on stderr'
check "$([[ ! -e "${marker}" ]] && printf true || printf false)" 'record mode does not invoke go'
check "$([[ ! -e "${output}" ]] && printf true || printf false)" 'record mode does not write baseline output'
check "$([[ ! -d "${fake_repository}/build" ]] && printf true || printf false)" 'record mode does not create build directory'
check "$([[ ! -s "${stdout_path}" ]] && printf true || printf false)" 'record mode does not emit success output'

canonical_output="${fake_repository}/benchmarks/../benchmarks/baseline.linux-amd64.json"
: >"${stdout_path}"
: >"${stderr_path}"
PATH="${stub_directory}:${PATH}" bash "${wrapper}" --output "${canonical_output}" >"${stdout_path}" 2>"${stderr_path}"
canonical_exit_code=$?
canonical_stdout="$(<"${stdout_path}")"
canonical_stderr="$(<"${stderr_path}")"
check "$([[ ${canonical_exit_code} -ne 0 ]] && printf true || printf false)" 'non-authoritative mode rejects canonical baseline path'
check "$([[ "${canonical_stderr}" == *'Non-authoritative benchmark runs must not write a canonical versioned baseline path.'* ]] && printf true || printf false)" 'canonical-path rejection is reported'
check "$([[ -z "${canonical_stdout}" ]] && printf true || printf false)" 'canonical-path rejection emits no success output'
check "$([[ ! -e "${marker}" ]] && printf true || printf false)" 'canonical-path rejection does not invoke go'
check "$([[ ! -e "${fake_repository}/benchmarks/baseline.linux-amd64.json" ]] && printf true || printf false)" 'canonical-path rejection does not write output'
check "$([[ ! -d "${fake_repository}/build" ]] && printf true || printf false)" 'canonical-path rejection does not create build directory'

if (( ${#failures[@]} == 0 )); then
  printf 'Status       : PASS\nCheckCount   : %d\nFailureCount : 0\nFailures     :\n' "${check_count}"
  exit 0
fi
printf 'Status       : FAIL\nCheckCount   : %d\nFailureCount : %d\nFailures     : %s\n' "${check_count}" "${#failures[@]}" "${failures[*]}"
exit 1
