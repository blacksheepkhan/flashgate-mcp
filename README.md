# fileserver-mcp

`fileserver-mcp` is a cross-platform Model Context Protocol (MCP) server written in Go.

It exposes secure filesystem operations to MCP-compatible clients through JSON-RPC over STDIO. The server is designed for predictable behavior, low operational overhead, clear security boundaries, and maintainable enterprise-style code.

## Status

The project currently implements the core MCP server loop, JSON-RPC routing, tool discovery, tool execution, filesystem abstraction, root-confined path handling, tests, and documentation.

Implemented tools:

```text
list_files
read_file
stat_path
exists_path
write_file
mkdir
delete_path
move_path
copy_path
rename_path
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

## Protocol and Transport

`fileserver-mcp` communicates through MCP using JSON-RPC over STDIO.

Standard output is reserved for protocol messages. Diagnostic output must not be written to standard output because it would corrupt the JSON-RPC stream.

The server currently supports MCP methods such as:

```text
initialize
tools/list
tools/call
```

Filesystem operations are exposed as MCP tools and invoked through `tools/call`.

## Security Model

All filesystem operations are constrained to a configured root directory.

The root directory is configured through the `MCP_ROOT` environment variable:

```text
MCP_ROOT
```

Tool arguments use relative paths below this root. Absolute paths, path traversal outside the configured root, and unsafe path forms are rejected by the filesystem and security layers.

Individual tools do not bypass the filesystem abstraction and do not call host filesystem APIs directly. This keeps path validation centralized and testable.

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

The catalog contains tool names, descriptions, input schemas, result schemas, and common error behavior.

## Project Structure

```text
cmd/
  server/               Application entry point

internal/
  config/               Runtime configuration
  fs/                   Filesystem abstraction and implementation
  mcp/                  MCP server, router, handlers, transport, and tools
  protocol/             JSON-RPC and MCP protocol types
  security/             Root-confined path validation

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
git clone https://github.com/blacksheepkhan/fileserver-mcp.git
cd fileserver-mcp
```

Build the server:

```bash
go build -o build/fileserver-mcp ./cmd/server
```

On Windows, the common build command is:

```powershell
go build -o build/fileserver-mcp.exe ./cmd/server
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
go build -o build/fileserver-mcp ./cmd/server
```

On Windows:

```powershell
go build -o build/fileserver-mcp.exe ./cmd/server
```

## Basic Usage

Set the filesystem root:

```powershell
$env:MCP_ROOT = "C:\Path\To\Allowed\Root"
```

Run the server:

```powershell
.\build\fileserver-mcp.exe
```

The process expects JSON-RPC messages on standard input and writes JSON-RPC responses to standard output.

## Example JSON-RPC Request

Example `tools/list` request:

```json
{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}
```

On Windows, a simple smoke test can be executed with a JSONL file:

```powershell
$json = '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}'
[System.IO.File]::WriteAllText("$PWD\request-tools-list.jsonl", $json + "`n", [System.Text.UTF8Encoding]::new($false))
cmd /c ".\build\fileserver-mcp.exe < request-tools-list.jsonl"
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
go build -o build/fileserver-mcp ./cmd/server
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
| `list_files` | Lists files and directories below the configured filesystem root. |
| `read_file` | Reads a text file below the configured filesystem root. |
| `stat_path` | Returns metadata for a file or directory. |
| `exists_path` | Checks whether a file or directory exists. |
| `write_file` | Writes a text file. |
| `mkdir` | Creates a directory. |
| `delete_path` | Deletes a file or directory. |
| `move_path` | Moves a file or directory. |
| `copy_path` | Copies a file or directory. |
| `rename_path` | Renames a file or directory. |

## Roadmap

Planned hardening and extension areas:

- integration tests
- CI pipeline
- release packaging
- benchmark tests
- optional read-only mode
- search tools
- append-file support
- larger-file streaming strategy
- structured logging to standard error
- client configuration examples

## License

This project is licensed under the GNU General Public License v3.0.

See `LICENSE` for details.
