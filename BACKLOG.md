# Backlog

This backlog tracks planned, open, and discovered work for `fileserver-mcp`.

The project goal is to provide a fast, maintainable, and secure MCP server for filesystem and local system operations. It is intended to cover the practical use cases of a filesystem MCP server and selected Desktop Commander style capabilities, while keeping security boundaries explicit and testable.

The backlog is maintained as part of the normal sprint workflow. New tasks discovered during implementation should be added here instead of being kept only in chat history.

## Working Rules

- Keep `README.md`, `CHANGELOG.md`, and `BACKLOG.md` current during each sprint.
- Add newly discovered work to this backlog before starting unrelated implementation.
- Prefer small, focused feature branches.
- Keep dangerous capabilities disabled or restricted by default.
- Do not write diagnostics to standard output during normal MCP operation.
- Use standard output only for JSON-RPC protocol messages.
- Use standard error for diagnostics, logging, and operational messages.

## Status Legend

| Status | Meaning |
|---|---|
| Ready | Clear and ready for implementation |
| Planned | Intended, but details still need refinement |
| Later | Useful, but not required soon |
| Blocked | Cannot continue without a decision or dependency |
| Done | Completed and retained for planning traceability |

## Current Sprint

### Sprint 3.35 - Read-only enforcement and tool capability gating

| ID | Status | Task | Reason | Notes |
|---|---|---|---|---|
| BL-049 | Done | Enforce `MCP_READ_ONLY` across filesystem tools | Config is parsed, but write-capable operations must not be available in read-only mode | Sprint 3.35 gates write-capable tools out of `tools/list` and direct `tools/call` by not registering them when `MCP_READ_ONLY=true` |

## Upcoming Sprints

| Sprint | Backlog IDs | Scope | Notes |
|---|---|---|---|
| Sprint 3.36 | BL-053, BL-054 | Root, realpath, and traversal hardening | Replace purely lexical checks with effective path validation |
| Sprint 3.37 | BL-050, BL-051, BL-052 | Hidden files, UNC paths, symlink/junction/reparse policy | Windows and Linux filesystem edge cases |
| Sprint 3.38 | BL-055 | JSON-RPC validation and error behavior | Validate request shape, IDs, methods, params, notifications, and invalid request handling |
| Sprint 3.39 | BL-008, BL-011, BL-034, BL-035 | Limits, logging, safe defaults, and secrets-aware behavior | Include read/list response limits and audit-oriented stderr logging |
| Sprint 3.40 | BL-001, BL-002, BL-040 | Windows/Linux test matrix and smoke tests | The earlier Linux JSON-RPC smoke-test CI item is deferred into Sprint 3.40 |
| Sprint 3.41 | BL-037, BL-038, BL-039 | Codex read-only activation preparation, without activation | Configuration examples, troubleshooting, and activation checklist |

## Epics

### Filesystem Tools

The filesystem toolset is the current foundation of the project.

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

Open work:

| ID | Status | Task | Reason | Notes |
|---|---|---|---|---|
| BL-004 | Planned | Add more JSON-RPC integration tests for filesystem tools | Current smoke test covers `initialize` and `tools/list` only | Add `tools/call` checks for read, write, list, stat, exists |
| BL-005 | Planned | Add benchmark tests | Measure large directory and file operation behavior | Useful before optimizing |
| BL-007 | Planned | Add append-file support | Common file operation not yet covered | Separate from overwrite-based `write_file` |
| BL-008 | Planned | Define larger-file streaming strategy | Current `read_file` is max-size limited | Avoid loading large files fully into memory |
| BL-009 | Planned | Add file metadata details | Improve file inspection usefulness | Include modified time and permissions where portable |
| BL-010 | Later | Add file hashing tool | Useful for diagnostics and comparisons | Needs size limits and algorithm choices |

### Search Tools

Search is a key capability for practical file-oriented MCP usage.

| ID | Status | Task | Reason | Notes |
|---|---|---|---|---|
| BL-011 | Planned | Design search tool model | Avoid unsafe or unbounded filesystem scans | Define path scope, recursion, limits, excludes |
| BL-012 | Planned | Add filename search tool | Common user workflow | Search by name or glob-like pattern |
| BL-013 | Planned | Add text search tool | Common code and document workflow | Needs max file size, encoding handling, binary detection |
| BL-014 | Later | Add ignore-file support | Avoid scanning build/vendor/cache folders | Consider `.gitignore` compatibility later |

### Process Tools

