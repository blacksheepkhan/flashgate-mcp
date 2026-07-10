# MCP Tool Conventions

This document defines the conventions used by FlashGate MCP. The current technical implementation remains `fileserver-mcp` until Sprint 3.42.

## Transport

`fileserver-mcp` uses MCP over JSON-RPC via STDIO.

Standard output is reserved for JSON-RPC protocol messages. Diagnostic output, logs, and process errors must be written to standard error.

JSON-RPC request envelopes are validated before dispatch. Requests must be JSON objects with `jsonrpc` set to `"2.0"`, a non-empty string `method`, and an optional `id` of type string, number, or null. Unsupported batch requests are rejected. Notifications do not receive responses and are not used to execute tools.

Single JSON-RPC messages are bounded by `MCP_MAX_JSONRPC_MESSAGE_BYTES`. Oversized messages are rejected as Invalid Request with `id:null` and are not dispatched.

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

Tool names use understandable, fully written verbs. Platform-specific or Unix jargon is avoided when a clear general term exists. Tool names become stability-sensitive after 1.0; justified breaking changes are allowed before 1.0 because no external compatibility contract exists yet.

## Input Schemas

Every tool exposes an `inputSchema`.

Schemas should:

- use JSON Schema object definitions
- set `additionalProperties` to `false`
- declare required fields explicitly
- include descriptions for all properties
- use relative path wording consistently
- validate future input and output contracts against JSON Schema 2020-12

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
| unknown tool name | `-32602` |
| read-only-gated write tool name | `-32602` |
| invalid filesystem operation | `-32602` |
| filesystem limit exceeded | `-32602` |

Filesystem errors are mapped centrally before being returned to the client.

Tools should not expose raw operating system errors directly unless they are intentionally wrapped and normalized.

Filesystem limit errors should use the generic client-visible message `filesystem error: limit exceeded`. Tool argument and `tools/call` payload limits should use generic Invalid params errors.

Protocol-level errors use generic JSON-RPC messages:

| Case | JSON-RPC code | Message |
|---|---:|---|
| invalid JSON | `-32700` | `parse error` |
| invalid request envelope | `-32600` | `invalid request` |
| unknown JSON-RPC method | `-32601` | `method not found` |
| invalid method params | `-32602` | `invalid params` |
| unexpected internal failure | `-32603` | `internal error` |

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
- enforce hidden, UNC, symlink, junction, and reparse policy through `internal/security.PathGuard`
- enforce configured message, argument, filesystem operation, and response limits
- map path and security policy denials to generic invalid-path tool errors
- map limit denials to generic limit errors
- keep protocol output and diagnostic output separated
- keep diagnostics on stderr and redact common secrets, credentials, connection strings, and host paths
- prefer explicit destructive flags such as `recursive` and `overwrite`

## Adding a New Tool

When adding a new tool:

1. Implement a dedicated tool file in `internal/mcp/tools`.
2. Add unit tests for definition, successful execution, malformed JSON, missing required arguments, and filesystem errors.
3. Register the tool in `cmd/server/main.go`.
4. Confirm `tools/list` order.
5. Add the tool to `docs/tools.md`.
6. Add the tool to `docs/mcp-tool-catalog.json`.

## Accepted Pre-1.0 Contract Conventions

The following conventions are accepted target rules. They do not change current tool schemas or registration in Sprint 3.41.

- One tool has one clear domain responsibility.
- Redundant tools are removed before 1.0 instead of receiving artificial compatibility aliases.
- Batch tools are added only when they measurably reduce calls or response size.
- Results are structured, compact, and bounded.
- Potentially large result sets provide cursor pagination.
- Long-running work uses opaque operation handles.
- MCP annotations describe behavior but never replace server-side authorization.
- Destructive operations use explicit semantics and, where useful, dry-run support or preconditions.
- Command execution has no free shell string as its standard parameter.
- A free-form workflow or shell language must not replace clearly defined tools.

Accepted planned filesystem cleanup:

| Current | Planned pre-1.0 decision |
|---|---|
| `list_files` | Rename to `list_directory` |
| `read_file` | Retain and later add bounded ranges |
| `stat_path` | Rename to `get_path_info` |
| `exists_path` | Remove; cover with `get_path_info` |
| `write_file` | Retain and add safe write modes |
| `mkdir` | Rename to `create_directory` |
| `delete_path` | Retain |
| `copy_path` | Retain |
| `move_path` | Retain and cover rename semantics |
| `rename_path` | Remove; cover with `move_path` |

The planned baseline is:

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

These changes require coordinated schema snapshots, tests, smoke tests, tool documentation, client examples, and changelog entries in Sprint 3.43. Before the first stable release, the project will define versioning, deprecation, and migration policy.
