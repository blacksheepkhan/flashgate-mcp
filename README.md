# FlashGate MCP

**Fast, secure and local-first host operations for MCP.**

FlashGate MCP is a resource-efficient cross-platform Model Context Protocol server for controlled filesystem, process, and operating-system operations. Deterministic work runs locally to minimize CPU, memory, latency, response size, model round trips, and token use.

> Sprint 3.42 completed the technical rename. FlashGate MCP uses repository
> `blacksheepkhan/flashgate-mcp`, module
> `github.com/blacksheepkhan/flashgate-mcp`, binary `flashgate-mcp`, and MCP
> server implementation name (`serverInfo.name`) `flashgate`.

It exposes secure filesystem operations to MCP-compatible clients through JSON-RPC over STDIO. The server is designed for predictable behavior, low operational overhead, clear security boundaries, and maintainable enterprise-style code.

## Status

The project currently implements the core MCP server loop, JSON-RPC routing and request validation, tool discovery, tool execution, filesystem abstraction, root-confined path handling, read-only tool gating, tests, and documentation.

The current implemented scope is filesystem operations. Search, process observation and management, allowlisted command execution, controlled system information, named roots, general capability profiles, and the Operations/Job Manager are accepted target architecture and remain planned work.

Implemented tools:

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

## Design Goals

- Secure-by-default filesystem access
- Root-confined path handling
- Deterministic MCP tool discovery
- Cross-platform support for Windows and Linux
- Single self-contained Go binary
- Minimal runtime dependencies
- Clear package boundaries
- Comprehensive unit tests
- Professional documentation for users and developers
- No external MCP framework dependency
- Measurable CPU, RAM, startup, latency, response-size, and token efficiency
- Local deterministic operations instead of transferring file content through the model

## FlashGate Security Gate

“Gate” means a server-enforced boundary: policies, capabilities, roots, limits, path and process validation, redaction, and audit events control access below the MCP adapter. Tool annotations and tool visibility are not authorization.

The current implementation enforces one configured root, read-only registration, path policies, hard limits, and redacted diagnostics. The target architecture plans multiple named roots and capability-based profiles while retaining server-side checks as authoritative. Dangerous capabilities remain disabled by default.

## Target Domains and Runtime

Accepted future domains are filesystem, search, process, execution, and system information. An optional shared Operations/Job Manager is planned for bounded long-running or managed work, cancellation, deadlines, progress, TTL, cleanup, and leak protection. Short synchronous work may remain directly in domain services; the manager is not currently implemented and does not own domain logic.

FlashGate MCP remains one repository and one primary binary unless benchmarks and threat models demonstrate a concrete isolation, deployment, performance, maintenance, platform, or release benefit from splitting it.

## Open-Source, Modules, and Protocol Extensions

FlashGate MCP is developed as a general, vendor-neutral open-source project. The core must not require Voxtronic paths, internal systems, proprietary dependencies, organization secrets, or company-specific permissions.

Public, community, vendor, organization-internal, and Voxtronic-specific FlashGate modules/providers may be considered later. No module/provider contract or runtime model is selected or implemented in Sprint 3.41, and future providers may not bypass central security controls.

MCP protocol extensions are separate negotiated wire-protocol features. The implemented protocol remains MCP `2025-11-25`; later features such as `io.modelcontextprotocol/tasks` inform compatibility planning but are not implemented. Deprecated MCP Roots is not the basis of FlashGate named roots.

## Protocol and Transport

`flashgate-mcp` communicates through MCP using JSON-RPC over STDIO.

Standard output is reserved for protocol messages. Diagnostic output must not be written to standard output because it would corrupt the JSON-RPC stream.

The server currently supports MCP methods such as:

```text
initialize
tools/list
tools/call
```

Filesystem operations are exposed as MCP tools and invoked through `tools/call`.

JSON-RPC request envelopes are validated before dispatch. Unsupported batch requests, invalid protocol versions, missing or invalid methods, invalid IDs, and malformed method params are rejected with generic JSON-RPC errors. Parse errors and invalid requests without a valid request ID serialize `id:null`. Notifications do not receive responses; `notifications/initialized` is accepted as a no-op, and other notifications are not executed.

