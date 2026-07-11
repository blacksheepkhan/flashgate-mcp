# FlashGate MCP Security Model

FlashGate MCP is designed as a secure-by-default local host-operations MCP server. Sprint 3.42 completed its technical rename; it currently exposes only the filesystem functionality described below.

Filesystem access is security-sensitive because MCP clients may request operations on local files. For this reason, all filesystem operations are restricted to a configured sandbox root.

## Core Security Principles

### Sandbox Root

All filesystem paths are resolved relative to a configured root directory.

The root is configured through:

```text
MCP_ROOT
```

`MCP_ROOT` is required. Production roots must be absolute, exist, be accessible under the current policy, and resolve to a directory. Missing, empty, whitespace-only, relative, non-existent, file, and policy-denied roots stop startup before Filesystem, Registry, Router or JSON-RPC processing.

Expected configuration/root failures use safe categories on stderr and exit code `3`; stdout remains empty. Unexpected bootstrap failures use `startup_failed` and exit code `1`. Raw roots and operating-system errors are not emitted.

The process working directory is available only for explicit development use:

```text
MCP_ROOT=.
MCP_ALLOW_CWD_ROOT=true
```

The opt-in accepts only lowercase `true` or `false`, never supplies a missing root, and never enables other relative roots. A successful CWD-development start emits one safe stderr warning. Production and Codex examples set `MCP_ALLOW_CWD_ROOT=false`.

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

The configured root must exist, pass policy, and resolve effectively to a directory when the server starts. This keeps root comparisons based on the effective filesystem location rather than only string-cleaned paths.

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

`Move()` does not overwrite existing targets unless `overwrite=true`. Replacement is restricted to file-to-file, revalidates the observed source and target identities immediately before the rename, and uses `os.Rename` without a separate target deletion. Directory targets are never explicitly removed.

The standard cross-platform API remains path-based: a concurrent writer could exchange the file at the already authorized target path after final revalidation. That residual race is bounded to the target path, cannot invoke directory removal or copy/delete fallback, and is narrower than the previous remove-then-rename sequence.

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

When `MCP_FOLLOW_SYMLINKS=false`, existing symlink path components are denied before filesystem operations. Create targets are denied when the nearest existing parent contains a symlink. `list_directory` filters symlink entries instead of exposing them.

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

When `MCP_ALLOW_HIDDEN_FILES=false`, path components whose names start with `.` are denied, except for `.` itself. Examples include `.git/config`, `.codex/settings`, and `dir/.secret`. Create targets with hidden names are denied. `list_directory` filters hidden entries instead of failing the parent directory.

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
| `MCP_MAX_LIST_ENTRIES` | `1000` | Maximum policy-visible `list_directory` entries. |
| `MCP_MAX_COPY_BYTES` | `10485760` | Maximum `copy_path` source file size. |
| `MCP_MAX_DELETE_ENTRIES` | `1000` | Maximum entries for recursive `delete_path`. |
| `MCP_MAX_RESPONSE_BYTES` | `16777216` | Maximum serialized JSON-RPC response size safety net. |

Limit violations use generic client-visible messages. Filesystem limit denials are mapped to Invalid params with `filesystem error: limit exceeded`. JSON-RPC messages above the configured message cap are rejected as Invalid Request with `id:null`.

`MCP_DEBUG=true` enables minimal stderr diagnostics. Diagnostics are redacted for common authorization headers, token/password/API-key/secret assignments, private-key markers, connection strings with credentials, and absolute host paths. Redaction is a diagnostic safeguard; client-visible security and protocol errors are still built generically instead of exposing raw OS errors.

## Security Testing

Security tests currently cover:

- empty root rejection
- missing, whitespace, relative and development-CWD root contracts
- root existence and directory type
- categorized startup errors, exit codes, empty startup stdout and host-path redaction
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
- filesystem traversal rejection across list/read/info/write/create/delete/copy/move
- JSON-RPC envelope validation
- explicit `id:null` error responses
- notification no-response and no tool execution behavior
- generic protocol error messages
- JSON-RPC message and tool argument limits
- filesystem read/write/list/copy/delete limits
- response-size safety net
- diagnostics redaction

Startup preflight completes before any tool Registry, Router or MCP server is created. Normal starts remain silent; diagnostics never share JSON-RPC stdout.

## Future Security Work

Planned future work:

- larger-file streaming strategy
- search tool limits and exclude model
- deeper cross-platform testing

## Accepted Target Security Architecture

This section records accepted target controls. Except where the current-state sections above say otherwise, these controls are planned and are not implemented in Sprint 3.41.

### Gate as a server-enforced boundary

The “Gate” in FlashGate means a server-enforced control boundary. Tool visibility, client claims, and MCP annotations are not authorization. Every operation must pass the applicable server-side capabilities, profile, root policy, path or process policy, limits, redaction, audit, and platform checks.

FlashGate modules/providers cannot bypass this boundary. Public, community, vendor, organization-internal, and Voxtronic-specific providers use the same central controls as the core. MCP protocol-extension negotiation is separate and does not grant authorization.

### Capability enforcement and tool profiles

Planned capabilities include:

```text
filesystem.read
filesystem.write
search.execute
process.observe
process.manage
process.control.external
command.execute
system.read
```

