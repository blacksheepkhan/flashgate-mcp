# Security Model

`fileserver-mcp` is designed as a secure-by-default filesystem MCP server.

Filesystem access is security-sensitive because MCP clients may request operations on local files. For this reason, all filesystem operations are restricted to a configured sandbox root.

## Core Security Principles

### Sandbox Root

All filesystem paths are resolved relative to a configured root directory.

The root is configured through:

```text
MCP_ROOT
```

If no root is provided, the default root is:

```text
.
```

Production deployments should explicitly set `MCP_ROOT`.

### No Direct Filesystem Access Outside `internal/fs`

Production code outside `internal/fs` must not directly call filesystem APIs such as:

- `os.ReadFile`
- `os.WriteFile`
- `os.ReadDir`
- `os.Remove`
- `os.Rename`
- `os.Open`
- `os.Stat`

All filesystem access must go through the `FileSystem` abstraction.

### Central Path Validation

Path validation is centralized in:

```text
internal/security/PathGuard
```

No other package should implement independent path validation.

## PathGuard

`PathGuard` validates and resolves user-provided paths against the sandbox root.

It protects against:

- empty root paths
- absolute user paths
- parent directory traversal
- resolved paths outside the sandbox root

Path validation is intentionally layered:

1. Lexical validation normalizes the user path and rejects absolute paths and leading parent traversal.
2. Effective path validation uses evaluated filesystem paths to confirm the final existing path, or the nearest existing parent for create targets, remains inside the evaluated sandbox root.

The configured root must exist and be evaluable when the server starts. This keeps root comparisons based on the effective filesystem location rather than only string-cleaned paths.

## SafePath

`SafePath` represents a path that has passed validation.

Only `PathGuard` creates `SafePath` values.

This prevents unvalidated user input from being passed directly to filesystem operations.

## Blocked Path Types

### Absolute Paths

User-provided absolute paths are rejected.

Examples:

```text
C:\Windows\System32
/etc/passwd
```

### Parent Traversal

Parent traversal is rejected.

Examples:

```text
..
../secret.txt
../../outside-root
```

### Outside Root Resolution

Even after cleaning and resolving a path, the final effective path must remain inside the configured sandbox root.

Existing paths are evaluated directly. Create targets that do not exist yet are checked by evaluating the nearest existing parent directory before the operation is allowed.

## Destructive Operations

Destructive operations are intentionally conservative.

### Write

`Write()` does not overwrite existing files unless `overwrite=true`.

### Delete

`Delete()` does not delete non-empty directories unless `recursive=true`.

### Move

`Move()` does not overwrite existing targets unless `overwrite=true`.

### Copy

`Copy()` does not overwrite existing targets unless `overwrite=true`.

Directory copy is currently unsupported by design.

## Symlinks

Sprint 3.36 rejects symlink-based escapes where an existing path, or the nearest existing parent for a create target, resolves outside the configured root.

Sprint 3.37 adds explicit symlink policy enforcement.

Configuration:

```text
MCP_FOLLOW_SYMLINKS
```

Default: `false`.

When `MCP_FOLLOW_SYMLINKS=false`, existing symlink path components are denied before filesystem operations. Create targets are denied when the nearest existing parent contains a symlink. `list_files` filters symlink entries instead of exposing them.

When `MCP_FOLLOW_SYMLINKS=true`, classic symlinks may be followed only if the effective target remains inside the effective root. Symlink escapes are still denied. Windows junctions and non-symlink reparse points remain denied by default.

## UNC Paths

UNC path policy is enforced for configured roots and user paths.

Configuration:

```text
MCP_ALLOW_UNC_PATHS
```

Default: `false`.

When `MCP_ALLOW_UNC_PATHS=false`, UNC roots and UNC-style user paths are rejected. When `MCP_ALLOW_UNC_PATHS=true`, UNC roots do not fail solely because they are UNC paths, but the root must still exist and all root containment, hidden-file, symlink, and reparse policies still apply.

## Hidden Files