JSON-RPC messages, tool arguments, filesystem payloads, and response sizes are bounded by configurable hard limits. A message that exceeds `MCP_MAX_JSONRPC_MESSAGE_BYTES` is rejected as an invalid request with `id:null` before tool dispatch. Oversized `tools/call` arguments return generic Invalid params errors, and filesystem limit violations return generic limit errors without host paths.

## Security Model

All filesystem operations are constrained to a configured root directory.

The root directory is configured through the `MCP_ROOT` environment variable:

```text
MCP_ROOT
```

`MCP_ROOT` is required. Production roots must be absolute, exist, satisfy the current root policy, and resolve to a directory. Missing, empty, whitespace-only, relative, non-existent, file, or policy-denied roots fail before tool registration and JSON-RPC processing. Expected configuration failures exit with code `3`, leave stdout empty, and report only a safe category on stderr.

The process working directory is available for development only when both `MCP_ROOT=.` and `MCP_ALLOW_CWD_ROOT=true` are set exactly. This opt-in never supplies a missing root and never enables other relative roots. It emits one safe warning on stderr.

Tool arguments use relative paths below this root. Absolute paths, path traversal outside the configured root, and unsafe path forms are rejected by the filesystem and security layers.

Path validation uses two stages:

1. Lexical validation rejects absolute user paths and leading parent traversal such as `..` or `../secret.txt`.
2. Effective path validation evaluates existing paths, and the nearest existing parent for create targets, to confirm the resulting filesystem location remains inside the configured root.

Individual tools do not bypass the filesystem abstraction and do not call host filesystem APIs directly. This keeps path validation centralized and testable.

Sprint 3.37 adds deny-by-default policy enforcement for hidden paths, UNC paths, symlinks, and Windows reparse points:

```text
MCP_ALLOW_HIDDEN_FILES=false
MCP_ALLOW_UNC_PATHS=false
MCP_FOLLOW_SYMLINKS=false
```

When hidden files are not allowed, dot-prefixed path components such as `.git/config` and Windows hidden-attribute paths are denied. `list_directory` filters hidden entries instead of failing the whole parent listing.

When UNC paths are not allowed, UNC roots and UNC-style user paths are rejected. When symlink following is not enabled, existing symlink components are denied and `list_directory` filters symlink or reparse entries. When symlink following is enabled, symlink targets are still constrained by effective root containment, and Windows junctions or non-symlink reparse points remain denied.

Security and path denials are mapped to generic invalid-path tool errors without exposing host absolute paths.

### Limits and diagnostics

Sprint 3.39 adds conservative hard limits:

| Environment variable | Default | Purpose |
|---|---:|---|
| `MCP_MAX_FILE_SIZE` | `10485760` | Maximum `read_file` bytes. Client `maxBytes` can only lower this cap. |
| `MCP_MAX_JSONRPC_MESSAGE_BYTES` | `16777216` | Maximum single JSON-RPC stdin message. |
| `MCP_MAX_TOOL_ARGUMENT_BYTES` | `12582912` | Maximum `tools/call` params or arguments payload. |
| `MCP_MAX_WRITE_BYTES` | `10485760` | Maximum `write_file` content bytes. |
| `MCP_MAX_LIST_ENTRIES` | `1000` | Maximum policy-visible entries returned by `list_directory`. |
| `MCP_MAX_COPY_BYTES` | `10485760` | Maximum `copy_path` source file size. |
| `MCP_MAX_DELETE_ENTRIES` | `1000` | Maximum entries allowed for recursive `delete_path`. |
| `MCP_MAX_RESPONSE_BYTES` | `16777216` | Safety net for serialized JSON-RPC responses. |

All limit values must be positive integers.

`MCP_DEBUG=true` enables minimal diagnostics on stderr. Diagnostics are redacted for common credentials, tokens, private-key markers, connection strings, and host paths. Normal MCP operation still writes only JSON-RPC protocol messages to stdout. No persistent logfiles are created.

### Read-only mode

Sprint 3.35 adds read-only enforcement for MCP tool discovery and direct tool calls.

Enable read-only mode with:

```text
MCP_READ_ONLY=true
```

When read-only mode is enabled, only these tools are registered and returned by `tools/list`:

```text
list_directory
read_file
get_path_info
```

Write-capable tools are not registered in read-only mode, so direct `tools/call` requests for `write_file`, `create_directory`, `delete_path`, `copy_path`, or `move_path` are rejected with a generic Invalid params error without revealing whether the tool exists in another mode.

