# Architecture

`fileserver-mcp` is a high-performance Model Context Protocol (MCP) server written in Go.

The project is structured as a layered application with strict separation between protocol handling, tool execution, security validation and filesystem access.

## Architectural Goals

The primary architectural goals are:

- predictable behavior
- low resource usage
- cross-platform support for Windows and Linux
- secure-by-default filesystem access
- clear testability
- long-term maintainability
- no dependency on external MCP protocol libraries

## High-Level Architecture

```text
MCP Client
  |
  | STDIO
  v
Transport
  |
  v
JSON-RPC / MCP Server
  |
  v
Router
  |
  v
Handlers
  |
  v
Tools
  |
  v
Filesystem Abstraction
  |
  v
PathGuard
  |
  v
Operating System Filesystem
```

## Layers

### `cmd/server`

Contains only the application entry point.

The `main` package must not contain application logic. It is responsible only for:

- loading configuration
- wiring dependencies
- starting the server
- reporting fatal startup errors

### `internal/config`

Contains immutable application configuration.

Configuration is loaded from defaults and environment variables. Future configuration sources may include configuration files and command-line flags.

### `internal/diagnostics`

Contains redaction and minimal stderr diagnostic helpers.

Diagnostics must never write to stdout and must not log raw request payloads, file contents, secrets, credentials, connection strings, or host paths.

### `internal/security`

Contains security-sensitive validation logic.

Currently this package provides:

- `PathGuard`
- `SafePath`
- sandbox root validation
- path traversal protection
- absolute path rejection

No package outside `internal/security` should implement its own path validation rules.

### `internal/fs`

Contains all direct filesystem access.

No production code outside this package should use filesystem operations such as:

- `os.ReadFile`
- `os.WriteFile`
- `os.ReadDir`
- `os.Remove`
- `os.Rename`
- `os.Open`

The filesystem package exposes a controlled `FileSystem` interface that is used by MCP tools.

### `internal/mcp`

Contains the MCP runtime components:

- transport
- router
- handlers
- tools
- server runtime

The MCP layer must not directly access the operating system filesystem.

### `internal/protocol`

Contains protocol data structures only.

This package should not contain business logic.

## Dependency Direction

Dependencies must point inward toward lower-level abstractions.

Allowed examples:

```text
mcp/tools -> fs
fs        -> security
server    -> router, transport, protocol
server    -> diagnostics
```

Disallowed examples:

```text
security -> fs
fs       -> mcp
protocol -> server
```

## Filesystem Core

The filesystem core currently supports:

- list
- read
- stat
- exists
- write
- mkdir
- delete
- move
- copy
- rename

All operations are sandboxed through `PathGuard`.

## Design Principles

### Single Responsibility

Each package and file should have one clear responsibility.

### Dependency Injection

Components should receive dependencies explicitly.

### Secure Defaults

Potentially destructive operations must be conservative by default.

Examples:

- file writes do not overwrite unless explicitly requested
- directory deletion is not recursive unless explicitly requested
- directory copy is intentionally unsupported in the initial implementation

### Small Files

Large implementation files should be split by responsibility.

The filesystem implementation is split into focused files:

```text
filesystem.go
list.go
read.go
stat.go
write.go
delete.go
move.go
copy.go
```

## Current Status

The foundation layer and MCP protocol boundary are implemented and tested. Further work should extend the existing layers without bypassing the configured limits, diagnostics redaction, or filesystem security boundary.
