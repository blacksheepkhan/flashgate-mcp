# Filesystem tool contract cleanup

Date: 2026-07-11
Sprint: 3.43 – Pre-1.0 filesystem tool contract cleanup

## Purpose

FlashGate MCP was pre-1.0 and not productively deployed when this breaking cleanup was completed. The sprint establishes one compact filesystem baseline before client activation. It changes MCP tool contracts only; repository identity, MCP protocol revision, root defaults, security policies, tool capability enforcement, and all `MCP_*` environment-variable names remain unchanged.

## Tool migration

| Previous MCP tool | Sprint 3.43 contract |
|---|---|
| `list_files` | Rename to `list_directory` |
| `read_file` | Retained |
| `stat_path` | Rename to `get_path_info` |
| `exists_path` | Removed; genuine absence is represented by `get_path_info` |
| `write_file` | Retained |
| `mkdir` | Rename to `create_directory` |
| `delete_path` | Retained |
| `copy_path` | Retained and explicitly file-only |
| `move_path` | Retained as the single Move/Rename contract |
| `rename_path` | Removed; use `move_path` |

The new default toolset, in discovery order, is:

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

The read-only profile exposes `list_directory`, `read_file`, and `get_path_info`. There are no aliases, compatibility handlers, or deprecation shims for the five previous names.

## Input and schema changes

Every tool input is exactly one JSON object with `additionalProperties:false`. Runtime decoding rejects malformed JSON, wrong top-level types, unknown fields, trailing JSON values or non-whitespace content, wrong field types, explicit `null` field values, missing required fields, and required paths that are empty or whitespace-only. Valid path strings are not trimmed or rewritten by argument validation; PathGuard remains authoritative.

| Tool | Required | Optional |
|---|---|---|
| `list_directory` | none | `path` |
| `read_file` | `path` | `maxBytes` |
| `get_path_info` | `path` | none |
| `write_file` | `path` | `content`, `overwrite` |
| `create_directory` | `path` | none |
| `delete_path` | `path` | `recursive` |
| `copy_path` | `source`, `target` | `overwrite` |
| `move_path` | `source`, `target` | `overwrite` |

For `list_directory`, only an omitted `path` defaults to `.`; explicit blank values are invalid. For `read_file`, omitted `maxBytes` uses the server limit, values must be at least 1, and values above the server limit are safely capped.

## Result and operation changes

### `get_path_info`

An existing path returns `path`, `exists:true`, `name`, `isDir`, and `size`. A genuinely missing path is a successful call:

```json
{ "path": "missing.txt", "exists": false }
```

The implementation uses one stat path rather than an existence pre-check. Traversal, root, hidden, UNC, symlink, junction, reparse, and other policy denials remain errors and are never masked as absence. Results never contain the resolved absolute host path.

### `create_directory`

Missing parents continue to be created. A newly created leaf returns `created:true`; an already existing directory returns `created:false`; an existing file returns a path-type error.

### `copy_path`

Only files are copied. Directory and recursive copy remain unsupported and are not implicitly emulated.

### `move_path`

`move_path` handles file and directory rename and same-volume movement. Cross-volume movement returns an unsupported-operation error and never falls back to copy/delete.

The default is `overwrite:false`, so an existing target is rejected. With `overwrite:true`, only file-to-existing-file replacement is allowed. Existing directory targets and File/Directory type mismatches are rejected. Before replacement, the implementation rejects textual or safely resolved same paths, Windows case aliases, `os.SameFile` matches such as hardlinks, and lexical or symlink-resolved movement of a directory into its own subtree. Source and target identities are revalidated immediately before `os.Rename`; there is no separate target deletion. Rejected SamePath, SameFile, changed-path, self-subtree, and cross-volume operations preserve the source and perform no partial move.

A concurrent writer can still exchange the file at the already authorized target path after final revalidation because the cross-platform path-based rename API cannot condition replacement on the previously observed file identity. The residual window is restricted to that target path, never invokes a directory-removal or copy/delete fallback, and is narrower than the previous explicit remove-then-rename behavior.

## Errors

The existing JSON-RPC architecture remains. Expected tool, argument, path, policy, not-found, already-exists, type, unsupported-operation, and limit failures use Invalid params (`-32602`); unexpected I/O uses Internal error (`-32603`). Sprint-local classification distinguishes `not_found`, `already_exists`, `access_denied`, `invalid_path`, `unsupported_path_type`, `unsupported_operation`, `limit_exceeded`, and `io_error`. Messages are normalized and do not expose host paths or raw operating-system errors.

Stable machine-readable MCP error payloads, `structuredContent`, and runtime `outputSchema` are not introduced by this migration.

## Client and smoke-test updates

Clients must refresh `tools/list`, replace the three renamed calls, remove direct existence and rename calls, and update expected results for `get_path_info` and `create_directory`. A client that needs existence information must call `get_path_info`; a client that needs rename must call `move_path`.

Default smoke tests now expect eight tools, read-only smoke tests expect three, and positive coverage includes directory listing, existing and missing path info, and Move-as-Rename. Negative smoke tests verify that removed names return generic Invalid params. Windows and Bash scripts use the same JSON-RPC contract and retain deterministic temporary-file cleanup and stdout redirection.

## Scope boundary

Sprint 3.43 does not implement directory copy, named roots, fail-closed root startup, search or batch tools, targeted edits, conditional writes, job/process/command/system features, provider contracts, general JSON Schema 2020-12 infrastructure, runtime output schemas, structured content, MCP Tasks, pagination, range reads, cross-volume fallback, or general API-versioning/deprecation infrastructure.

Sprint 3.44 remains responsible for read-only client activation preparation using this cleaned three-tool read-only baseline. Later backlog items own all other deferred capabilities.
