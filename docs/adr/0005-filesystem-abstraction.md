# ADR-0005: Filesystem Abstraction

## Status

Accepted

## Context

`fileserver-mcp` exposes filesystem operations to MCP clients.

Direct filesystem access is security-sensitive. If filesystem APIs are used across multiple packages, security checks may become inconsistent or be bypassed accidentally.

The project must support:

- Windows
- Linux
- sandboxed filesystem access
- future read-only mode
- future audit logging
- future search functionality
- unit testing without touching real user files

## Decision

All production filesystem operations must go through a central `FileSystem` abstraction in:

```text
internal/fs
```

The main interface is:

```go
type FileSystem interface {
    List(path string) ([]Entry, error)
    Read(path string, maxBytes int64) ([]byte, error)
    Stat(path string) (Metadata, error)
    Exists(path string) (bool, error)
    Write(path string, content []byte, overwrite bool) error
    Mkdir(path string) error
    Delete(path string, recursive bool) error
    Move(source string, target string, overwrite bool) error
    Copy(source string, target string, overwrite bool) error
    Rename(source string, target string, overwrite bool) error
}
```

`internal/fs` is the only package that may directly use operating system filesystem calls in production code.

All paths are validated through `internal/security.PathGuard`.

## Consequences

### Positive

- security checks are centralized
- MCP tools remain simple
- filesystem behavior is easier to test
- future implementations can provide an in-memory filesystem or mock filesystem
- audit logging can be added in one place
- read-only mode can be enforced centrally

### Negative

- initial implementation requires more structure
- simple tools must depend on an abstraction rather than calling `os.*` directly
- some operations require careful policy decisions, such as overwrite and recursive delete behavior

## Security Impact

This decision significantly reduces the risk of accidental sandbox bypasses.

All MCP tools must use the `FileSystem` interface instead of direct filesystem access.

## Current Implementation

The current local filesystem implementation is:

```text
internal/fs/LocalFileSystem
```

It supports:

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

Directory copy is intentionally unsupported in the initial implementation.

The current `copy_path` tool copies files only. Directory copy remains unsupported.

## Amendment - 2026-07-10

The public project name is FlashGate MCP; the technical rename remains planned for Sprint 3.42. Read-only tool registration, centralized limits, and redacted stderr diagnostics are now implemented. Directory copy remains planned work and is not implied by the current `copy_path` tool. See ADR-0006 through ADR-0013 for the current architecture, security, and MCP compatibility direction.
