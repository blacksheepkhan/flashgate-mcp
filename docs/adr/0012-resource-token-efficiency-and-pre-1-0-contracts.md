# ADR-0012: Resource/Token Efficiency and Pre-1.0 Contracts

## Status

Accepted

## Context

FlashGate MCP aims to reduce local resource cost and model round trips. The current project has no production deployment or external compatibility contract, so preserving redundant early tool names would create avoidable long-term cost.

## Decision

Treat startup time, idle RSS, peak memory, CPU, allocations, p50/p95 latency, byte counters, response size, `tools/list` size, calls per reference workflow, and approximate schema/response tokens as measurable quality goals. Server-side bytes, characters, calls, entries, results, duration, and memory remain primary; token estimates are supplemental.

Prefer local deterministic operations, pagination, limits, filtering, sorting, field selection, ranges, head/tail reads, batch inspection, bounded trees, targeted edits, atomic/conditional writes, dry runs, bounded filesystem plans, and compact structured results. Do not introduce a free-form workflow or shell language.

Before 1.0, justified tool names, parameters, and results may change or be removed. Planned cleanup is `list_files` to `list_directory`, `stat_path` to `get_path_info`, `mkdir` to `create_directory`, removal of `exists_path` in favor of `get_path_info`, and removal of `rename_path` in favor of `move_path`. Other base tools remain, with safe contract enhancements planned.

`get_path_info` is planned to represent a missing path structurally with `exists: false`, while distinct normalized errors remain for access, root, path-type, and I/O failures.

## Rationale

Local server work avoids transferring content through the model. Early contract cleanup is cheaper and clearer than artificial deprecation for unpublished tools.

## Consequences

- Sprint 3.41 changes no schema or tool registration.
- Breaking changes require changelog, tests, tool docs, examples, and smoke-test updates.
- A versioning, deprecation, and migration strategy is required before stable release.
- Benchmarks and snapshot tests become contract gates.

## Security Impact

Bounded inputs/results and local operations reduce data exposure and denial-of-service risk. MCP annotations do not replace authorization.

## Implementation Guidance

Perform the tool cleanup in its dedicated sprint. Record baselines before optimization. Use machine-readable errors and schema snapshots, and test all changes through MCP `tools/call`.

## Deferred Decisions

- final normalized error-code set
- exact range, batch, cursor, and plan schemas
- post-1.0 compatibility and deprecation policy
- model-specific token estimation method

## Implementation Amendment - 2026-07-11

Sprint 3.43 implements the decided pre-1.0 baseline as exactly eight default tools in this order: `list_directory`, `read_file`, `get_path_info`, `write_file`, `create_directory`, `delete_path`, `copy_path`, and `move_path`. The read-only profile exposes the first three read tools. The five replaced or redundant names have no aliases.

Inputs are closed objects and are strictly decoded at runtime. `get_path_info` uses one metadata operation and represents genuine absence structurally; policy denials remain errors. `create_directory` reports actual creation state. `copy_path` remains file-only. `move_path` is the single same-volume Move/Rename contract with pre-replacement SameFile, effective-subtree, identity-revalidation, and overwrite safeguards. Stable machine-readable MCP error payloads, runtime output schemas, structured content, directory copy, and general compatibility infrastructure remain deferred.

## Implementation Amendment - 2026-07-12

Sprint 3.45a chooses MCP `CallToolResult` variant B for every successful call from all eight filesystem tools: one `TextContent` block contains compact deterministic JSON, and `structuredContent` contains the same domain object. Protocol conformance takes precedence over minimizing the duplicated payload. Reproducible benchmarks record response bytes, runtime, and allocations for the historical direct form, text-only form, and selected text-plus-structured form without setting CI budgets from a single machine.

Typed domain results and the filesystem core remain independent of MCP. Runtime `outputSchema` is a separate all-eight-tool gate. The existing safe JSON-RPC tool-error model is not partially migrated; BL-203 remains responsible for a complete normalized `isError=true` contract.
