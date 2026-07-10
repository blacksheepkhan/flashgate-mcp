# ADR-0013: MCP Version and Extension Compatibility

## Status

Accepted

## Context

FlashGate MCP currently implements MCP protocol version `2025-11-25`. The MCP specification continues to evolve through protocol revisions and negotiated protocol extensions. Later specifications deprecate MCP Roots and define the official Tasks Extension, while JSON Schema 2020-12 is the default schema dialect. FlashGate needs a compatibility strategy that does not couple its local system core to one wire version or incorrectly claim support for later protocol features.

FlashGate modules/providers are optional local project extensions. MCP protocol extensions are separately specified, identified, and negotiated additions to the MCP wire protocol. They are not the same concept.

## Decision

The local FlashGate core remains independent of MCP protocol versions. The MCP adapter owns protocol-version negotiation, extension negotiation, version-specific DTOs, and mapping between internal domain results and MCP wire contracts.

The implemented protocol remains `2025-11-25`. Later MCP features are not considered supported until their adapter implementation, tests, compatibility review, and changelog entry are complete.

New MCP features must not be inserted directly into core domains. The internal Operations/Job Manager will be designed so the MCP adapter can later map eligible internal jobs to the official MCP Tasks Extension `io.modelcontextprotocol/tasks`. No custom operation status/result/cancel tools are accepted as the primary MCP job contract while the Tasks-extension and client-compatibility decision remains open.

Internal operation state may be more detailed than the external MCP Task state. The adapter therefore requires an explicit, tested mapping for status, result, cancellation, errors, TTL, and client-visible messages. A bounded fallback for clients without Tasks support must be decided before asynchronous MCP behavior is exposed.

MCP Roots is deprecated in the later specification line and is not a foundation of FlashGate named roots. Named roots are authoritative server configuration addressed by explicit root IDs and relative tool paths. Legacy MCP Roots behavior may be evaluated only as optional compatibility for a supported legacy client, with no architectural dependency.

MCP protocol extensions use the official vendor-prefix/slash identifier and capability-negotiation contract, for example `io.modelcontextprotocol/tasks`. FlashGate module/provider identifiers and runtime loading rules remain undecided and must not reuse MCP extension rules by implication.

All future tool input and output schemas will be validated against JSON Schema 2020-12. Any protocol upgrade requires protocol-version and extension negotiation tests, client compatibility review, schema/conformance checks, and changelog documentation.

## Rationale

Keeping version-specific protocol details at the adapter boundary preserves a stable local core and permits deliberate support for multiple client eras. Aligning future asynchronous behavior with the official Tasks Extension avoids creating a competing FlashGate-specific MCP contract. Server-configured named roots remain enforceable even when client support for deprecated MCP Roots is absent.

## Consequences

- MCP `2025-11-25` remains the only implemented protocol revision in Sprint 3.41.
- Final 2026 MCP features inform planning but are not claimed as supported.
- The Operations/Job Manager has no dependency on MCP types.
- The MCP adapter needs explicit internal-job-to-Task mapping and fallback decisions.
- Named roots have no dependency on MCP Roots.
- FlashGate modules/providers and MCP protocol extensions use separate terminology and contracts.
- JSON Schema 2020-12 becomes a validation requirement for future tool contracts.

## Security Impact

Extension negotiation is not authorization. Tasks and legacy compatibility paths must still enforce capabilities, root policies, caller/handle binding, limits, redaction, and audit controls. Opaque task or operation identifiers must not permit enumeration or cross-caller access. Deprecated MCP Roots must never weaken authoritative server root configuration.

## Implementation Guidance

- Keep version negotiation and extension declarations in the MCP adapter.
- Add version-specific compatibility tests before advertising a protocol revision or extension.
- Map internal states to official Task states deliberately; do not leak unrestricted internal diagnostics.
- Return synchronous results when supported and appropriate until Tasks compatibility is decided.
- For clients without a required extension, define a bounded synchronous or explicit capability-error fallback rather than an ad hoc job-tool surface.
- Validate all future input/output schemas as JSON Schema 2020-12 and evaluate official MCP conformance tooling.

## Decision Gates

- protocol versions supported in addition to `2025-11-25`
- Tasks Extension support and internal lifecycle mapping
- bounded behavior for clients without Tasks support
- optional legacy MCP Roots compatibility for a demonstrated client need
- FlashGate module/provider contract and runtime model

## Official References

- [MCP 2025-11-25 schema](https://modelcontextprotocol.io/specification/2025-11-25/schema)
- [SEP-1613: JSON Schema 2020-12](https://modelcontextprotocol.io/seps/1613-establish-json-schema-2020-12-as-default-dialect-f)
- [SEP-2133: Extensions](https://modelcontextprotocol.io/seps/2133-extensions)
- [SEP-2577: Deprecate Roots, Sampling, and Logging](https://modelcontextprotocol.io/seps/2577-deprecate-roots-sampling-and-logging)
- [SEP-2663: Tasks Extension](https://modelcontextprotocol.io/seps/2663-tasks-extension)
