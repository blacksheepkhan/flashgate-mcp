# FlashGate MCP Architecture

FlashGate MCP is a fast, secure, resource-efficient, local-first MCP server for controlled filesystem, process, and operating-system operations.

## Project identity and transition

The public project name is **FlashGate MCP** as of Sprint 3.41. Flash describes low latency, efficient local processing, and compact responses. Gate describes the server-enforced boundary formed by policies, capabilities, roots, limits, redaction, and audit events.

The current repository, Go module, binary, MCP server implementation name (`serverInfo.name`), package paths, scripts, workflows, and machine-readable tool catalog still use the legacy technical identifier `fileserver-mcp`. Sprint 3.42 is the dedicated technical rename. Until then, both names appear only where this transition or the current implementation must be described.

Planned technical identifiers, not implemented in Sprint 3.41:

| Item | Target |
|---|---|
| Repository | `flashgate-mcp` |
| Binary | `flashgate-mcp` |
| MCP server implementation name (`serverInfo.name`) | `flashgate` |
| Go module | `github.com/blacksheepkhan/flashgate-mcp` |
| Short name | FlashGate |

FlashGate MCP is not a web-hosting service and not a remote-shell replacement. It is a controlled local boundary between MCP clients and explicitly enabled host functionality.

## Architectural goals

- predictable, testable behavior
- low latency, CPU use, memory use, and response size
- local deterministic work instead of model round trips
- cross-platform Windows and Linux support
- secure-by-default access
- explicit policy, capability, limit, redaction, and audit boundaries
- clear dependency direction and domain ownership
- one repository and one primary binary until evidence supports a split
- vendor-neutral open-source core
- no dependency on external MCP protocol libraries

## Current state

The current implementation is a layered Go application using MCP JSON-RPC over STDIO. It provides:

- configuration from environment variables
- JSON-RPC validation, routing, MCP initialization, `tools/list`, and `tools/call`
- ten filesystem tools documented in `tools.md`
- one configured filesystem root through `MCP_ROOT`
- read-only tool registration through `MCP_READ_ONLY`
- central `PathGuard` validation and filesystem abstraction
- hard protocol, argument, filesystem, and response limits
- redacted debug diagnostics on stderr
- Windows and Linux tests and smoke tests
- MCP protocol version `2025-11-25`

The implemented dependency path is:

```text
MCP Client
    |
STDIO transport
    |
JSON-RPC / MCP server -> router -> handlers -> tools
    |
Filesystem abstraction
    |
PathGuard and current policies
    |
Operating-system filesystem
```

The current implementation does not yet contain named roots, general capability profiles, an Operations/Job Manager, search, process management, command execution, system-information tools, a FlashGate module/provider system, MCP protocol extensions, or separate Windows/Linux domain adapters.

## Accepted target architecture

```text
                         MCP Adapter
       Version/Extension Negotiation, JSON-RPC, Schemas
                              |
                  Profiles and Capabilities
                              |
            Policy, Limits, Audit, Redaction
                              |
 ┌────────────┬──────────┬─────────┬───────────┬─────────┐
 │ Filesystem │  Search  │ Process │ Execution │ System  │
 └────────────┴──────────┴─────────┴───────────┴─────────┘
          |                                   |
          | direct short operations           + - - - - - - - - +
          v                                   optional lifecycle  |
                 Windows and Linux Adapters              v
                                            Operations / Job Manager
                                            (generic runtime service)
```

This diagram is the accepted target, not the current package layout.

### Dependency direction

The local system core is separated conceptually from the MCP/JSON-RPC adapter. Core components must not depend on MCP, JSON-RPC, tool schema, or tool response types.

The MCP adapter owns:

- transport
- tool registration and tool schemas
- input validation
- error translation
- structured result serialization

Operating-system logic stays below that adapter. Go components in this codebase reuse the core directly. Future MCPs must not use MCP-to-MCP calls as their internal architecture. Local IPC is introduced only when process isolation or another programming language justifies it. The core must remain extractable, but no prematurely stable public Go API is promised.

Target dependency examples:

```text
mcp adapter -> profiles/capabilities -> policy/limits/audit
domain services -> platform adapters
domain orchestration -. optional long-running lifecycle .-> operations/jobs
filesystem/search/process/execution/system -> shared policy services
platform adapters -X-> mcp adapter
core domains     -X-> JSON-RPC or tool response types
operations/jobs  -X-> MCP tool types or domain business logic
```

### Domain model

#### Filesystem

