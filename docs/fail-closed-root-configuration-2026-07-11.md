# Fail-closed root configuration

Date: 2026-07-11
Sprint: 3.44 – Codex read-only activation preparation

## Previous behavior

When `MCP_ROOT` was absent or empty, FlashGate used the process working directory (`.`). General relative root values were also resolved against that working directory. A normal file could pass initial root construction because root type was not checked explicitly.

## New behavior

`MCP_ROOT` is required. Production roots must be absolute, exist, be accessible under the existing policy, and resolve to a directory.

| Configuration | Result |
| --- | --- |
| `MCP_ROOT` absent | startup fails with `missing_root`, exit 3 |
| empty or whitespace-only | `invalid_root`, exit 3 |
| general relative root or `..` | `invalid_root`, exit 3 |
| non-existent root | `root_not_found`, exit 3 |
| file instead of directory | `root_not_directory`, exit 3 |
| policy/permission denial | `root_not_allowed`, exit 3 |

Expected startup failures write no JSON-RPC output. stderr contains one safe category and no raw operating-system error or absolute root path.

## Development current-directory opt-in

The process working directory is available only when both values are exact:

```text
MCP_ROOT=.
MCP_ALLOW_CWD_ROOT=true
```

`MCP_ALLOW_CWD_ROOT` accepts only lowercase `true` and `false`. It never supplies a missing root and never enables any other relative path. A successful Development-CWD start emits exactly one safe stderr warning.

## Client changes

Every STDIO client must now provide an explicit absolute `MCP_ROOT`. Codex read-only preparation must additionally set:

```text
MCP_READ_ONLY=true
MCP_ALLOW_CWD_ROOT=false
```

Existing clients that relied on the launch working directory must migrate to an absolute root. Development scripts may use the explicit two-variable opt-in, but production and activation examples must not.

See [Codex read-only activation preparation](codex-read-only-activation.md) for Windows, Linux, Claude Desktop, general STDIO, validation and rollback examples.

## Rollback and development guidance

If an unconverted client can no longer start, disable its FlashGate entry and restore the previous client configuration backup. Do not weaken PathGuard policy or silently reintroduce an implicit root. For temporary local development only, use the explicit CWD opt-in and retain the warning.

## Scope boundary

This change keeps one authoritative `MCP_ROOT`. It does not implement named roots, root IDs, MCP Roots, a general profile framework, a new wire error API, or real Codex activation.
