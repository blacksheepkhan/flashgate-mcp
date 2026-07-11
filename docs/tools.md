# Filesystem MCP tools

FlashGate MCP exposes eight filesystem tools in the default profile, in this exact order:

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

The read-only profile exposes only `list_directory`, `read_file`, and `get_path_info`. Write-capable tools are not registered in read-only mode, and calls to unavailable or unknown names return generic JSON-RPC Invalid params.

For later Codex activation, `MCP_READ_ONLY=true` must be explicit and `MCP_ROOT` must be an absolute preflighted directory. See [Codex read-only activation preparation](codex-read-only-activation.md). Sprint 3.44 does not activate a client.

All paths are relative to the configured root. Absolute paths, traversal, denied hidden/UNC paths, and denied symlink, junction, or reparse access remain server-side errors. Inputs are strict JSON objects: unknown properties, malformed JSON, trailing JSON values, wrong field types, explicit `null` field values, missing required fields, and blank required paths are rejected.

## `list_directory`

Lists one directory. `path` is optional; omission means `.`, while an explicitly empty or whitespace-only value is invalid.

```json
{
  "path": "docs"
}
```

```json
{
  "entries": [
    { "name": "tools.md", "isDir": false, "size": 123 }
  ]
}
```

No pagination, filtering, recursion, or batch behavior is provided.

## `read_file`

Required: `path`. Optional: `maxBytes` with a minimum of 1. When omitted, the configured server limit is used; a larger client value is capped at that limit.

```json
{
  "path": "README.md",
  "maxBytes": 4096
}
```

```json
{
  "content": "...",
  "size": 123
}
```

Range reads are not implemented.

## `get_path_info`

Required: `path`. A single metadata lookup provides both existence and metadata without an existence pre-check.

Existing path:

```json
{
  "path": "README.md",
  "exists": true,
  "name": "README.md",
  "isDir": false,
  "size": 123
}
```

Missing path:

```json
{
  "path": "missing.txt",
  "exists": false
}
```

Only genuine missing-path errors become `exists:false`. Security and policy denials remain errors, and no absolute host path is returned.

## `write_file`

Required: `path`. Optional: `content` (empty is allowed) and `overwrite` (default `false`). Existing limits and root/security enforcement apply.

```json
{
  "path": "output.txt",
  "content": "text",
  "overwrite": false
}
```

```json
{
  "path": "output.txt",
  "size": 4,
  "written": true
}
```

## `create_directory`

Required: `path`. Missing parents are created. `created` describes the actual leaf state.

```json
{ "path": "output/archive" }
```

New leaf:

```json
{ "path": "output/archive", "created": true }
```

Existing directory:

```json
{ "path": "output/archive", "created": false }
```

An existing file at the target is a path-type error.

## `delete_path`

Required: `path`. Optional: `recursive`, default `false`. Non-empty directories require explicit recursive deletion and remain bounded by the configured deletion limit.

```json
{ "path": "output/archive", "recursive": true }
```

```json
{ "path": "output/archive", "deleted": true }
```

## `copy_path`

Required: `source` and `target`. Optional: `overwrite`, default `false`.

```json
{ "source": "a.txt", "target": "b.txt", "overwrite": false }
```

```json
{ "source": "a.txt", "target": "b.txt", "copied": true }
```

This contract copies files only. Directory copy and recursive copy are not supported.

## `move_path`

Required: `source` and `target`. Optional: `overwrite`, default `false`. This is the single contract for both move and rename.

```json
{ "source": "old.txt", "target": "new.txt", "overwrite": false }
```

```json
{ "source": "old.txt", "target": "new.txt", "moved": true }
```

Files and directories may be renamed or moved on the same volume. Cross-volume moves are rejected without copy/delete fallback. With `overwrite:true`, only file-to-existing-file replacement is allowed. Existing directory targets and all directory replacement combinations are rejected. Same-path, same-file (including detectable hardlinks and Windows case aliases), and lexical or symlink-resolved directory-into-own-subtree operations are rejected before replacement.

The source and target identities are revalidated immediately before the operating-system rename, and existing files are replaced through `os.Rename` without a separate target deletion. The remaining race is limited to a concurrent change at the already authorized target path after final revalidation; the path-based cross-platform API cannot condition replacement on the previously observed file identity. Such a race cannot trigger a directory-removal fallback, and this behavior is narrower than the previous explicit remove-then-rename sequence.

## Errors

Parse, invalid-request, and method errors use the standard JSON-RPC codes. Expected argument, path, policy, not-found, already-exists, path-type, unsupported-operation, and limit failures use `-32602`. Unexpected I/O failures use `-32603`. Error messages are normalized and do not expose absolute host paths or raw operating-system details.

Stable machine-readable MCP tool-error payloads, `structuredContent`, and runtime `outputSchema` are not part of Sprint 3.43.

The previous pre-1.0 contract and required client changes are documented in [filesystem tool contract cleanup](filesystem-tool-contract-cleanup-2026-07-11.md).
