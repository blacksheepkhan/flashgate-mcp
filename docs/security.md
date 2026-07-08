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

## Future Security Work

Planned future work:

- maximum response size
- maximum directory listing size
- audit logging