Profiles combine functional capabilities and policy and determine tool registration. Example profile names include `safe-read`, `filesystem-write`, and `admin`. `high-risk`, `destructive`, and `interactive` are risk classifications or additional policy conditions, not universal capabilities. The current `MCP_READ_ONLY` registration behavior is the first restricted-profile case, not the final profile model. Final names remain open for Sprint 3.50. Authorization is checked again during execution so direct calls cannot bypass registration rules.

### Per-root policies

The accepted named-root target is based on authoritative FlashGate server configuration and explicit root IDs plus relative paths. Each root may define read/write permission, file and result size, allowed file types, capability mapping, symlink/reparse rules, and process working-directory permission. Deprecated MCP Roots is not an architectural dependency; optional legacy compatibility may be evaluated only for a supported client and never overrides server policy.

### Operations and Job Manager controls

The planned Operations/Job Manager is an optional runtime service for long-running or managed work. Short synchronous work may run directly in domain services. The generic manager does not own domain logic or MCP tool types. Managed operations use opaque handles, bounded queues, global and per-domain concurrency, maximum runtime, bounded results and temporary data, TTL cleanup, controlled shutdown, and job-leak protection.

The server owns deadlines. Workers receive cancellable contexts and check cancellation regularly; a watchdog enforces expiry. External workers are terminated if required. Temporary resources are removed or marked incomplete, status becomes `timed_out`, and errors/results remain bounded. A worker is never trusted as the only deadline enforcer.

### Managed process identity and control

Server-started processes receive opaque process handles. Handles, not PIDs alone, identify managed instances for status, output, wait, and stop operations. PID reuse must not cause a request to control the wrong process.

Stopping defaults to server-managed processes. Controlling an external PID requires a distinct functional capability such as `process.control.external` plus a high-risk policy classification and is outside standard profiles. Process listings, details, command lines, and environment information must be filtered and redacted.

### Command execution boundary

Planned command execution uses configured executable IDs resolved to server-approved absolute program paths. Arguments are separate arrays; the standard interface has no free shell string. Working directories must be inside permitted roots. Environment variables use an allowlist or explicit propagation rules. stdout and stderr remain separate and bounded. Runtime, output, and process concurrency are limited.

A future synchronous `run_command` will wrap the same Managed Process Engine. A second execution engine is prohibited. Interactive shells remain disabled and require distinct interactive/high-risk policy decisions. Platform-specific Windows and Linux isolation must preserve the same policy outcome and least-privilege intent.

### Secret redaction and audit events

Secret redaction applies to client-visible results, diagnostics, process command lines, environment data, audit fields, and module/provider output where applicable. Redaction complements data minimization; it does not make unrestricted collection acceptable.

A planned audit event model will record bounded, structured security-relevant facts such as effective capability, root ID, operation type, policy decision, handle, duration, outcome, and redacted failure classification. Audit data must not include raw secrets, unrestricted file content, full environments, or unnecessary absolute host paths.

### FlashGate module/provider and MCP extension boundaries

All FlashGate modules/providers must use shared:

- capabilities
- root policies
- limits
- path validation
- process policies
- secret redaction
- audit events
- platform adapters

A provider's vendor, deployment location, or official/community label grants no implicit trust. The later provider-runtime decision must include isolation, update, dependency, and supply-chain analysis. MCP protocol extensions use their own negotiated wire contract, but negotiation remains distinct from capability authorization.

### MCP version and Tasks compatibility

The implemented protocol remains MCP `2025-11-25`. No later MCP feature is supported in Sprint 3.41. The MCP adapter owns protocol and extension negotiation. If eligible internal jobs are later exposed through `io.modelcontextprotocol/tasks`, each Task request must be bound to the authorized caller and internal handle, and internal states/results must be mapped and redacted deliberately. Custom operation tools are not the accepted primary contract while this compatibility decision remains open.

## Threat Model Workstreams

Separate threat models are required before the corresponding target domains become generally available:

| Workstream | Minimum concerns |
|---|---|
| Filesystem | traversal, symlink/reparse escapes, races, overwrite/delete semantics, resource exhaustion, data disclosure |
| Search | unbounded recursion, scanned-byte cost, binary/encoding behavior, ignored sensitive paths, result leakage |
| Processes | PID reuse, lifecycle races, command-line/environment disclosure, external PID control, orphan cleanup |
| Command execution | executable substitution, argument and environment injection, working-directory escape, output/resource exhaustion, platform isolation |
| FlashGate modules/providers | policy bypass, capability inflation, dependency/supply-chain risk, metadata trust, in-process versus IPC isolation |
| MCP protocol extensions | negotiation downgrade/mismatch, capability confusion, version compatibility, authorization separation |

Stateful components additionally require race-detector coverage, restart/shutdown analysis, handle lifecycle tests, negative capability tests, and cleanup verification.

## Deferred Security Decisions

- final profile and capability configuration schema
- exact audit event schema and retention
- CPU/RAM isolation mechanism per platform
- FlashGate module/provider runtime and isolation model
- MCP protocol-extension compatibility beyond `2025-11-25`
- external PID control design
- interactive input and shell design
- privacy-sensitive network-information scope
