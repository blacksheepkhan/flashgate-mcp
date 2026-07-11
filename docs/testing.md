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

The default smoke test validates `initialize`, the exact eight-tool `tools/list`, `list_directory`, `get_path_info` for existing and missing paths, and `move_path` rename behavior. The read-only variant verifies the exact three-tool profile and invokes all five write-capable names, requiring the same generic Invalid params response without filesystem changes. The negative smoke validates all five removed legacy names in addition to malformed JSON, unknown methods, invalid `tools/call` params, and notification no-response behavior.

The startup-negative smoke covers missing/empty/whitespace/relative roots, `.` with and without the development opt-in, invalid development/read-only values, missing and file roots, a valid absolute root, exit codes, empty stdout, safe stderr categories and cleanup.

GitHub Actions runs default, read-only, negative JSON-RPC, and startup-negative smoke variants on both `windows-latest` and `ubuntu-latest`. The smoke scripts create per-run artifacts under `build/` and clean them before exit. Script output is CI diagnostic output; server stdout remains reserved for redirected JSON-RPC protocol messages.

Limit and redaction behavior is primarily covered by Go unit tests. Additional limit-negative smoke coverage can be added later if it can be done without broad smoke-script refactoring.

Focused contract tests compare runtime tool definitions with `docs/mcp-tool-catalog.json` for name, title, description, and the complete input schema. Filesystem tests cover Missing normalization, truthful directory creation, same-path/SameFile protection, overwrite type combinations, lexical and effective self-subtree rejection, target-state revalidation, deterministic cross-volume no-fallback behavior, and platform error classification.

### Planned MCP Compatibility Testing

The implemented protocol remains MCP `2025-11-25`. Future protocol or extension support requires version-negotiation, extension-negotiation, client fallback, and compatibility tests before it is advertised. Future input and output schemas will be validated against JSON Schema 2020-12, and official MCP conformance tooling will be evaluated. These checks are planned in Sprint 3.45 and are not current implementation claims.

### Benchmarks

Benchmarks will be added for performance-sensitive operations such as:

- directory listing
- file reading
- file copying
- search

Benchmark command:

```bash
go test -bench=. ./...
```

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