Owns files, directories, metadata, reads, writes, copying, moving, deletion, hashing, directory size, and controlled file changes. Files and directories remain one domain. Long filesystem work still belongs to this domain even when a job executes it.

#### Search

Owns path and filename search, metadata filtering, literal and regular-expression text search, bounded recursion, and optional future accelerators or indexes.

#### Process

Owns process observation, details, process trees, managed process instances, opaque process handles, status, and lifecycle.

#### Execution

Owns controlled program execution, executable allowlisting, argument validation, working-directory and environment policies, runtime/output limits, and platform isolation.

#### System information

Owns only explicitly released operating-system data, architecture, host resources, scoped disk usage, filtered environment information, and possibly later restricted network information.

#### Operations and jobs

Operations/jobs are a shared, optional technical runtime service, not a competing business domain and not a mandatory layer between every domain and its platform adapter. Short synchronous work may run directly in a domain service. Long-running or managed work may use the service for lifecycle, cancellation, deadlines, progress, and cleanup while domain logic and ownership remain outside the generic manager.

#### Cross-cutting components

- policy and capabilities
- limits
- diagnostics, audit, and secret redaction
- Windows platform adapters
- Linux platform adapters
- MCP adapter

## Operations and Job Manager

Long-running internal work may use an opaque handle shaped like `op_<opaque-id>`. The internal handle and state model are not an accepted custom MCP tool contract. Future MCP exposure must first decide compatibility with the official Tasks Extension `io.modelcontextprotocol/tasks` and supported clients.

Accepted statuses:

- `queued`
- `running`
- `completed`
- `failed`
- `cancelled`
- `timed_out`

The manager is planned to retain the handle, operation type, owning domain, start/end timestamps, deadline, status, bounded progress, byte counters, bounded error/result data, temporary resources, TTL, and cleanup status.

The manager owns generic lifecycle mechanics only. Domain services provide the operation and retain validation, business rules, result meaning, and domain ownership. The manager knows neither MCP tool types nor domain-specific response types. The MCP adapter may later map eligible internal jobs to negotiated MCP Tasks. Internal states may be more detailed than external Task states, so the adapter requires a defined mapping and a bounded fallback decision for clients without Tasks support.

Deadline enforcement remains server-controlled:

1. The server sets a deadline.
2. A worker receives a cancellable context.
3. The operation checks cancellation regularly.
4. A watchdog observes the deadline.
5. The server cancels work on expiry.
6. External worker processes are terminated when required.
7. Temporary resources are removed or marked incomplete in a controlled way.
8. Status becomes `timed_out`.
9. Results and diagnostics remain bounded.

The normal execution unit is a controlled Go goroutine using `context.Context`, bounded buffers and concurrency, streaming, and deterministic cleanup. A separate operating-system process is justified only for an external program, hard CPU/memory isolation, crash isolation, work that cannot be cancelled reliably in-process, another user/security context, or an external platform mechanism.

Planned resource controls include global and per-domain concurrency, queue length, runtime, result size, temporary-data limits, shutdown cleanup, and job-leak prevention.

## Managed processes and command execution

A planned Managed Process Registry gives each server-started process an opaque process handle. Handles are the primary identity for status, output, waiting, and stopping. PIDs may be diagnostic fields but cannot be the sole security identity because PID reuse must not associate a request with the wrong process.

Stopping is limited by default to server-managed processes. Controlling an external PID requires a separate functional capability such as `process.control.external` plus a high-risk policy classification and is outside standard profiles. Process command lines and environments must be filtered and redacted.

Command execution resolves configured executable IDs to allowed absolute program paths. Arguments are a separate array; a free shell string is not a standard parameter. Working directories stay within allowed roots. Environment propagation follows an allowlist or defined rules. stdout and stderr are separate and bounded, as are runtime and concurrency. A future `run_command` is a synchronous wrapper over the Managed Process Engine, not a second execution engine. Interactive shells remain disabled and require an interactive/high-risk policy decision.

## Capabilities, profiles, and named roots

Not every tool is permanently exposed. Tool profiles and capabilities determine registration, while server-side checks remain authoritative. MCP tool annotations are metadata, not authorization.

Illustrative functional capabilities:

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

Profiles combine functional capabilities and policy. Possible names include `safe-read`, `filesystem-write`, and `admin`. Risk classifications such as `high-risk`, `destructive`, and `interactive` are additional policy conditions, not universal permissions. Final capability, profile, and risk-classification names remain open for Sprint 3.50. The current read-only behavior is the first restricted-profile case.