Hidden file policy is enforced for lexical dot-paths and Windows hidden attributes where available through the standard library.

Configuration:

```text
MCP_ALLOW_HIDDEN_FILES
```

Default: `false`.

When `MCP_ALLOW_HIDDEN_FILES=false`, path components whose names start with `.` are denied, except for `.` itself. Examples include `.git/config`, `.codex/settings`, and `dir/.secret`. Create targets with hidden names are denied. `list_files` filters hidden entries instead of failing the parent directory.

When `MCP_ALLOW_HIDDEN_FILES=true`, hidden and dotfile paths are allowed if all other policies pass.

## JSON-RPC Boundary

Sprint 3.38 adds JSON-RPC request validation before MCP dispatch.

Requests must be object-shaped JSON-RPC 2.0 messages. Invalid JSON, invalid request envelopes, unsupported batch requests, invalid IDs, missing methods, and malformed method params are rejected with generic JSON-RPC errors.

Responses for parse errors or invalid requests without a valid request ID include:

```json
{"jsonrpc":"2.0","id":null,"error":{"code":-32600,"message":"invalid request"}}
```

Notifications do not receive JSON-RPC responses. `notifications/initialized` is accepted as a no-op lifecycle notification. Other notifications are not executed, so `tools/call` without an `id` cannot trigger filesystem operations.

Unexpected handler panics are contained at the request boundary and returned as generic Internal error responses when the request requires a response.

## Limits and Redaction

Sprint 3.39 adds configurable hard limits for protocol input, tool arguments, filesystem payloads, and response size.

| Environment variable | Default | Scope |
|---|---:|---|
| `MCP_MAX_FILE_SIZE` | `10485760` | Hard cap for `read_file`; client `maxBytes` can only lower it. |
| `MCP_MAX_JSONRPC_MESSAGE_BYTES` | `16777216` | Maximum single JSON-RPC message read from stdin. |
| `MCP_MAX_TOOL_ARGUMENT_BYTES` | `12582912` | Maximum `tools/call` params or arguments payload. |
| `MCP_MAX_WRITE_BYTES` | `10485760` | Maximum `write_file` content size. |
| `MCP_MAX_LIST_ENTRIES` | `1000` | Maximum policy-visible `list_files` entries. |
| `MCP_MAX_COPY_BYTES` | `10485760` | Maximum `copy_path` source file size. |
| `MCP_MAX_DELETE_ENTRIES` | `1000` | Maximum entries for recursive `delete_path`. |
| `MCP_MAX_RESPONSE_BYTES` | `16777216` | Maximum serialized JSON-RPC response size safety net. |

Limit violations use generic client-visible messages. Filesystem limit denials are mapped to Invalid params with `filesystem error: limit exceeded`. JSON-RPC messages above the configured message cap are rejected as Invalid Request with `id:null`.

`MCP_DEBUG=true` enables minimal stderr diagnostics. Diagnostics are redacted for common authorization headers, token/password/API-key/secret assignments, private-key markers, connection strings with credentials, and absolute host paths. Redaction is a diagnostic safeguard; client-visible security and protocol errors are still built generically instead of exposing raw OS errors.

## Security Testing

Security tests currently cover:

- empty root rejection
- absolute path rejection
- path traversal rejection
- root normalization
- effective root validation
- symlink escape rejection
- symlink deny/follow policy
- Windows reparse point deny behavior where safely detectable
- UNC root and user path denial
- hidden dot-path and Windows hidden-attribute denial
- create-target parent validation
- safe path metadata
- filesystem traversal rejection across list/read/stat/exists/write/mkdir/delete/move/copy
- JSON-RPC envelope validation
- explicit `id:null` error responses
- notification no-response and no tool execution behavior
- generic protocol error messages
- JSON-RPC message and tool argument limits
- filesystem read/write/list/copy/delete limits
- response-size safety net
- diagnostics redaction

## Future Security Work

Planned future work:

- larger-file streaming strategy
- search tool limits and exclude model
- deeper cross-platform testing
