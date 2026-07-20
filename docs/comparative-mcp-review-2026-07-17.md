# Comparative MCP Architecture Review - 2026-07-17

## Purpose

This review records which ideas from comparable filesystem, shell, and desktop MCP projects were considered for FlashGate and how they affect the Version 1.0 plan.

It is an architecture and planning input, not a claim that FlashGate currently implements the compared features or outperforms the other projects.

## Evaluation criteria

The comparison prioritizes:

1. response latency;
2. low model-token consumption;
3. low RAM and CPU consumption;
4. no interpreter runtime when avoidable;
5. bounded local deterministic work;
6. server-enforced security;
7. cross-platform Windows/Linux behavior;
8. maintainable native distribution.

## Projects considered

### Official MCP filesystem reference server

Useful patterns:

- line-head and line-tail reads;
- multi-file reads with per-item failures;
- dry-run edits and diffs;
- directory trees;
- MIME-aware media handling;
- complete tool annotations.

Not adopted as architecture:

- Node.js runtime dependency;
- client-provided MCP Roots as the authorization foundation;
- broad reference-example assumptions without FlashGate's threat model and resource budgets.

### Native Rust filesystem MCP

Useful patterns:

- native self-contained binary;
- read-only safe default;
- ability to disable tools and reduce catalog size;
- asynchronous/bounded I/O direction.

Adopted direction:

- safe read-only default after roots are configured;
- profile-specific tool exposure and catalog budgets;
- native interpreter-free release artifacts.

### Go filesystem MCP implementations

Useful patterns:

- embeddable/testable Go core;
- MCP resources for file content;
- MIME and binary detection;
- inline/base64 size limits;
- batch reads and metadata.

Adopted direction:

- explicit core/adapter separation;
- opaque large-result/resource abstraction;
- MIME/content classes and bounded transfer modes;
- batch inspection.

FlashGate does not promise a stable public Go API before Version 1.0 unless separately decided.

### Desktop Commander-style MCPs

Useful patterns:

- managed long-running processes;
- incremental output cursors;
- broad operational diagnostics and audit visibility;
- explicit awareness of large tool-catalog token cost.

Not adopted:

- unrestricted generic shell;
- large default tool catalog;
- assumption that application path/block lists secure an otherwise arbitrary shell;
- desktop-automation breadth in the FlashGate core.

Adopted direction:

- managed process registry;
- separate bounded stdout/stderr ring buffers;
- cursor-based output;
- tool-catalog measurement;
- audit correlation.

### Configured shell/command MCPs

Useful pattern:

- named, configured operations rather than arbitrary command strings.

Adopted with stricter rules:

- typed command definitions;
- server-resolved executable IDs and paths;
- fixed/allowed arguments and values;
- no shell string;
- environment, working-directory, timeout, output, and network policies;
- optional executable identity pinning.

Not adopted:

- dynamic shell templates as the primary contract;
- interpreter or policy-language dependency;
- arbitrary user-supplied command lines.

## FlashGate strengths retained

The comparison confirms the existing FlashGate direction:

- one native Go binary;
- low measured startup and memory footprint;
- compact profile-based catalogs;
- server-side path/root/capability enforcement;
- no general shell;
- bounded Operations/Job Manager;
- direct STDIO plus optional local system service;
- Windows/Linux adapters behind a protocol-independent core;
- benchmarked resource and token budgets.

## Gaps identified and planning result

