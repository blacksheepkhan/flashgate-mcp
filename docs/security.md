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

A later security sprint will define the full symlink, junction, and reparse policy based on configuration.

Planned configuration:

```text
MCP_FOLLOW_SYMLINKS
```

## UNC Paths

UNC path handling is not yet fully implemented.

Planned configuration:

```text
MCP_ALLOW_UNC_PATHS
```

The default will remain secure: UNC paths are not allowed unless explicitly enabled.

## Hidden Files

Hidden file handling is not yet fully implemented.

Planned configuration:

```text
MCP_ALLOW_HIDDEN_FILES
```

## Security Testing

Security tests currently cover:

- empty root rejection
- absolute path rejection
- path traversal rejection
- root normalization
- effective root validation
- symlink escape rejection
- create-target parent validation
- safe path metadata
- filesystem traversal rejection across list/read/stat/exists/write/mkdir/delete/move/copy

## Future Security Work

Planned future work:

- full symlink, junction, and reparse policy
- Windows UNC path validation
- hidden file policy
- maximum response size
- maximum directory listing size
- audit logging
