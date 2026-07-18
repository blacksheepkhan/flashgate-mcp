#!/usr/bin/env bash
set -euo pipefail

quick=false
record_baseline=false
output_path=""
performance_contaminated=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --quick)
      quick=true
      ;;
    --record-baseline)
      record_baseline=true
      ;;
    --output)
      shift
      output_path="${1:?--output requires a path}"
      ;;
    *)
      printf 'Unknown argument: %s\n' "$1" >&2
      exit 2
      ;;
  esac
  shift
done

if [[ "${record_baseline}" == true ]]; then
  printf 'Authoritative baseline recording is not supported by scripts/benchmark.sh. Use the documented two-phase prebuilt workflow with a native Linux checkout under /home and the Windows controller workspace outside synchronized storage.\n' >&2
  exit 2
fi

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd -- "${script_dir}/.." && pwd)"
build_dir="${repo_root}/build"
server_binary="${build_dir}/flashgate-mcp"
benchmark_binary="${build_dir}/flashgate-benchmark"
budget_path="${repo_root}/benchmarks/budgets.json"

if [[ -n "${output_path}" ]]; then
  if [[ "${output_path}" != /* ]]; then
    output_path="${repo_root}/${output_path}"
  fi
  output_path="$(realpath -m -- "${output_path}")"
  if [[ "${output_path}" == "${repo_root}"/benchmarks/baseline.*-*.json ]]; then
    printf 'Non-authoritative benchmark runs must not write a canonical versioned baseline path.\n' >&2
    exit 1
  fi
fi

# shellcheck source=benchmark-window.sh
source "${script_dir}/benchmark-window.sh"

if measurement_window_is_blocked; then
  performance_contaminated=true
fi

working_tree_dirty="$(git status --porcelain --untracked-files=all)"

if [[ -z "${output_path}" ]]; then
  output_path="${build_dir}/benchmark-current.linux-$(go env GOARCH).json"
fi

mkdir -p "${build_dir}"
go build -o "${server_binary}" ./cmd/server
go build -o "${benchmark_binary}" ./cmd/benchmark
commit="$(git rev-parse HEAD)"

arguments=(
  -binary "${server_binary}"
  -output "${output_path}"
  -commit "${commit}"
  -budgets "${budget_path}"
)
if [[ -n "${working_tree_dirty}" ]]; then
  arguments+=(-working-tree-dirty)
fi
if [[ "${quick}" == true ]]; then
  arguments+=(-quick)
fi

"${benchmark_binary}" "${arguments[@]}"

if measurement_window_is_blocked; then
  performance_contaminated=true
fi
if [[ "${performance_contaminated}" == true ]]; then
  measurement_window_contaminated_warning >&2
fi
printf 'Benchmark result: %s\n' "${output_path}"