Process tools are part of the intended Desktop Commander replacement scope, but require stricter security controls than filesystem read operations.

| ID | Status | Task | Reason | Notes |
|---|---|---|---|---|
| BL-015 | Planned | Design process tool security model | Listing and controlling processes can expose sensitive information | Define defaults, permissions, and platform differences |
| BL-016 | Planned | Add `list_processes` tool | Desktop Commander style capability | Include PID, executable name, command line if allowed |
| BL-017 | Planned | Add process detail tool | Useful for diagnostics | Must avoid exposing sensitive environment data by default |
| BL-018 | Planned | Add stop process tool | Needed for process control workflows | Should require explicit capability and safe errors |
| BL-019 | Planned | Add start process tool | Needed for controlled local automation | Should not be a generic shell by default |
| BL-020 | Later | Add process tree view | Useful for diagnostics | Platform-specific behavior must be handled carefully |

### Command Execution

Command execution is powerful but dangerous. It should not become an unrestricted shell by default.

| ID | Status | Task | Reason | Notes |
|---|---|---|---|---|
| BL-021 | Planned | Design command execution security model | Required before implementing command tools | Prefer allowlisted commands, not arbitrary shell execution |
| BL-022 | Planned | Add allowlisted command execution | Safer alternative to free-form shell | Config-driven allowlist |
| BL-023 | Planned | Add command timeout handling | Prevent hanging processes | Include default and maximum timeout |
| BL-024 | Planned | Add command output limits | Prevent excessive memory use and huge MCP responses | Separate stdout and stderr limits |
| BL-025 | Planned | Define command result schema | Stable client behavior | Include exit code, stdout, stderr, timeout flag |
| BL-026 | Planned | Add working directory handling | Needed for practical command execution | Restrict to configured roots |
| BL-027 | Later | Add interactive shell/session support | Desktop Commander style capability | High risk; defer until core security model is mature |

### System Information

System information tools can help replace diagnostic Desktop Commander workflows.

| ID | Status | Task | Reason | Notes |
|---|---|---|---|---|
| BL-028 | Planned | Add `system_info` tool | Basic diagnostics | OS, architecture, hostname if allowed, server version |
| BL-029 | Planned | Add disk usage tool | Useful for filesystem diagnostics | Scope to configured roots where possible |
| BL-030 | Planned | Add environment inspection with filtering | Useful for diagnostics | Must redact or exclude sensitive values by default |
| BL-031 | Later | Add network information tool | Useful but privacy-sensitive | Requires explicit design |

### Security and Capability Model

The project must remain secure-by-default as functionality expands beyond file operations.

| ID | Status | Task | Reason | Notes |
|---|---|---|---|---|
| BL-032 | Planned | Define capability model for dangerous tools | Needed before process and command tools | Example capabilities: filesystem.write, process.control, command.execute |
| BL-033 | Planned | Add configuration flags for tool groups | Allow operators to disable risky features | Filesystem write, process tools, command tools |
| BL-034 | Planned | Add audit-oriented stderr logging | Improve traceability without corrupting STDIO | Must never write logs to stdout |
| BL-035 | Planned | Add safe defaults for non-developer users | Improve usability and safety | Dangerous tools disabled unless explicitly enabled |
| BL-036 | Later | Review workflow pinning strategy | Supply-chain hardening | Major tags are currently used; SHA pinning could be considered later |

### Client Compatibility

The server should be easy to use with common MCP clients.

| ID | Status | Task | Reason | Notes |
|---|---|---|---|---|
| BL-037 | Planned | Add Claude Desktop configuration example | User adoption | Include Windows and Linux examples |
| BL-038 | Planned | Add general MCP client configuration examples | Reduce setup friction | Keep examples clear and minimal |
| BL-039 | Later | Add troubleshooting guide | Support non-developer users | Include MCP_ROOT, JSON-RPC, and permission issues |

### CI, Release, and Quality

CI and release automation should keep validating real behavior, not only compilation.

| ID | Status | Task | Reason | Notes |
|---|---|---|---|---|
| BL-040 | Ready | Run Linux smoke test in Ubuntu CI | Validate real STDIO/JSON-RPC path on Linux | Use `scripts/smoke-jsonrpc.sh` |
| BL-041 | Planned | Add release notes or tag-based release workflow | Release maturity | Could produce GitHub Releases later |
| BL-042 | Planned | Add artifact verification step | Release quality | Verify `--version` on built artifacts |
| BL-043 | Later | Add cross-platform script validation | Prevent script regressions | PowerShell and Bash linting if useful |

