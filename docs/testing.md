# Testing

FlashGate MCP is the public project name from Sprint 3.41. Commands and artifact paths in this document still use the technical `fileserver-mcp` identifier until Sprint 3.42.

`fileserver-mcp` uses Go's standard testing framework.

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
go build -o build/fileserver-mcp ./cmd/server
```

On Windows, the build command is usually:

```powershell
go build -o build/fileserver-mcp.exe ./cmd/server
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

The default smoke test validates `initialize` and `tools/list`. The read-only variant verifies that write-capable tools are not registered when `MCP_READ_ONLY=true`. The negative smoke test validates malformed JSON, unknown methods, invalid `tools/call` params, and notification no-response behavior.

GitHub Actions runs default, read-only, and negative JSON-RPC smoke variants on both `windows-latest` and `ubuntu-latest`. The smoke scripts create per-run JSONL request and response files under `build/` and clean them up before exit. Script output is CI diagnostic output; server stdout remains reserved for redirected JSON-RPC protocol messages.

Limit and redaction behavior is primarily covered by Go unit tests. Additional limit-negative smoke coverage can be added later if it can be done without broad smoke-script refactoring.

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