The current single root remains technically compatible. The target supports multiple named roots based on authoritative FlashGate server configuration and explicit root IDs plus relative paths, so absolute host paths need not reach the model. Each root may define read/write access, size/result limits, file types, capability mapping, symlink/reparse policy, and process working-directory permission. Deprecated MCP Roots is not an architectural dependency. Optional legacy compatibility may be evaluated only for a supported client and never overrides server configuration.

## MCP compatibility boundary

The local core is protocol-version independent. The MCP adapter owns version and extension negotiation, version-specific wire types, JSON Schema validation, and mappings from internal operations to external contracts. The implemented revision remains `2025-11-25`; no 2026 protocol feature is implemented in Sprint 3.41.

Future asynchronous exposure will evaluate the negotiated MCP Tasks Extension `io.modelcontextprotocol/tasks`. The adapter must map detailed internal states to external Task states and decide bounded behavior for clients without Tasks support. All future input/output schemas are validated against JSON Schema 2020-12. Protocol upgrades require compatibility tests, conformance evaluation, and changelog documentation. See ADR-0013.

## Resource and token efficiency

Efficiency is a measurable quality goal. Server-side primary measurements include bytes, characters, calls, entries, results, duration, allocations, and memory. Benchmarks are planned for startup, idle RSS, peak memory, CPU, allocations, p50/p95 latency, scanned/read/written bytes, serialized response size, `tools/list` size, calls per reference task, and approximate schema/response tokens. Model-specific token estimates supplement but do not replace byte and resource measurements.

Deterministic work should run locally. A `copy_path`-style local operation is preferred over reading content through a model and sending it back for writing.

Planned efficiency mechanisms include cursor pagination, limits, filtering, sorting, field selection, line and byte ranges, head/tail reads, batch reads/stats/hashing, bounded directory trees, targeted and atomic/conditional writes, dry runs, bounded filesystem plans, and compact structured results. No free-form workflow or shell language replaces clearly defined tools.

## Repository and deployment model

One repository and one primary binary remain the accepted baseline. Separation first occurs through packages, interfaces, and capabilities. Separate binaries require a measurable advantage in security isolation, deployment, startup time, memory, maintainability, platform separation, or independent releases. Benchmarks and threat models must support such a split. Microservices or IPC are not introduced without concrete benefit.

## Open-source, modules/providers, and MCP extensions

FlashGate MCP is a general, vendor-neutral open-source project. The core must not require Voxtronic paths, company-specific tool names, proprietary dependencies, product permissions, organization secrets, internal URLs, or infrastructure values.

Future public standard, community, vendor, organization-internal, or Voxtronic-specific **FlashGate modules/providers** may be possible as optional local project extensions. Their origin never weakens the security model. Every provider must use shared capabilities, root policies, limits, path validation, process policies, secret redaction, audit events, and platform adapters.

Possible module/provider metadata includes name, version, vendor, declared capabilities, tools, configuration schema, security classification, platform requirements, and dependencies. No FlashGate module identifier syntax is selected.

No module/provider mechanism is selected or implemented in Sprint 3.41. Before the first external provider, a decision must choose among statically linked extension packages registered at build time, registered in-process providers, or isolated out-of-process providers over local IPC. A Go module may be a source/versioning form for statically linked packages; it is not a runtime loading or isolation model.

An **MCP protocol extension** is different: it is a negotiated addition to the MCP wire protocol and uses the official vendor-prefix/slash identifier contract, such as `io.modelcontextprotocol/tasks`. FlashGate module/provider contracts do not imply MCP protocol-extension support.

## Planned work

The authoritative sequence is maintained in [BACKLOG.md](../BACKLOG.md). Near-term work is:

1. Sprint 3.42: technical project rename.
2. Sprint 3.43: pre-1.0 filesystem tool contract cleanup.
3. Sprint 3.44: Codex read-only activation preparation.
4. Sprint 3.45 onward: benchmarks/contracts, Operations/Job Manager, efficient filesystem work, search, named roots/capabilities, process, execution, and system information.

## Deferred decisions

- exact package boundaries and stable internal interfaces
- public Go API stability
- separate binaries or worker processes beyond the stated gates
- FlashGate module/provider contract and runtime model
- MCP protocol versions and extensions beyond `2025-11-25`
- final profile and operation-tool names
- local search index after benchmarks
- interactive shell support
- privacy-sensitive network information
- post-1.0 versioning, deprecation, and migration policy

Accepted decisions are recorded in [docs/adr](adr/). Planned components in this document must not be interpreted as implemented features.