| Gap | Version 1.0 decision |
|---|---|
| Heavy result duplicated in text and structured fields | Add payload classes and single-transmission rule |
| Missing useful-byte/wire-amplification metric | Add benchmark metrics and CI budget |
| Service execution identity unresolved | Adopt hybrid per-root architecture; implement Variant A only |
| Per-user service fairness insufficient | Add per-principal quotas and fair scheduling |
| Executable allowlist does not fully constrain arguments | Add typed command-definition schema |
| Tool visibility not explicitly token-budgeted per profile | Add catalog/instruction budgets and fingerprints |
| No compact server workflow instructions | Add bounded profile-specific instructions |
| Large/binary results lack a generic handoff model | Add opaque result/resource handles and fallback |
| Protocol planning did not fully account for stateless 2026 core | Add stateless-adapter, cache/TTL, and final Tasks planning |
| Audit model did not define operational lifecycle | Add correlation, rotation, retention, backpressure, and disk-full behavior |
| Native OS access preference was not a binding gate | Add native-adapter/no-interpreter policy |
| Release supply-chain evidence incomplete | Add checksums, signing plan, SBOM, provenance, and rollback |
| Performance advantage not compared reproducibly | Add pinned cross-project benchmark |
| Stable Version 1.0 boundary not explicit | Define Planned versus Later semantics and release gate |

## Deliberately deferred optimizations

The following ideas are valuable but not necessary for Version 1.0:

- conditional read/not-modified responses;
- per-user service workers;
- persistent user-scoped service hosts;
- ripgrep adapter;
- persistent search index;
- external PID control;
- interactive shell;
- provider/plugin ecosystem.

Their interfaces are kept feasible, but the initial release is not expanded to implement them without a demonstrated need.

## Non-goals confirmed

FlashGate will not become:

- a general remote shell;
- a broad desktop automation suite;
- a runtime package manager;
- an interpreter host;
- an unbounded file-indexing service;
- a remote/cloud service through this work package;
- a product with dozens of default tools irrespective of profile.

## Source baseline

The review used primary project and protocol documentation available on 2026-07-17. External project features change over time; the Version 1.0 cross-project benchmark must pin exact releases or commits and archive the compared configuration and results.

Relevant protocol planning references include:

- MCP `2025-11-25` Tools and Resources specifications;
- the `2026-07-28` MCP specification release candidate;
- SEP-2577 deprecation of Roots, Sampling, and Logging;
- the final Tasks Extension SEP;
- official and selected native filesystem MCP repositories.

## External references reviewed

Protocol and reference sources reviewed on 2026-07-17:

- [MCP Tools specification 2025-11-25](https://modelcontextprotocol.io/specification/2025-11-25/server/tools)
- [MCP Resources specification 2025-11-25](https://modelcontextprotocol.io/specification/2025-11-25/server/resources)
- [MCP schema reference 2025-11-25](https://modelcontextprotocol.io/specification/2025-11-25/schema)
- [MCP 2026-07-28 release candidate](https://blog.modelcontextprotocol.io/posts/2026-07-28-release-candidate/)
- [SEP-2577: Deprecate Roots, Sampling, and Logging](https://modelcontextprotocol.io/seps/2577-deprecate-roots-sampling-and-logging)
- [MCP Tasks Extension](https://tasks.extensions.modelcontextprotocol.io/seps/2663-tasks-extension)
- [Official MCP filesystem reference server](https://github.com/modelcontextprotocol/servers/tree/main/src/filesystem)
- [Rust MCP Filesystem](https://github.com/rust-mcp-stack/rust-mcp-filesystem)
- [Mark3labs Go MCP Filesystem Server](https://github.com/mark3labs/mcp-filesystem-server)
- [Desktop Commander MCP](https://github.com/wonderwhy-er/DesktopCommanderMCP)
- [MCPShell](https://github.com/inercia/MCPShell)

The comparison records architectural ideas, not a stable feature guarantee for external projects. The `BL-259` benchmark must pin exact releases or commits and preserve the evaluated configuration.

## Related documents

- [Efficiency improvement plan](efficiency-improvement-plan.md)
- [Version 1.0 scope](version-1-scope-and-release-boundary.md)
- [Architecture](architecture.md)
- [Authoritative backlog](../BACKLOG.md)

## Backlog ID correction - 2026-07-20

This dated comparative review retains the identifier valid when it was written.
The active reproducible cross-project efficiency benchmark task moved from historical `BL-259` to current `BL-261`.

The authoritative mapping is recorded in `docs/backlog-id-migration-2026-07-20.md`.
