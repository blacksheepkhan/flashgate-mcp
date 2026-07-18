# Testing

FlashGate MCP uses Go's standard testing framework. Sprint 3.42 updated commands and artifact paths to the `flashgate-mcp` binary.

The project aims for high test coverage in security-sensitive and filesystem-related code.

## Test Commands

Run all tests:

```bash
go test ./...
```

Run tests with the race detector:

```bash
go test -race ./...
```

The authoritative race gate runs on a platform with a supported race toolchain.
The current Windows host has no CGO/GCC race toolchain, so Windows race is reported
as an infrastructure limitation rather than worked around. Native Linux
`go test -race ./...` is required. Functional serialization, payload, fixture, and
budget-contract tests remain active under race; only the `testing.AllocsPerRun`
budget assertion is skipped because race instrumentation changes allocation
behavior. Ordinary non-race tests continue to enforce the unchanged allocation
budgets.

Run tests for a specific package:

```bash
go test ./internal/fs
```

Run focused MCP protocol tests:

```bash
go test -v ./internal/protocol ./internal/mcp/server ./internal/mcp/router ./internal/mcp/tools ./internal/mcp/initialize
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Required Quality Checks

Before committing, run:

```bash
go fmt ./...
go vet ./...
go test ./...
golangci-lint run
go build -o build/flashgate-mcp ./cmd/server
```

On Windows, the build command is usually:

```powershell
go build -o build/flashgate-mcp.exe ./cmd/server
```

## Test Strategy

### Unit Tests

Unit tests are required for:

- configuration loading
- path validation
- filesystem operations
- MCP protocol routing
- tool execution

### Filesystem Tests

Filesystem tests use:

```go
t.TempDir()
```

This ensures that tests do not modify real user data.

Test helper functions may use `os.WriteFile`, `os.ReadFile`, `os.MkdirAll` and `os.Stat` to create and verify test fixtures.

This is allowed in tests.

Production code outside `internal/fs` must not use direct filesystem operations.

### Security Tests

Security tests must cover:

- path traversal
- absolute path rejection
- sandbox escape attempts
- destructive operation defaults
- overwrite behavior
- recursive delete behavior
- JSON-RPC message and tool argument limits
- filesystem read, write, list, copy, and recursive delete limits
- diagnostics redaction

### Integration Tests

JSON-RPC smoke tests exercise the built server binary over STDIO.

On Windows:

```powershell
.\scripts\smoke-jsonrpc.ps1
$env:MCP_READ_ONLY = "true"
.\scripts\smoke-jsonrpc.ps1
Remove-Item Env:\MCP_READ_ONLY
.\scripts\smoke-jsonrpc-negative.ps1
```

On Linux:

```bash
bash scripts/smoke-jsonrpc.sh
MCP_READ_ONLY=true bash scripts/smoke-jsonrpc.sh
bash scripts/smoke-jsonrpc-negative.sh
```

Run fail-closed startup validation on Windows:

```powershell
.\scripts\smoke-startup-negative.ps1
```

On Linux:

```bash
bash scripts/smoke-startup-negative.sh
```

The default smoke test validates `initialize`, the exact eight-tool `tools/list`, `list_directory`, `read_file`, `get_path_info` for existing and missing paths, and `move_path` rename behavior. Every positive result must have only the `CallToolResult` envelope fields, exactly one text block with only `type` and `text`, valid compact object JSON, and a deeply equal `structuredContent` object. The read-only variant verifies the exact three-tool profile and invokes all five write-capable names, requiring the same generic Invalid params response without filesystem changes. The negative smoke validates all five removed legacy names in addition to malformed JSON, unknown methods, invalid `tools/call` params, and notification no-response behavior.

The startup-negative smoke covers missing/empty/whitespace/relative roots, `.` with and without the development opt-in, invalid development/read-only values, missing and file roots, a valid absolute root, exit codes, empty stdout, safe stderr categories and cleanup.

GitHub Actions runs default, read-only, negative JSON-RPC, and startup-negative smoke variants on both `windows-latest` and `ubuntu-latest`. The smoke scripts create per-run artifacts under `build/` and clean them before exit. Script output is CI diagnostic output; server stdout remains reserved for redirected JSON-RPC protocol messages.

Limit and redaction behavior is primarily covered by Go unit tests. Additional limit-negative smoke coverage can be added later if it can be done without broad smoke-script refactoring.

Focused contract tests compare runtime tool definitions with `docs/mcp-tool-catalog.json` for name, title, description, complete input schema, and deeply equal runtime `outputSchema`/catalog `resultSchema`. Targeted tests require exactly eight runtime output schemas, object roots, valid required/property relationships, expected project property types, representative successful `structuredContent`, both `get_path_info` variants, and the `read_file` outer-array/inner-string distinction. The tests-only structural checker covers only `type`, `properties`, `required`, `additionalProperties`, `items`, `oneOf`, and `const` as currently emitted; it is not a complete JSON Schema 2020-12 validator.

The `tools/list` JSON-RPC wire test checks schema exposure for both profiles and records deterministic UTF-8 JSONL sizes with and without output schemas. Sprint 3.45b records 1239/2134 bytes for read-only and 3850/5657 bytes for default; no regression budget is enforced.

### MCP Compatibility Testing

The implemented protocol remains MCP `2025-11-25`. Explicit `CallToolResult` DTO tests, a strict project-local decoder, legacy unwrapped negative fixtures, all-eight-tool adapter coverage, and full JSON-RPC wire tests cover success and the unchanged error contract. The decoder intentionally validates the exact FlashGate-emitted subset (one text block, required object `structuredContent`, optional boolean `isError`, no `_meta`) rather than claiming to decode every standard-conformant MCP result. Windows and Bash positive smokes enforce the same shape.

Future protocol or extension support still requires version-negotiation, extension-negotiation, client fallback, and compatibility tests before it is advertised. Complete JSON Schema 2020-12 validation and official MCP conformance tooling remain planned.

### Benchmarks

Sprint 3.45d benchmarks performance-sensitive operations including:

- directory listing
- file reading
- file copying
- tool-result wrapping and serialized payload forms
- search

Benchmark command:

```bash
go test -bench=. ./...
```

The deterministic benchmark tests also enforce the exact tool-profile and workflow measurement sets, the six payload/allocation budgets from `benchmarks/budgets.json`, initialized-notification framing, initialization-result validation, zero `scanned_bytes` for ordinary reads, partial Linux procfs metrics, host-path redaction, and clean versioned artifacts. Platform baseline generation is a separate two-phase operation after the implementation commit. The diagnostic wrappers recognize `-RecordBaseline` and `--record-baseline` only to reject them fail-closed before any Go invocation or write.

Functional gates such as tests, builds, vet, lint, protocol smokes, and parser checks are independent of the host measurement window. Their timing and resource consumption are not performance evidence. Performance gates are valid only when the entire measurement series runs outside the primary development host's scheduled-load block from 19:00 inclusive until 04:00 exclusive in `Europe/Vienna`; the preferred safety-margin window is 04:15–18:45, and the series must finish before 19:00. The blocked interval is a formal baseline blocker.

Every performance measurement report records the `Europe/Vienna` time window, start and end times, and whether known or unusual additional host load was present. A baseline is rejected if such load is known or observed even inside the nominally allowed interval. Contaminated runs are retained as diagnosis evidence without being approved, compared for regression, or used to tune budgets. Ordinary wrapper runs may continue but are marked contaminated; they cannot record a baseline.

## Current Tested Packages

Currently tested:

- `internal/config`
- `internal/diagnostics`
- `internal/security`
- `internal/fs`
- `internal/protocol`
- `internal/mcp/server`
- `internal/mcp/router`
- `internal/mcp/transport`
- `internal/mcp/initialize`
- `internal/mcp/tools`

## Version 1.0 Planned Validation Matrix

The current tests above describe the implemented filesystem baseline. Version 1.0 adds the following required gates.

### Payload and catalog tests

- payload-class selection for metadata, structured pages, heavy text, media/binary, and large results;
- heavy payload appears only once across MCP result fields;
- wire-amplification and useful-byte budgets;
- bounded base64 thresholds;
- opaque resource handles contain no host path and enforce owner/TTL/service-generation checks;
- fallback behavior for clients without resource-link support;
- deterministic tool ordering and catalog fingerprint;
- profile-specific `tools/list`, schema, description, and initialization-instruction budgets;
- safe read-only catalog when roots exist and no explicit profile is selected.

### Operations and multi-principal tests

- opaque handles bound to principal, profile, root, execution backend, and service generation;
- cross-principal status/result/cancel/cache/resource denial;
- global, per-domain, and per-principal concurrency limits;
- global/per-principal queue caps and fair scheduling;
- deterministic overload behavior;
- TTL cleanup, restart invalidation, shutdown, and leak detection;
- slow-reader and audit/log backpressure behavior.

### Typed command tests

- executable ID resolves only to approved absolute binary;
- fixed subcommand and allowed flags/value rules;
- path arguments remain under allowed roots;
- no shell interpretation;
- response files, hooks, plugins, loaders, config overrides, and unapproved environment are rejected;
- stdout/stderr, runtime, process count, and network policy limits;
- Windows/Linux isolation outcomes and redaction.

### System service and execution-identity tests

Version 1.0 tests Variant A only:

- Windows SCM and Linux systemd lifecycle;
- Named Pipe ACL and Unix socket ownership/mode;
- OS-derived peer identity cannot be overridden by payload;
- caller authorization independent of service-account filesystem permission;
- allowed FlashGate policy plus denied service-account ACL fails safely;
- denied caller plus available service-account ACL fails before execution;
- service-account root backend and dedicated identity;
- no LocalSystem/root convenience default;
- unsupported `user-worker` configuration fails closed;
- no in-process impersonation path;
- caller and effective backend identity both appear in bounded audit events;
- service restart invalidates generation-bound handles/resources;
- `auto` never falls back after managed denial or incompatibility;
- proxy/client stdout remains MCP-only.

Variant B worker tests are post-Version 1.0 and require a separate implementation gate.

### Protocol compatibility tests

Before Version 1.0, publish and test the supported MCP revision matrix:

- current `2025-11-25` behavior;
- any later final revision only after implementation;
- stateless-core behavior where selected;
- deterministic list cache/TTL invalidation;
- final Tasks Extension mapping without mixing the 2025 experimental lifecycle;
- extension downgrade/mismatch;
- JSON Schema 2020-12 validation;
- deprecated Roots never overrides server roots.

### Audit and failure-path tests

- immutable event/correlation IDs;
- proxy/service/backend/job/process correlation;
- redaction before output;
- rotation and retention;
- slow sink and bounded buffering;
- disk-full behavior;
- log-injection handling;
- shutdown flush/drop policy;
- no secret, full payload, unrestricted environment, or unnecessary host-path leakage.

### Release and supply-chain tests

- artifact version/help/platform/name checks;
- no interpreter runtime dependency;
- service asset syntax and install/remove dry validation;
- checksums;
- SBOM and dependency inventory;
- build provenance;
- signing verification where configured;
- reproducible-build comparison;
- pinned/validated workflow policy;
- atomic rollback documentation and smoke procedure.

### Cross-project benchmark

The Version 1.0 benchmark compares pinned FlashGate, official Node.js filesystem, selected native Rust filesystem, and selected Go filesystem MCP versions on the same host and corpus. The report must separate feature/security differences from measured performance and must not claim results for unmeasured operations.

See [Efficiency Improvement Plan](efficiency-improvement-plan.md), [Execution Identity Backends](execution-identity-backends.md), and [Version 1.0 Scope](version-1-scope-and-release-boundary.md).

<!-- FLASHGATE_PERFORMANCE_WORKSPACE_POLICY_START -->
## Authoritative benchmark workspace gate

On the primary Windows development host, an authoritative benchmark attempt is
blocked unless its Windows working area is below:

`C:\Voxtronic\Codex\Temp\Benchmarks`

Before the attempt, verify that the root is a fixed local NTFS path, contains no
reparse point, and is not below OneDrive, Dropbox, a redirected user directory,
a network share, or other synchronized storage.

All source bundles, isolated Windows checkouts, prepared binaries, measurement
outputs, logs, verification files, and controller data stay in that local area
through the final host-load gate. Archival copying to synchronized storage occurs
only afterward and is not part of the measured phase.

The native Linux checkout and temporary output remain on the distribution's
native ext4 filesystem under `/home`. A path below `/mnt` or `/media` is a formal
baseline blocker.

All validation, test, vet, lint, build, linker, and parser work finishes before
the authoritative host gate. After the last such operation, wait at least 180
seconds without Git, Go, scan, archive, or analysis activity. Run the authoritative
three-block CPU/disk/RAM/per-process-delta preflight exactly once. If it passes,
invoke the prepared binaries directly without rebuilding. A 15-second intermediate
gate precedes native Linux measurement and a final host gate precedes any result
copy, hash scan, JSON verification, report, archive, or OneDrive access.

`scripts/benchmark.ps1 -RecordBaseline` and
`scripts/benchmark.sh --record-baseline` are deliberately blocked compatibility
flags, not an authoritative workflow. A separately prepared controller implements
the two-phase attempt; no wrapper-side shortcut or time override is permitted.
<!-- FLASHGATE_PERFORMANCE_WORKSPACE_POLICY_END -->
