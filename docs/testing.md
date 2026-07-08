# Testing

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
.\scripts\smoke-jsonrpc-negative.ps1
```

The default smoke test validates `initialize` and `tools/list`. The negative smoke test validates malformed JSON, unknown methods, invalid `tools/call` params, and notification no-response behavior.

Limit and redaction behavior is primarily covered by Go unit tests. Additional limit-negative smoke coverage can be added later if it can be done without broad smoke-script refactoring.

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
