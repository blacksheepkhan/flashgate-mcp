# FlashGate MCP Project Identity

## Name

The binding public project name is **FlashGate MCP**.

Tagline: **Fast, secure and local-first host operations for MCP.**

Working description: FlashGate MCP is a fast, secure, resource-efficient, local-first MCP server for controlled filesystem, process, and operating-system operations.

## Meaning

- Flash means low latency, efficient local processing, and compact responses.
- Gate means controlled access through policies, capabilities, roots, limits, redaction, and audit events.

FlashGate MCP is not a web-hosting service and not a remote-shell replacement. It is a controlled local boundary between MCP clients and explicitly enabled operating-system functions.

## Scope and non-goals

The current implemented scope is a root-confined filesystem MCP server. Accepted future domains include search, managed processes, controlled execution, and system information.

Unrestricted host access, free-form remote shell behavior, implicit network exposure, and bypasses around central policies are non-goals.

## Open-source direction

The project is intended to be a general, vendor-neutral open-source project. Its core has no mandatory Voxtronic-specific paths, tools, proprietary systems, permissions, secrets, or infrastructure values.

Public, community, vendor, organization-internal, and Voxtronic-specific **FlashGate modules/providers** may be considered later as optional local project extensions. Provider origin never changes the central security boundary. Sprint 3.41 defines no module/provider contract, identifier syntax, or runtime model.

An **MCP protocol extension** is a separate negotiated addition to the MCP wire protocol and follows the official vendor-prefix/slash identifier contract, for example `io.modelcontextprotocol/tasks`. FlashGate modules/providers do not automatically define or implement MCP protocol extensions.

## Technical-name transition

Public name from Sprint 3.41: FlashGate MCP.

Current technical identifiers remain `fileserver-mcp` until Sprint 3.42:

- repository and local directory
- Go module and import paths
- binary and release artifacts
- MCP server implementation name (`serverInfo.name`)
- scripts, workflows, examples, and machine-readable catalog

Planned Sprint 3.42 targets:

| Item | Target |
|---|---|
| Repository | `flashgate-mcp` |
| Binary | `flashgate-mcp` |
| MCP server implementation name (`serverInfo.name`) | `flashgate` |
| Go module | `github.com/blacksheepkhan/flashgate-mcp` |
| Short name | FlashGate |

No technical rename is part of Sprint 3.41. After Sprint 3.42, old identifiers should remain only where required for changelog history, migration guidance, and historical records.
