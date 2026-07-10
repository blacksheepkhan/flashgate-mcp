# Changelog

All notable changes to this project will be documented in this file.

The format follows the spirit of [Keep a Changelog](https://keepachangelog.com/), and this project uses semantic versioning once releases begin.

## [Unreleased]

### Added

- FlashGate MCP project identity and transition documentation.
- Sprint 3.41 architecture baseline.
- Vendor-neutral open-source and FlashGate module/provider direction, separated from MCP protocol extensions.
- Accepted Operations and Job Manager target architecture.
- ADRs for project identity, domain-separated core, deployment/modules/providers, capability profiles, operations/jobs, managed processes/commands, resource/token efficiency, and MCP version/extension compatibility.
- Backlog ID migration documentation for the continuous `BL-001` sequence.
- Sprint 3.38 JSON-RPC request validation and error behavior hardening.
- Sprint 3.39 configurable hard limits, redacted diagnostics, and secrets-aware behavior.
- Sprint 3.40 Windows/Linux JSON-RPC smoke-test matrix coverage.
- Sprint 3.37 hidden, UNC, symlink, junction, and reparse policy enforcement.
- Sprint 3.36 root, realpath, and traversal hardening for filesystem access.
- Sprint 3.35 read-only tool capability gating for filesystem MCP tools.
- Initial Go module setup.
- Project structure for a professional MCP server implementation.
- Immutable configuration package.
- Environment-based configuration loading.
- Secure path validation with `PathGuard`.
- `SafePath` abstraction for validated filesystem paths.
- Local filesystem abstraction.
- Filesystem support for:
  - list
  - read
  - stat
  - exists
  - write
  - mkdir
  - delete
  - move
  - copy
  - rename
- MCP protocol types for JSON-RPC and MCP messages.
- MCP server loop for JSON-RPC over STDIO.
- MCP initialize handler.
- MCP router and handler abstraction.
- MCP tool registry.
- MCP `tools/list` support.
- MCP `tools/call` support.
- JSON-RPC envelope validation for protocol version, method shape, IDs, notifications, unsupported batches, and method-specific params.
- Configurable limits for JSON-RPC messages, tool arguments, write payloads, list entries, copy source size, recursive delete entries, and response size.
- Central diagnostics redaction for common tokens, credentials, private-key markers, connection strings, and host paths.
- Negative JSON-RPC smoke test coverage for malformed JSON, unknown methods, invalid `tools/call` params, and notification no-response behavior.
- Bash negative JSON-RPC smoke test script for Linux CI.
- Deterministic MCP tool discovery order.
- Filesystem MCP tools:
  - `list_files`
  - `read_file`
  - `stat_path`
  - `exists_path`
  - `write_file`
  - `mkdir`
  - `delete_path`
  - `move_path`
  - `copy_path`
  - `rename_path`
- Unit tests for configuration, security, filesystem, protocol, router, transport, server, initialize, tools, CLI, version, and bootstrap packages.
- Package documentation files.
- Human-readable MCP tool documentation in `docs/tools.md`.
- Tool implementation conventions in `docs/tool-conventions.md`.
- Machine-readable MCP tool catalog in `docs/mcp-tool-catalog.json`.
- PowerShell scripts for build, lint, and test workflows.
- Windows JSON-RPC smoke test script in `scripts/smoke-jsonrpc.ps1`.
- Linux/macOS JSON-RPC smoke test script in `scripts/smoke-jsonrpc.sh`.
- GitHub Actions CI workflow for formatting, vetting, tests, linting, and build validation.
- Windows and Ubuntu CI JSON-RPC smoke test execution for default, read-only, and negative smoke variants.
- Manual GitHub Actions release build workflow.
- Release build artifacts for Windows and Linux.
- Release artifact summary and retention configuration.
- Version metadata package.
- Build metadata embedding through linker flags.
- `--version` CLI mode.
- `--help` and `-h` CLI modes.
- CLI argument validation with dedicated invalid-argument exit behavior.
- GNU General Public License v3.0 license file.
- Project backlog in `BACKLOG.md`.

### Changed

- Public project name changed from Fileserver MCP to FlashGate MCP; the technical repository, module, binary, MCP server implementation name (`serverInfo.name`), scripts, workflows, and catalog remain unchanged until Sprint 3.42.
- Backlog consolidated into one canonical continuously numbered task catalog without a separate `BL-D` series.
- Former Sprint 3.41 Codex read-only preparation shifted to Sprint 3.44.
- Architecture expanded toward a domain-separated local system core while clearly distinguishing current and planned components.
- Long-running work is modeled through a shared planned Operations/Job layer without changing domain ownership.
- Pre-1.0 tool contracts are explicitly allowed to change before a stable external contract exists.
- Vendor-neutral open-source core and optional FlashGate module/provider direction documented without selecting a runtime model; MCP protocol extensions are treated separately.
- Sprint 3.41 review corrections align planning with SEP-1613, SEP-2133, SEP-2577, SEP-2663, and MCP `2025-11-25` implementation-name semantics.
- Roadmap clarified as a high-level pointer to authoritative `BACKLOG.md` planning.
- Unknown tool names in `tools/call`, including read-only-gated write tools, now return generic JSON-RPC Invalid params errors instead of Method not found.
- `MCP_MAX_FILE_SIZE` is now a hard server cap for `read_file`; client `maxBytes` can reduce but not increase it.
- Minimal debug diagnostics are now gated by `MCP_DEBUG` and written only to stderr after redaction.
- Smoke scripts now use per-run JSONL request/response files under `build/` and clean them up before exit.
- JSON-RPC protocol errors now use generic messages such as `parse error`, `invalid request`, `invalid params`, `method not found`, and `internal error`.
- `PathGuard` now accepts an explicit filesystem security policy while keeping the default constructor compatible.
- `list_files` now filters hidden and denied link/reparse entries according to policy.
- Moved protocol definitions from `pkg/protocol` to `internal/protocol`.
- Removed the public `pkg` package layout in favor of internal packages.
- Refactored filesystem operations into focused files.
- Updated supported MCP protocol version from `2025-06-18` to `2025-11-25`.
- Split CLI and server bootstrap responsibilities.
- Extracted tool registry bootstrap.
- Extracted router bootstrap.
- Updated README with build, usage, CLI, release, tool, and smoke-test information.
- Updated roadmap handling so planned work is tracked in `BACKLOG.md`.
- Updated GitHub Actions versions to Node-24-compatible major versions:
  - `actions/checkout@v7`
  - `actions/setup-go@v6`
  - `actions/upload-artifact@v6`

### Fixed

- Stabilized golangci-lint execution in CI by installing the expected linter version.
- Removed Node.js 20 deprecation annotations from CI and release workflows.
- Improved release artifact visibility in GitHub Actions.
- Validated JSON-RPC smoke-test behavior locally and in CI.

### Security

- Sprint 3.38 validates JSON-RPC envelopes before dispatch, rejects unsupported batches, suppresses responses for notifications, prevents `tools/call` notification execution, serializes unknown IDs as `id:null`, and converts handler panics to generic Internal error responses.
- Sprint 3.39 bounds JSON-RPC messages, tool arguments, filesystem operation payloads, recursive delete scope, and serialized responses with generic limit errors.
- Sprint 3.39 adds centralized redaction before debug diagnostics reach stderr.
- Sprint 3.36 adds effective path validation through `PathGuard` using evaluated existing paths and evaluated nearest existing parents for create targets.
- Sprint 3.36 rejects symlink-based filesystem escapes that resolve outside the configured root.
- Sprint 3.36 maps security/path denials to generic invalid-params tool errors without exposing host paths.
- Sprint 3.35 enforces `MCP_READ_ONLY=true` at tool registration time by exposing only `list_files`, `read_file`, `stat_path`, and `exists_path`.
- Sprint 3.35 prevents direct `tools/call` execution of write-capable tools in read-only mode because those tools are not registered.
- Filesystem paths are resolved through a sandbox root.
- Absolute user paths are rejected.
- Parent directory traversal is rejected.
- Destructive operations use conservative defaults.
- Tool implementations do not bypass the filesystem abstraction.
- Normal MCP operation reserves standard output for JSON-RPC protocol messages.
- Diagnostic output is kept separate from the JSON-RPC stream.