## Tool Documentation

The MCP tool interface is documented in:

```text
docs/tools.md
```

Tool implementation and naming conventions are documented in:

```text
docs/tool-conventions.md
```

A machine-readable MCP tool catalog is available at:

```text
docs/mcp-tool-catalog.json
```

Preparation for a later, separately approved Codex read-only activation is documented in [docs/codex-read-only-activation.md](docs/codex-read-only-activation.md). Sprint 3.44 does not modify Codex configuration or register FlashGate as an MCP server.

The catalog contains tool names, descriptions, input schemas, result schemas, and common error behavior.

## Project Planning

Planned work is tracked in:

```text
BACKLOG.md
```

The backlog is the authoritative planning document for upcoming filesystem, process, command execution, system information, security, CI, release, and documentation work.

Project history is tracked in:

```text
CHANGELOG.md
```

`README.md`, `CHANGELOG.md`, and `BACKLOG.md` should be kept current as part of the normal sprint workflow.

Architecture and security references:

- [Architecture](docs/architecture.md)
- [Security model](docs/security.md)
- [Architecture decisions](docs/adr/)
- [Authoritative backlog](BACKLOG.md)
- [High-level roadmap](docs/roadmap.md)
- [Project identity](docs/project-identity.md)

## Project Structure

```text
cmd/
  server/               Application entry point

internal/
  config/               Runtime configuration
  diagnostics/          Redacted stderr diagnostics helpers
  fs/                   Filesystem abstraction and implementation
  mcp/                  MCP server, router, handlers, transport, and tools
  protocol/             JSON-RPC and MCP protocol types
  security/             Root-confined path validation
  version/              Build and release metadata

docs/
  mcp-tool-catalog.json Machine-readable MCP tool catalog
  tool-conventions.md   Tool implementation and interface conventions
  tools.md              Human-readable MCP tool reference
```

## Requirements

- Go 1.26 or newer

## Building

Clone the repository:

```bash
git clone https://github.com/blacksheepkhan/flashgate-mcp.git
cd flashgate-mcp
```

Build the server:

```bash
go build -o build/flashgate-mcp ./cmd/server
```

On Windows, the common build command is:

```powershell
go build -o build/flashgate-mcp.exe ./cmd/server
```

## Testing and Quality Checks

Run the standard validation chain:

```bash
go fmt ./...
go vet ./...
go test ./...
golangci-lint run
```

Build after validation:

```bash
go build -o build/flashgate-mcp ./cmd/server
```

On Windows:

```powershell
go build -o build/flashgate-mcp.exe ./cmd/server
```

Run the JSON-RPC smoke test on Windows:

```powershell
.\scripts\smoke-jsonrpc.ps1
```

Run the same smoke test in Windows read-only mode:

```powershell
$env:MCP_READ_ONLY = "true"
.\scripts\smoke-jsonrpc.ps1
Remove-Item Env:\MCP_READ_ONLY
```

On Linux, build the non-`.exe` binary and run:

```bash
bash scripts/smoke-jsonrpc.sh
MCP_READ_ONLY=true bash scripts/smoke-jsonrpc.sh
```

The default smoke test starts the built server binary and sends JSON-RPC requests through STDIO. It verifies that the server responds to:

```text
initialize
tools/list
```

The script also validates the negotiated protocol version and the registered MCP tool list.

Run the negative JSON-RPC smoke test on Windows:

```powershell
.\scripts\smoke-jsonrpc-negative.ps1
```

Run fail-closed startup validation:

```powershell
.\scripts\smoke-startup-negative.ps1
```

On Linux:

```bash
bash scripts/smoke-jsonrpc-negative.sh
```

Linux startup validation uses:

```bash
bash scripts/smoke-startup-negative.sh
```

The negative smoke test verifies malformed JSON, unknown methods, invalid `tools/call` params, and notification no-response behavior.

The smoke scripts create per-run JSONL request and response files under `build/` and remove them before exit. Script status output goes to the shell or CI log; the server process still writes only JSON-RPC protocol data to its redirected stdout stream.

## Release Builds

The repository contains a manual GitHub Actions workflow for release builds:

```text
.github/workflows/release-build.yml
```

The workflow can be started from GitHub:

```text
Actions → Release Build → Run workflow
```

The workflow accepts an optional version label. This value is embedded into the generated binaries and shown by the `--version` command.

