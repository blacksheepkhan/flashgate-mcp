# MCP Tool Reference

This document describes the MCP tools exposed by `fileserver-mcp`.

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

Sprint 3.35 adds read-only tool capability gating. In read-only mode, write-capable tools are not registered. Direct `tools/call` requests for `write_file`, `mkdir`, `delete_path`, `move_path`, `copy_path`, or `rename_path` return a tool-not-found error.

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
- invalid path
- path outside the configured root
- target already exists and `overwrite` is `false`
- directory is not empty and `recursive` is `false`

---

## `list_files`

Lists files and directories below the configured filesystem root.

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

Copies a file or directory below the configured filesystem root.

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
