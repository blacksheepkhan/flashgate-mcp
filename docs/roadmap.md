# FlashGate MCP Roadmap

`BACKLOG.md` is the only authoritative planning and steering document. This roadmap provides only the high-level sequence and intentionally does not duplicate the task catalog.

## Current direction

FlashGate MCP is the binding project name. The repository, Go module, binary, and MCP server implementation name (`serverInfo.name`) remain technically named `fileserver-mcp` until the dedicated Sprint 3.42 rename.

The current implementation focuses on secure, root-confined filesystem tools. Filesystem contract and efficiency work comes before process and execution work because it exercises shared policies, limits, structured results, benchmarks, named roots, capabilities, and the Operations/Job Manager with a smaller risk surface.

## Sequence

| Phase | Sprints | Direction |
|---|---|---|
| Architecture and identity | 3.41 | FlashGate identity, ADR baseline, authoritative backlog consolidation |
| Technical transition | 3.42-3.44 | Technical rename, pre-1.0 filesystem contract cleanup, Codex read-only preparation |
| Shared foundations | 3.45-3.46 | Efficiency/contract baselines and Operations/Job Manager |
| Filesystem and search | 3.47-3.49 | Efficient inspection, safe edits/plans, bounded search |
| Policy model | 3.50 | Named roots, capabilities, and profiles |
| Process and execution | 3.51-3.54 | Threat model, observation, managed execution, allowlisted commands/isolation |
| System information | 3.55 | Scoped and redacted host information |

Operations/jobs are optional shared runtime infrastructure, not a separate business domain or mandatory layer. Short synchronous work may run directly in domain services. Long-running or managed filesystem, search, process, execution, or system work retains its domain ownership while optionally using the common lifecycle manager.

See [BACKLOG.md](../BACKLOG.md) for canonical tasks, sprint IDs, status, and scope. See [architecture.md](architecture.md) for current state, accepted target architecture, planned components, and deferred decisions.
