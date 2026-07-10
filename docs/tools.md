# MCP Tool Reference

This document describes the tools currently exposed by the technical `fileserver-mcp` implementation of FlashGate MCP.

> The current tool names remain active until the dedicated pre-1.0
> tool contract cleanup sprint. This document distinguishes implemented
> tools from accepted planned changes.

## Contract Status

FlashGate MCP is pre-1.0, is not yet used in production, and has no external user compatibility contract. Tool names, parameters, and result schemas may therefore be changed or removed before the first stable release. Changes still require changelog, tests, tool documentation, client examples, and smoke-test updates.

Sprint 3.41 changes no actual tool schema. The authoritative machine-readable current catalog remains unchanged in `mcp-tool-catalog.json`. The implemented protocol remains MCP `2025-11-25`; future schemas will be validated as JSON Schema 2020-12.

### Accepted planned cleanup

| Current implemented tool | Accepted planned decision |
|---|---|
| `list_files` | Rename to `list_directory` |
| `read_file` | Keep; later add line/byte/head/tail ranges |
| `stat_path` | Rename to `get_path_info` |
| `exists_path` | Remove; replace with `get_path_info` |
| `write_file` | Keep; later add safe write modes |
| `mkdir` | Rename to `create_directory` |
| `delete_path` | Keep |
| `copy_path` | Keep |
| `move_path` | Keep and define both move and rename semantics |
| `rename_path` | Remove; replace with `move_path` |

Planned cleaned baseline:

```text
list_directory
read_file
get_path_info
write_file
create_directory
delete_path
copy_path
move_path
```

A future `get_path_info` is planned to represent a missing path structurally:

```json
{
  "exists": false,
  "path": "missing.txt"
}
```

Other failures remain distinguishable, with candidate categories `not_found`, `access_denied`, `outside_allowed_root`, `invalid_path`, `unsupported_path_type`, and `io_error`. Final error codes are deferred to the tool-contract sprint.

### Planned new operation types

Future, not implemented operation types include paginated listing, ranged and batch reads, batch path inspection and hashing, bounded trees and search, targeted/conditional writes, internal Operations/Job handles, managed process operations, allowlisted command execution, and controlled system information. Custom status/result/cancel tools are not the accepted primary MCP job contract. The MCP adapter will first evaluate mapping eligible internal jobs to the negotiated Tasks Extension `io.modelcontextprotocol/tasks` and decide bounded behavior for clients without Tasks support.

`fileserver-mcp` uses MCP over JSON-RPC via STDIO. Tools are invoked through the MCP `tools/call` method.

All paths are relative to the configured filesystem root.

The filesystem root is configured through:

```text
MCP_ROOT
```

## Tool Order

`tools/list` returns tools in deterministic registration order.

Default mode:

```text
list_files
read_file
stat_path
exists_path
write_file
mkdir
delete_path
move_path
copy_path
rename_path
```

Read-only mode with `MCP_READ_ONLY=true`:

```text
list_files
read_file
stat_path
exists_path
```

Sprint 3.35 adds read-only tool capability gating. In read-only mode, write-capable tools are not registered. Direct `tools/call` requests for `write_file`, `mkdir`, `delete_path`, `move_path`, `copy_path`, or `rename_path` return a generic Invalid params error, matching unknown tool names without revealing whether the tool exists in another mode.

## Common Error Behavior

Most invalid tool requests return JSON-RPC `Invalid params`:

```json
{
  "code": -32602,
  "message": "..."
}
```

Common causes:

- malformed JSON arguments
- missing required arguments
- unknown tool name
- direct call for a write-capable tool while `MCP_READ_ONLY=true`
- invalid path
- path outside the configured root
- hidden, UNC, symlink, junction, or reparse policy denial
- target already exists and `overwrite` is `false`
- directory is not empty and `recursive` is `false`

Protocol-level JSON-RPC errors are generic. Invalid request envelopes return `invalid request`, unknown JSON-RPC methods return `method not found`, invalid method params return `invalid params`, and unexpected server errors return `internal error`.