### Documentation

Documentation is part of the project deliverable and should be updated continuously.

| ID | Status | Task | Reason | Notes |
|---|---|---|---|---|
| BL-044 | Ready | Reference backlog and changelog from README | Make planning documents discoverable | Add project planning section |
| BL-045 | Ready | Keep CHANGELOG updated per sprint | Preserve technical history | Update under `[Unreleased]` |
| BL-046 | Ready | Keep BACKLOG updated per sprint | Preserve discovered work | Move completed items to Done Summary when useful |
| BL-047 | Planned | Add non-developer smoke-test documentation | Support users without Go installed | Explain script-based validation |
| BL-048 | Planned | Expand tool documentation for future process and command tools | Needed before implementation | Document security implications clearly |

## Security Sprint Tracking

| ID | Status | Task | Reason | Notes |
|---|---|---|---|---|
| BL-049 | Done | Enforce `MCP_READ_ONLY` across filesystem tools | Config is parsed, but write-capable operations must not be available in read-only mode | Completed in Sprint 3.35 |
| BL-050 | Ready | Enforce `MCP_ALLOW_HIDDEN_FILES` in filesystem access | Config is parsed, but hidden-file policy may not be applied | Define portable Windows/Linux hidden-file behavior |
| BL-051 | Ready | Enforce `MCP_ALLOW_UNC_PATHS` on Windows roots and paths | Config is parsed, but UNC path policy may not be applied | Reject UNC roots and user paths unless explicitly allowed |
| BL-052 | Ready | Enforce `MCP_FOLLOW_SYMLINKS` consistently | Config is parsed, but filesystem operations may still follow symlinks | Use no-follow checks where needed before file operations |
| BL-053 | Ready | Close symlink escape risk from configured root | Symlinks inside the root may resolve outside the allowed tree | Validate resolved real paths before read, write, copy, move, rename, delete, and stat |
| BL-054 | Ready | Replace purely lexical path checks with real-path validation | Lexical root checks alone are insufficient for filesystem security | Keep lexical checks, then verify evaluated paths remain under root |
| BL-055 | Ready | Harden JSON-RPC request validation | Current validation is minimal | Validate JSON-RPC version, method shape, IDs, notifications, and invalid request handling |

## Done Summary

This section is intentionally not a full commit history. Detailed chronological changes are tracked in `CHANGELOG.md`.

| ID | Status | Task | Completed In | Notes |
|---|---|---|---|---|
| BL-D001 | Done | Establish Go project foundation | Early foundation commits | Module, package layout, configuration, security, filesystem abstraction |
| BL-D002 | Done | Implement root-confined filesystem tools | Multiple filesystem tool commits | list, read, stat, exists, write, mkdir, delete, move, copy, rename |
| BL-D003 | Done | Implement MCP server foundation | Multiple MCP commits | JSON-RPC, initialize, tools/list, tools/call, server loop |
| BL-D004 | Done | Add package and tool documentation | Documentation commits | README, package docs, tool catalog |
| BL-D005 | Done | Add CI pipeline | `c2ddf07` and follow-up commits | test, vet, build, lint |
| BL-D006 | Done | Add release build workflow | `ce50104` and follow-up commits | Windows and Linux release artifacts |
| BL-D007 | Done | Add version and help CLI modes | CLI/version commits | `--version`, `--help`, argument validation |
| BL-D008 | Done | Add Windows JSON-RPC smoke test script | `357f60f` | `scripts/smoke-jsonrpc.ps1` |
| BL-D009 | Done | Run JSON-RPC smoke test in Windows CI | `d8d5eec` | CI checks real STDIO/JSON-RPC path on Windows |
| BL-D010 | Done | Document JSON-RPC smoke test | `07816f8` | README updated |
| BL-D011 | Done | Add Linux/macOS JSON-RPC smoke test script | `8da6d0b` | `scripts/smoke-jsonrpc.sh` |
| BL-D012 | Done | Update GitHub Actions versions | `55418ed` | `checkout@v7`, `setup-go@v6` |
| BL-D013 | Done | Update upload-artifact action | `7540c67` | `upload-artifact@v6` |
| BL-D014 | Done | Add optional read-only mode | Sprint 3.35 | `MCP_READ_ONLY=true` disables write-capable filesystem tools at registration time |
| BL-D015 | Done | Add filesystem write capability gating | Sprint 3.35 | Read-only mode exposes only `list_files`, `read_file`, `stat_path`, and `exists_path` |
