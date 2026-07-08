# MCP Tool Conventions

This document defines the conventions used by `fileserver-mcp` for exposing filesystem operations as MCP tools.

## Transport

`fileserver-mcp` uses MCP over JSON-RPC via STDIO.

Standard output is reserved for JSON-RPC protocol messages. Diagnostic output, logs, and process errors must be written to standard error.

## Filesystem Root

All filesystem tools operate below the configured root directory.

The root directory is configured through:

```text
MCP_ROOT
```

All tool path arguments are interpreted as paths relative to this root.

Tools must not accept or operate on paths outside the configured root. Path validation is centralized in the filesystem and security layers.

## Path Rules

Tool path arguments must be relative paths.

The following path forms are invalid or must be rejected by the filesystem layer:

- absolute paths
- paths escaping the configured root
- paths using parent traversal to leave the root
- paths whose effective filesystem location resolves outside the configured root
- invalid or empty paths where the tool requires a path

The filesystem layer validates paths in two stages. Lexical validation rejects unsafe path forms before filesystem access. Effective validation then evaluates existing paths, or the nearest existing parent for create targets, and confirms the evaluated location remains below the configured root.

Examples:

```json
{ "path": "project/README.md" }
```

```json
{ "source": "project/old.txt", "target": "project/new.txt" }
```

## Tool Naming

Tool names use lower snake case.

Examples:

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

Tool names should be stable once published because clients may depend on them.

## Input Schemas

Every tool exposes an `inputSchema`.

Schemas should:

- use JSON Schema object definitions
- set `additionalProperties` to `false`
- declare required fields explicitly
- include descriptions for all properties
- use relative path wording consistently

Example:

```json
{
  "type": "object",
  "required": ["path"],
  "additionalProperties": false,
  "properties": {
    "path": {
      "type": "string",
      "description": "Relative file path below the configured filesystem root."
    }
  }
}
```

## Result Objects

Tools return structured JSON objects.

Result objects should:

- include the primary path or source/target paths
- include a boolean operation status when useful
- avoid returning host-specific absolute paths
- avoid leaking implementation details

Examples:

```json
{
  "path": "project/file.txt",
  "written": true,
  "size": 12
}
```

```json
{
  "source": "project/old.txt",
  "target": "project/new.txt",
  "renamed": true
}
```

## Error Handling

Tool argument errors and filesystem validation errors are returned as JSON-RPC errors.

Common cases:

| Case | JSON-RPC code |
|---|---:|
| malformed tool arguments | `-32602` |
| missing required argument | `-32602` |
| invalid filesystem operation | `-32602` |

Filesystem errors are mapped centrally before being returned to the client.

Tools should not expose raw operating system errors directly unless they are intentionally wrapped and normalized.

## Overwrite Behavior

Tools that may replace an existing target use an `overwrite` boolean argument.

Default:

```json
{ "overwrite": false }
```

When `overwrite` is `false`, existing targets must be protected.

When `overwrite` is `true`, the target may be replaced if the filesystem layer supports the operation.

Affected tools:

- `write_file`
- `move_path`
- `copy_path`
- `rename_path`

## Recursive Behavior

Destructive directory deletion uses an explicit `recursive` boolean argument.

Default:

```json
{ "recursive": false }
```

Affected tool:

- `delete_path`

Non-empty directories must not be deleted unless `recursive` is explicitly set to `true`.

## Tool Order

`tools/list` returns tools in deterministic registration order.

Current order:

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

The order should remain stable to simplify client behavior, debugging, and snapshot testing.

## Security Requirements

Filesystem tools must follow these rules:

- never bypass the configured root
- never use unchecked absolute paths
- never duplicate path validation in individual tools
- keep path validation centralized
- keep protocol output and diagnostic output separated
- prefer explicit destructive flags such as `recursive` and `overwrite`

## Adding a New Tool

When adding a new tool:

1. Implement a dedicated tool file in `internal/mcp/tools`.
2. Add unit tests for definition, successful execution, malformed JSON, missing required arguments, and filesystem errors.
3. Register the tool in `cmd/server/main.go`.
4. Confirm `tools/list` order.
5. Add the tool to `docs/tools.md`.
6. Add the tool to `docs/mcp-tool-catalog.json`.