Limit violations are also generic. Oversized JSON-RPC messages return Invalid Request with `id:null`. Oversized `tools/call` arguments return Invalid params. Filesystem operation limits return Invalid params with:

```text
filesystem error: limit exceeded
```

---

## `list_files`

Lists files and directories below the configured filesystem root.

The number of policy-visible entries is capped by `MCP_MAX_LIST_ENTRIES`. If more entries would be returned, the tool fails with a limit error instead of silently truncating the listing.

When hidden files or symlinks/reparse points are denied by policy, matching child entries are filtered from the result instead of failing the whole directory listing. The parent directory itself must still pass all path and policy checks.

### Input

```json
{
  "path": "relative/directory"
}
```

### Fields

| Field | Type | Required | Description |
|---|---|---:|---|
| `path` | string | no | Relative directory path below the configured filesystem root. Defaults to `"."`. |

### Result

```json
{
  "entries": [
    {
      "name": "README.md",
      "isDir": false,
      "size": 1234
    }
  ]
}
```

### Example

```json
{
  "name": "list_files",
  "arguments": {
    "path": "fileserver-mcp"
  }
}
```

---

## `read_file`

Reads a text file below the configured filesystem root.

The maximum returned content is capped by `MCP_MAX_FILE_SIZE`. The optional `maxBytes` argument can lower this cap for a request, but cannot raise it.

### Input

```json
{
  "path": "relative/file.txt",
  "maxBytes": 8192
}
```

### Fields

| Field | Type | Required | Description |
|---|---|---:|---|
| `path` | string | yes | Relative file path below the configured filesystem root. |
| `maxBytes` | integer | no | Maximum number of bytes to read. Defaults to the configured maximum file size. |

### Result

```json
{
  "content": "file content",
  "size": 12
}
```

### Example

```json
{
  "name": "read_file",
  "arguments": {
    "path": "fileserver-mcp/README.md",
    "maxBytes": 8192
  }
}
```

---

## `stat_path`

Returns metadata for a file or directory below the configured filesystem root.

### Input

```json
{
  "path": "relative/path"
}
```

### Fields

| Field | Type | Required | Description |
|---|---|---:|---|
| `path` | string | yes | Relative file or directory path below the configured filesystem root. |

### Result

```json
{
  "name": "README.md",
  "isDir": false,
  "size": 1234
}
```

### Example

```json
{
  "name": "stat_path",
  "arguments": {
    "path": "fileserver-mcp/README.md"
  }
}
```

---

## `exists_path`

Checks whether a file or directory exists below the configured filesystem root.

### Input

```json
{
  "path": "relative/path"
}
```

### Fields

| Field | Type | Required | Description |
|---|---|---:|---|
| `path` | string | yes | Relative file or directory path below the configured filesystem root. |

### Result

```json
{
  "exists": true
}
```

### Example

```json
{
  "name": "exists_path",
  "arguments": {
    "path": "fileserver-mcp/README.md"
  }
}
```

---

## `write_file`

Writes a text file below the configured filesystem root.

Content size is capped by `MCP_MAX_WRITE_BYTES`.

### Input

```json
{
  "path": "relative/file.txt",
  "content": "file content",
  "overwrite": false
}
```

### Fields

| Field | Type | Required | Description |
|---|---|---:|---|
| `path` | string | yes | Relative file path below the configured filesystem root. |
| `content` | string | no | Text content to write. An empty string is allowed. |
| `overwrite` | boolean | no | Whether an existing file may be overwritten. Defaults to `false`. |

### Result

```json
{
  "path": "relative/file.txt",
  "size": 12,
  "written": true
}
```

### Example

```json
{
  "name": "write_file",
  "arguments": {
    "path": "fileserver-mcp/tmp.txt",
    "content": "hello",
    "overwrite": false
  }
}
```

---

## `mkdir`

Creates a directory below the configured filesystem root.