Release builds currently produce the following artifacts:

```text
flashgate-mcp-linux-amd64
flashgate-mcp-windows-amd64
```

Artifacts are uploaded by GitHub Actions and can be downloaded from the completed workflow run.

Release binaries are built with:

```text
-trimpath
-ldflags="-s -w"
```

The release workflow also embeds build metadata through linker flags:

```text
version
commit
date
```

## Basic Usage

Set the filesystem root:

```powershell
$env:MCP_ROOT = "C:\Path\To\Allowed\Root"
```

The value must be an explicit absolute directory. For development only, `MCP_ROOT=.` additionally requires `MCP_ALLOW_CWD_ROOT=true`; client activation examples do not use that opt-in.

For read-only operation:

```powershell
$env:MCP_READ_ONLY = "true"
```

Run the server:

```powershell
.\build\flashgate-mcp.exe
```

The process expects JSON-RPC messages on standard input and writes JSON-RPC responses to standard output.

## CLI Help and Version Information

The binary supports a dedicated help mode:

```powershell
.\build\flashgate-mcp.exe --help
```

Short form:

```powershell
.\build\flashgate-mcp.exe -h
```

Example output:

```text
flashgate-mcp

Usage:
  flashgate-mcp
  flashgate-mcp --version
  flashgate-mcp --help

Environment:
  MCP_ROOT             Required absolute root directory exposed to MCP clients
  MCP_READ_ONLY        Set to true to expose only read-only filesystem tools
  MCP_ALLOW_CWD_ROOT   Development only: set to true with MCP_ROOT=.
```

The binary also supports a dedicated version mode:

```powershell
.\build\flashgate-mcp.exe --version
```

A local development build without embedded release metadata prints default values:

```text
flashgate-mcp
version: dev
commit: unknown
date: unknown
```

Release builds embed the version label, Git commit SHA, and UTC build date:

```text
flashgate-mcp
version: v0.1.0-test
commit: d9342ef1f1c4ebf03c2716f11d10b7fdb8dd316a
date: 2026-07-05T20:02:54Z
```

The `--version` mode is intended for diagnostics and artifact traceability. Normal MCP operation still communicates exclusively through JSON-RPC over STDIO.

## Example JSON-RPC Request

Example `tools/list` request:

```json
{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}
```

On Windows, a simple smoke test can be executed with a JSONL file:

```powershell
$json = '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}'
[System.IO.File]::WriteAllText("$PWD\request-tools-list.jsonl", $json + "`n", [System.Text.UTF8Encoding]::new($false))
cmd /c ".\build\flashgate-mcp.exe < request-tools-list.jsonl"
Remove-Item request-tools-list.jsonl
```

## Development Workflow

The project uses small, focused feature branches.

Typical workflow:

```bash
git checkout -b feature-sprint-name
go fmt ./...
go vet ./...
go test ./...
golangci-lint run
go build -o build/flashgate-mcp ./cmd/server
git add .
git commit -m "Describe focused change"
git checkout main
git merge feature-sprint-name
git branch -d feature-sprint-name
git push
```

Each feature should include:

- focused implementation
- unit tests
- documentation updates when interfaces change
- successful formatting, vetting, tests, linting, and build

## Implemented MCP Tools

| Tool | Description |
|---|---|
| `list_directory` | Lists files and directories below the configured filesystem root. |
| `read_file` | Reads a text file below the configured filesystem root. |
| `get_path_info` | Returns existence and metadata; missing paths return `exists:false`. |
| `write_file` | Writes a text file. |
| `create_directory` | Creates a directory and reports whether it was newly created. |
| `delete_path` | Deletes a file or directory. |
| `copy_path` | Copies a file. Directory copy is currently unsupported. |
| `move_path` | Moves or renames a file or directory on the same volume. |

When `MCP_READ_ONLY=true`, only `list_directory`, `read_file`, and `get_path_info` are exposed.

## Roadmap

Planned work is maintained authoritatively in [BACKLOG.md](BACKLOG.md). [docs/roadmap.md](docs/roadmap.md) contains only the high-level sequence.

The backlog covers filesystem tools, search tools, process tools, command execution, system information, security and capability controls, client compatibility, CI, release automation, and documentation.

## License

This project is licensed under the GNU General Public License v3.0.

See `LICENSE` for details.