### Input

```json
{
  "path": "relative/new-directory"
}
```

### Fields

| Field | Type | Required | Description |
|---|---|---:|---|
| `path` | string | yes | Relative directory path below the configured filesystem root. |

### Result

```json
{
  "path": "relative/new-directory",
  "created": true
}
```

### Example

```json
{
  "name": "mkdir",
  "arguments": {
    "path": "fileserver-mcp/tmp-dir"
  }
}
```

---

## `delete_path`

Deletes a file or directory below the configured filesystem root.

Recursive deletes are capped by `MCP_MAX_DELETE_ENTRIES`. If the target tree exceeds the limit, no delete is performed.

### Input

```json
{
  "path": "relative/path",
  "recursive": false
}
```

### Fields

| Field | Type | Required | Description |
|---|---|---:|---|
| `path` | string | yes | Relative file or directory path below the configured filesystem root. |
| `recursive` | boolean | no | Whether a non-empty directory may be deleted recursively. Defaults to `false`. |

### Result

```json
{
  "path": "relative/path",
  "deleted": true
}
```

### Example

```json
{
  "name": "delete_path",
  "arguments": {
    "path": "fileserver-mcp/tmp.txt",
    "recursive": false
  }
}
```

---

## `move_path`

Moves a file or directory below the configured filesystem root.

### Input

```json
{
  "source": "relative/source.txt",
  "target": "relative/target.txt",
  "overwrite": false
}
```

### Fields

| Field | Type | Required | Description |
|---|---|---:|---|
| `source` | string | yes | Relative source file or directory path below the configured filesystem root. |
| `target` | string | yes | Relative target file or directory path below the configured filesystem root. |
| `overwrite` | boolean | no | Whether an existing target may be overwritten. Defaults to `false`. |

### Result

```json
{
  "source": "relative/source.txt",
  "target": "relative/target.txt",
  "moved": true
}
```

### Example

```json
{
  "name": "move_path",
  "arguments": {
    "source": "fileserver-mcp/source.txt",
    "target": "fileserver-mcp/target.txt",
    "overwrite": false
  }
}
```

---

## `copy_path`

Copies a file below the configured filesystem root.

Source file size is capped by `MCP_MAX_COPY_BYTES`. Directory copy remains unsupported.

### Input

```json
{
  "source": "relative/source.txt",
  "target": "relative/target.txt",
  "overwrite": false
}
```

### Fields

| Field | Type | Required | Description |
|---|---|---:|---|
| `source` | string | yes | Relative source file or directory path below the configured filesystem root. |
| `target` | string | yes | Relative target file or directory path below the configured filesystem root. |
| `overwrite` | boolean | no | Whether an existing target may be overwritten. Defaults to `false`. |

### Result

```json
{
  "source": "relative/source.txt",
  "target": "relative/target.txt",
  "copied": true
}
```

### Example

```json
{
  "name": "copy_path",
  "arguments": {
    "source": "fileserver-mcp/source.txt",
    "target": "fileserver-mcp/copy.txt",
    "overwrite": false
  }
}
```

---

## `rename_path`

Renames a file or directory below the configured filesystem root.

### Input

```json
{
  "source": "relative/old-name.txt",
  "target": "relative/new-name.txt",
  "overwrite": false
}
```

### Fields

| Field | Type | Required | Description |
|---|---|---:|---|
| `source` | string | yes | Relative source file or directory path below the configured filesystem root. |
| `target` | string | yes | Relative target file or directory path below the configured filesystem root. |
| `overwrite` | boolean | no | Whether an existing target may be overwritten. Defaults to `false`. |

### Result

```json
{
  "source": "relative/old-name.txt",
  "target": "relative/new-name.txt",
  "renamed": true
}
```

### Example

```json
{
  "name": "rename_path",
  "arguments": {
    "source": "fileserver-mcp/old.txt",
    "target": "fileserver-mcp/new.txt",
    "overwrite": false
  }
}
```
