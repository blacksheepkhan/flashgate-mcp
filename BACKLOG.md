# FlashGate MCP Backlog

This is the authoritative planning and steering document for FlashGate MCP.

The public project name is FlashGate MCP from Sprint 3.41. The repository, local directory, Go module, binary, MCP server implementation name (`serverInfo.name`), scripts, workflows, and machine-readable catalog still use the technical identifier `fileserver-mcp` until Sprint 3.42. No technical rename is performed in Sprint 3.41.

## Working rules

- Keep `README.md`, `CHANGELOG.md`, and `BACKLOG.md` current in every sprint.
- Define every canonical task exactly once in the task catalog; sprint tables reference IDs only.
- Keep dangerous capabilities disabled or restricted by default.
- Keep protocol output on stdout and diagnostics on stderr.
- Mark target and planned behavior explicitly; do not present it as implemented.
- Preserve completed tasks with status `Done`.
- Complete the current sprint's ID migration document before merge. After merge, dated migration files are immutable history; any later full renumbering creates a new dated migration file and may add a small migration index.

> Neue Aufgaben werden an der fachlich richtigen Position eingefügt. Danach werden alle nachfolgenden BL-IDs repositoryweit fortlaufend umnummeriert.

## Status legend

| Status | Meaning |
|---|---|
| Ready | Clear and ready for implementation |
| Planned | Accepted work that still needs implementation or refinement |
| Later | Accepted but intentionally deferred |
| Blocked | Waiting for a decision or dependency |
| Done | Completed and retained for traceability |

## Completed sprint baseline

### Sprint 3.40 - Windows/Linux test matrix and smoke tests

Backlog IDs: `BL-023`, `BL-024`, `BL-025`.

Completed sprint numbers through Sprint 3.40 are historical and are not renumbered. Task IDs were migrated into the continuous catalog below.

## Planned sprint sequence

| Sprint | Backlog IDs | Scope |
|---|---|---|
| Sprint 3.41 | BL-026–BL-035 | FlashGate architecture baseline and backlog consolidation |
| Sprint 3.42 | BL-227–BL-243 | Technical project rename to FlashGate MCP |
| Sprint 3.43 | BL-244–BL-257 | Pre-1.0 filesystem tool contract cleanup |
| Sprint 3.44 | BL-174, BL-258–BL-266 | Codex read-only activation preparation |
| Sprint 3.45 | BL-189–BL-212 | Efficiency benchmarks and MCP tool contract foundation |
| Sprint 3.46 | BL-084–BL-099 | Operations and Job Manager foundation |
| Sprint 3.47 | BL-036–BL-049 | Efficient filesystem listing, reading and batch inspection |
| Sprint 3.48 | BL-050–BL-061, BL-063–BL-067 | Targeted edits, conditional writes and bounded filesystem plans |
| Sprint 3.49 | BL-068–BL-080, BL-082 | Filesystem and text search |
| Sprint 3.50 | BL-100–BL-111 | Named roots, capabilities and tool profiles |
| Sprint 3.51 | BL-113, BL-129, BL-132, BL-162, BL-165 | Process architecture and security model |
| Sprint 3.52 | BL-114–BL-118 | Process observation |
| Sprint 3.53 | BL-119–BL-126, BL-130–BL-131, BL-133–BL-135 | Managed process execution |
| Sprint 3.54 | BL-136–BL-149, BL-151–BL-152 | Allowlisted command execution and OS isolation |
| Sprint 3.55 | BL-153–BL-157 | System information |

Sprint 3.44 replaces the former Sprint 3.41 Codex preparation plan and must use the FlashGate technical names created in Sprint 3.42 and the cleaned tool names created in Sprint 3.43.

## Canonical task catalog

### Completed foundation and current implementation

| ID | Status | Task | Scope and acceptance notes |
|---|---|---|---|
| BL-001 | Done | Establish Go project foundation | Module, package layout, configuration, security, and filesystem abstraction |
| BL-002 | Done | Implement root-confined filesystem tools | Current list/read/stat/exists/write/mkdir/delete/move/copy/rename implementation |
| BL-003 | Done | Implement MCP server foundation | JSON-RPC, initialize, `tools/list`, `tools/call`, routing, and server loop |
| BL-004 | Done | Add package and tool documentation | README, package docs, human and machine-readable tool references |
| BL-005 | Done | Add CI pipeline | Formatting, vet, tests, lint, and build validation |
| BL-006 | Done | Add release build workflow | Windows and Linux artifacts under current technical names |
| BL-007 | Done | Add version and help CLI modes | `--version`, `--help`, and argument validation |
| BL-008 | Done | Add Windows JSON-RPC smoke script | Real STDIO path on Windows |
| BL-009 | Done | Run Windows JSON-RPC smoke in CI | Windows CI integration |
| BL-010 | Done | Document JSON-RPC smoke testing | README usage and validation description |
| BL-011 | Done | Add Linux/macOS JSON-RPC smoke script | Bash STDIO validation |
| BL-012 | Done | Update GitHub Actions major versions | Node-24-compatible action versions |
| BL-013 | Done | Update artifact upload action | Current artifact workflow version |
| BL-014 | Done | Add optional read-only mode | `MCP_READ_ONLY=true` restricted registration |
| BL-015 | Done | Enforce filesystem write capability gating | Read-only mode exposes only current read tools |
| BL-016 | Done | Harden effective-root and traversal validation | Real-path validation for existing paths and create parents |
| BL-017 | Done | Enforce hidden, UNC, symlink, junction, and reparse policy | Deny-by-default cross-platform path policy |
| BL-018 | Done | Harden JSON-RPC validation and error behavior | Envelopes, IDs, notifications, batches, params, unknown tools, and panic boundary |
| BL-019 | Done | Enforce protocol message and tool-argument limits | Bounded JSON-RPC and `tools/call` input |
| BL-020 | Done | Enforce filesystem operation and response limits | Bounded reads, writes, listing, copy, recursive delete, and responses |
| BL-021 | Done | Add centralized redaction and safe stderr diagnostics | Secret/host-path redaction, debug gating, and no stdout diagnostics |
| BL-022 | Done | Add safe defaults for non-developer users | Deny-by-default limits and conservative behavior |
| BL-023 | Done | Run Linux smoke test in Ubuntu CI | Sprint 3.40 real STDIO validation |
| BL-024 | Done | Run JSON-RPC smoke matrix on Windows and Linux | Default, read-only, and negative variants |
| BL-025 | Done | Isolate smoke JSONL artifacts per run | Unique files and deterministic cleanup |

### Project identity and architecture baseline

| ID | Status | Task | Scope and acceptance notes |
|---|---|---|---|
| BL-026 | Done | Adopt FlashGate MCP public identity | Name, tagline, meaning, scope, non-goals, and technical-name transition |
| BL-027 | Done | Document current and accepted target architecture | Diagram and explicit current/planned/deferred separation |
| BL-028 | Done | Define domain-separated local system core | Filesystem, search, process, execution, system, jobs, policy, limits, diagnostics, adapters, MCP |
| BL-029 | Done | Define core reuse and deployment baseline | Direct Go reuse, one repository/binary, evidence gates for IPC or split |
| BL-030 | Done | Define vendor-neutral open-source core | No mandatory Voxtronic assumptions, secrets, paths, or proprietary dependencies |
| BL-031 | Done | Define FlashGate module/provider direction and decision gate | Shared controls, no identifier or runtime model yet; distinct from MCP protocol extensions |
| BL-032 | Done | Define Operations/Job Manager architecture | Handles, states, deadlines, cancellation, resources, cleanup, goroutine/process gates |
| BL-033 | Done | Define capability profiles and named-root direction | Server-side enforcement, profile examples, root policy model |
| BL-034 | Done | Define managed process and execution architecture | Handles, PID rules, allowlists, no-shell default, single engine |
| BL-035 | Done | Define efficiency and pre-1.0 contract policy | Metrics, local work, planned cleanup, no artificial compatibility |

### Filesystem epic

| ID | Status | Task | Scope and acceptance notes |
|---|---|---|---|
| BL-036 | Planned | Add filesystem tests through MCP `tools/call` | Read, write, list, info, missing path, and security cases |
| BL-037 | Planned | Add paginated `list_directory` | Bounded pages rather than fail-or-return-all behavior |
| BL-038 | Planned | Define stable cursor semantics | Opaque cursor, invalidation, ordering, and policy changes |
| BL-039 | Planned | Add listing sort and filters | Portable name/type/metadata behavior |
| BL-040 | Planned | Add listing field selection | Return only requested portable fields |
| BL-041 | Planned | Add line-range reads | Bounded text line windows |
| BL-042 | Planned | Add byte-range and head/tail reads | Bounded binary-safe offsets and edge semantics |
| BL-043 | Planned | Expand portable path metadata | Modified time, permissions where portable, type, and size |
| BL-044 | Planned | Add large-file streaming strategy | Avoid whole-file memory loading and unbounded responses |
| BL-045 | Planned | Define binary-file read behavior | Detection, encoding, explicit modes, and limits |
| BL-046 | Planned | Add batch `read_files` | Per-item bounded results and partial-failure model |
| BL-047 | Planned | Add batch `get_paths_info` | Reduce repeated stat/existence round trips |
| BL-048 | Planned | Add batch hashing | Bounded algorithms, byte accounting, and job handoff |
| BL-049 | Planned | Add bounded directory tree | Depth, entry, byte, field, and pagination controls |
| BL-050 | Planned | Add exact targeted file changes | Explicit ranges or match-based edits without model retransmission |
| BL-051 | Planned | Add expected-match-count checks | Reject ambiguous or stale targeted edits |
| BL-052 | Planned | Add atomic writes | Same-filesystem replace and deterministic cleanup semantics |
| BL-053 | Planned | Add conditional write preconditions | Hash, modified-time, and path-type checks |
| BL-054 | Planned | Add dry-run support | Structured preview for destructive or multi-step changes |
| BL-055 | Planned | Add append-file support | Explicit bounded append distinct from overwrite |
| BL-056 | Planned | Add bounded filesystem plans | Limited known operations; no free-form workflow language |
| BL-057 | Planned | Enforce plan operation, entry, and byte limits | Preflight and runtime accounting |
| BL-058 | Planned | Define cross-volume move behavior | Copy/verify/delete, cancellation, and partial-state rules |
| BL-059 | Planned | Define conflict strategy | Fail, skip, replace, and explicit per-operation rules |
| BL-060 | Planned | Support directory copy and move | Bounded traversal, jobs, conflicts, and cleanup |
| BL-061 | Planned | Add directory-size operation | Streaming scan, limits, progress, and jobs |
| BL-062 | Planned | Add scoped disk-usage operation | Root-scoped capacity and privacy-safe results |
| BL-063 | Planned | Extend `write_file` safe modes | Create-only, replace-only, and explicit atomic/conditional behavior |
| BL-064 | Planned | Integrate long filesystem work with jobs | Preserve filesystem domain ownership |
| BL-065 | Planned | Threat-model bounded filesystem plans | Path races, rollback limits, conflicts, and partial completion |
| BL-066 | Planned | Add Windows/Linux filesystem MCP integration tests | Exercise real `tools/call` contracts and path behavior |
| BL-067 | Planned | Create representative filesystem benchmark corpus | Small/large files, deep/wide trees, binary data, and cross-volume cases |

### Search epic

| ID | Status | Task | Scope and acceptance notes |
|---|---|---|---|
| BL-068 | Planned | Define search model and threat model | Scope, recursion, data exposure, resource budgets, and errors |
| BL-069 | Planned | Add path search | Relative root-scoped paths |
| BL-070 | Planned | Add filename search | Literal and pattern matching |
| BL-071 | Planned | Add metadata filters | Type, size, and time with portable semantics |
| BL-072 | Planned | Add literal text search | Bounded content scanning |
| BL-073 | Planned | Add regular-expression search | Complexity and scan limits |
| BL-074 | Planned | Add include/exclude patterns | Deterministic precedence and relative matching |
| BL-075 | Planned | Enforce search depth, file, and scanned-byte limits | Server-side hard caps and counters |
| BL-076 | Planned | Enforce total and per-file match limits | Bounded results and diagnostics |
| BL-077 | Planned | Add bounded context lines | Per-match and aggregate response limits |
| BL-078 | Planned | Define binary detection and encoding behavior | Skip/error/explicit modes and reporting |
| BL-079 | Planned | Add search pagination | Opaque cursors and stable ordering |
| BL-080 | Planned | Add ignore-file support | Explicit policy and optional gitignore-compatible behavior |
| BL-081 | Later | Evaluate optional ripgrep adapter | Allowlisted local adapter with version/security checks |
| BL-082 | Planned | Provide pure-Go search fallback | Portable baseline without external dependency |
| BL-083 | Later | Decide on local search index after benchmarks | No index without measured need and privacy/lifecycle design |

### Operations and Job Manager epic

| ID | Status | Task | Scope and acceptance notes |
|---|---|---|---|
| BL-084 | Planned | Implement Operation Registry | Thread-safe ownership and lifecycle |
| BL-085 | Planned | Generate opaque operation handles | `op_<opaque-id>` without guessable internals |
| BL-086 | Planned | Implement operation status model | queued/running/completed/failed/cancelled/timed_out |
| BL-087 | Planned | Add context cancellation | Cooperative cancellation contract |
| BL-088 | Planned | Add server deadlines and watchdog | Server-controlled timeout enforcement |
| BL-089 | Planned | Define progress and byte counters | Read/written/scanned bytes and bounded domain progress |
| BL-090 | Planned | Add bounded result storage and TTL | Expiry and retrieval behavior |
| BL-091 | Planned | Add temporary-resource cleanup | Success, failure, cancellation, timeout, and incomplete markers |
| BL-092 | Planned | Limit global and per-domain parallel jobs | Configurable safe defaults |
| BL-093 | Planned | Limit operation queue length | Deterministic overload response |
| BL-094 | Planned | Define controlled server shutdown | Cancellation, grace period, worker termination, final state |
| BL-095 | Planned | Prevent and detect job leaks | TTL sweep, ownership checks, and metrics |
| BL-096 | Planned | Preserve domain ownership | Jobs execute work without becoming its business domain |
| BL-097 | Planned | Implement goroutine/subprocess decision rules | External, isolation, cancellation, identity, and platform gates |
| BL-098 | Planned | Add Windows and Linux job tests | Deadline, cancellation, cleanup, shutdown, and temporary files |
| BL-099 | Planned | Add job security tests and race detector | Handles, limits, lifecycle races, and leak checks |

### Named roots, capabilities, and profiles epic

| ID | Status | Task | Scope and acceptance notes |
|---|---|---|---|
| BL-100 | Planned | Implement functional capability model | Functional rights kept separate from profiles and risk classifications |
| BL-101 | Planned | Support multiple named roots | Preserve compatible single-root migration path |
| BL-102 | Planned | Use root ID and relative path in target contracts | Avoid model-visible absolute host paths |
| BL-103 | Planned | Implement profile and risk-policy configuration | Final profile names, composition, risk classifications, validation, and safe defaults |
| BL-104 | Planned | Add per-root read/write policy | Independent permissions per root |
| BL-105 | Planned | Add per-root size and result limits | File, scan, response, and temporary data policies |
| BL-106 | Planned | Add per-root allowed file types | Explicit portable matching |
| BL-107 | Planned | Add per-root symlink/reparse rules | Preserve root confinement and platform semantics |
| BL-108 | Planned | Add per-root capability mapping | Tool and operation authorization |
| BL-109 | Planned | Add process working-directory permission per root | Execution policy integration |
| BL-110 | Planned | Implement dynamic tool registration | Effective profile/capability controls `tools/list` |
| BL-111 | Planned | Add negative capability and `tools/list` tests | Hidden/unavailable tools and direct-call denial |
| BL-112 | Later | Evaluate deprecated MCP Roots only for legacy client compatibility | No architectural dependency; server configuration and explicit root IDs remain authoritative; implement only for a supported legacy client; document deprecation and protocol-version scope |

### Process epic

| ID | Status | Task | Scope and acceptance notes |
|---|---|---|---|
| BL-113 | Planned | Create process threat model | Observation, identity, control, disclosure, races, and cleanup |
| BL-114 | Planned | Add paginated process list | Bounded stable results |
| BL-115 | Planned | Add process field selection | Minimize sensitive output |
| BL-116 | Planned | Add process details | Explicit fields and access errors |
| BL-117 | Planned | Add process tree | Bounded depth and partial platform data |
| BL-118 | Planned | Enforce `process.observe` | Registration and execution checks |
| BL-119 | Planned | Implement Managed Process Registry | Thread-safe server-started process ownership |
| BL-120 | Planned | Generate opaque process handles and prevent PID reuse errors | Handles are primary identity; PID is diagnostic only |
| BL-121 | Planned | Define managed process status | Starting/running/exited/failed/stopped/timed-out states |
| BL-122 | Planned | Add `start_process` | Managed engine only, with policy checks |
| BL-123 | Planned | Add `wait_process` | Context, timeout, and final result semantics |
| BL-124 | Planned | Add `read_process_output` | Cursor/bounded incremental output |
| BL-125 | Planned | Add separate stdout/stderr ring buffers | Size limits, truncation markers, and cleanup |
| BL-126 | Planned | Add `stop_process` | Managed handles by default |
| BL-127 | Later | Evaluate external PID control | Requires `process.control.external` or equivalent plus high-risk policy; absent from standard profiles |
| BL-128 | Later | Evaluate `write_process_input` | Separate policy and lifecycle gate |
| BL-129 | Planned | Define process cleanup, restart, and orphan behavior | Shutdown, server crash/restart, TTL, and diagnostics |
| BL-130 | Planned | Limit managed process count and concurrency | Global and profile budgets |
| BL-131 | Planned | Enforce process runtime limits | Defaults, maxima, cancellation, and termination |
| BL-132 | Planned | Define CPU and RAM limit strategy | Windows/Linux capabilities and fallbacks |
| BL-133 | Planned | Redact process command lines and environments | Minimize output and audit data |
| BL-134 | Planned | Implement Windows and Linux process adapters | Equivalent policy outcomes with platform-specific internals |
| BL-135 | Planned | Add process race, lifecycle, and restart tests | Registry races, PID reuse assumptions, cleanup, and limits |

### Command execution epic

| ID | Status | Task | Scope and acceptance notes |
|---|---|---|---|
| BL-136 | Planned | Create command-execution threat model | Executable substitution, injection, environment, roots, output, isolation |
| BL-137 | Planned | Add executable allowlist with IDs | Resolve IDs server-side to approved absolute paths |
| BL-138 | Planned | Enforce no-shell default and argument arrays | No free shell string in standard profile |
| BL-139 | Planned | Restrict working directories to allowed roots | Named-root process permission integration |
| BL-140 | Planned | Define timeout defaults and maxima | Server-enforced deadlines |
| BL-141 | Planned | Limit stdout and stderr separately | Bounded buffers/results and truncation markers |
| BL-142 | Planned | Define stable command result schema | Exit, output references, timeout, status, and bounded diagnostics |
| BL-143 | Planned | Implement `run_command` over Managed Process Engine | Synchronous wrapper only |
| BL-144 | Planned | Add environment allowlist and propagation rules | Minimal explicit environment |
| BL-145 | Planned | Add execution environment and output redaction | Results, diagnostics, and audit events |
| BL-146 | Planned | Implement Windows execution isolation | Least privilege and documented platform limits |
| BL-147 | Planned | Implement Linux execution isolation | Least privilege and documented platform limits |
| BL-148 | Planned | Limit parallel command processes | Profile/global resource budgets |
| BL-149 | Planned | Define least-privilege execution identity | Avoid inherited privileges where practical |
| BL-150 | Later | Keep interactive shell disabled | Separate interactive/high-risk policy decision and threat model |
| BL-151 | Planned | Prevent a second execution engine | Registry and policy invariants/tests |
| BL-152 | Planned | Add Windows/Linux execution security tests | Allowlist, args, roots, env, output, timeout, and isolation |

### System information epic

| ID | Status | Task | Scope and acceptance notes |
|---|---|---|---|
| BL-153 | Planned | Add `system_info` | Controlled OS, architecture, version, and allowed host facts |
| BL-154 | Planned | Expose scoped disk usage | Reuse root-scoped operation `BL-062` |
| BL-155 | Planned | Add filtered environment information | Allowlist and secret exclusion |
| BL-156 | Planned | Add field selection and redaction | Minimize results and host identifiers |
| BL-157 | Planned | Enforce `system.read` capability | Registration and server-side execution checks |
| BL-158 | Later | Evaluate restricted network information | Separate privacy-sensitive decision gate |

### Security epic

| ID | Status | Task | Scope and acceptance notes |
|---|---|---|---|
| BL-159 | Planned | Enforce capabilities server-side | Registration is not the authorization boundary |
| BL-160 | Planned | Add negative capability tests | Every high-risk path and direct-call attempt |
| BL-161 | Planned | Add per-root policy security tests | Read/write, types, limits, links/reparse, working directory |
| BL-162 | Planned | Define process policy model | Observe/manage/external-control separation, managed identity, risk classifications, and policy inputs; enforcement remains in concrete implementation tasks |
| BL-163 | Planned | Define and enforce execution policies | Executables, args, roots, environment, limits, isolation |
| BL-164 | Planned | Add Operations/Job limit security tests | Queue, concurrency, time, temporary data, cleanup, handles |
| BL-165 | Planned | Maintain stateful domain threat models | Filesystem, search, processes, execution, FlashGate modules/providers, and MCP protocol extensions |
| BL-166 | Planned | Define structured audit event model | Bounded redacted decisions, operations, outcomes, and identifiers |
| BL-167 | Planned | Extend secret redaction across new domains | Process, execution, system, jobs, audit, and module/provider outputs |
| BL-168 | Planned | Validate least-privilege execution | Server and child process permissions |
| BL-169 | Planned | Enforce FlashGate module/provider security boundaries | No bypass of policies, functional capabilities, roots, limits, audit, or adapters |
| BL-170 | Planned | Document sandbox boundaries and residual risk | OS, process, filesystem, link/reparse, and configuration limits |
| BL-171 | Planned | Verify MCP annotations are never authorization | Documentation and negative tests |
| BL-172 | Later | Review workflow pinning strategy | Supply-chain hardening and SHA-pin trade-offs |
| BL-173 | Planned | Maintain public security policy | Reporting, supported versions, disclosure, and release gate |
| BL-174 | Planned | Fail closed when no root is explicitly configured | Sprint 3.44 prerequisite: missing root must not expose process working directory; current-directory root requires explicit development opt-in; Sprint 3.41 documents but does not implement |

### Open source and FlashGate modules/providers epic

| ID | Status | Task | Scope and acceptance notes |
|---|---|---|---|
| BL-175 | Done | Confirm current open-source license | Repository and README currently declare GNU GPL v3.0; no license change in Sprint 3.41 |
| BL-176 | Planned | Review license and distribution compatibility before external module contract | Separate factual compatibility review before first external FlashGate module/provider contract; no legal conclusions in backlog |
| BL-177 | Planned | Define governance model | Decision authority, releases, and stewardship |
| BL-178 | Planned | Define maintainer rules | Roles, review, security, and succession |
| BL-179 | Planned | Expand contribution guidelines | Development, testing, documentation, and DCO/CLA decision |
| BL-180 | Later | Add Code of Conduct before public community release | Community expectations and enforcement |
| BL-181 | Planned | Hold FlashGate module/provider contract decision gate | Required before first external provider |
| BL-182 | Planned | Decide FlashGate module/provider identifier rules | No concrete syntax before the contract decision |
| BL-183 | Planned | Define FlashGate module/provider metadata | Name, version, vendor, tools, config, platforms, dependencies |
| BL-184 | Planned | Define module/provider capability declarations | Required and optional functional capabilities with least privilege |
| BL-185 | Planned | Define module/provider security classification | Risk categories and review requirements |
| BL-186 | Planned | Distinguish public, community, vendor, and internal providers | Distribution and support labels do not change security |
| BL-187 | Planned | Define official versus community provider policy | Trust language, signing/update expectations, and support |
| BL-188 | Later | Decide FlashGate provider runtime model | Choose among statically linked packages registered at build time, registered in-process providers, or isolated out-of-process providers over local IPC; a Go module is source/versioning only, not a runtime model |

### Efficiency and MCP contract foundation

| ID | Status | Task | Scope and acceptance notes |
|---|---|---|---|
| BL-189 | Planned | Add startup benchmark | Repeatable cold/warm measurement |
| BL-190 | Planned | Measure idle RSS | Windows/Linux baseline |
| BL-191 | Planned | Measure peak memory, CPU time, and allocations | Representative operations and concurrency |
| BL-192 | Planned | Measure p50 and p95 latency | Reference calls and workflows |
| BL-193 | Planned | Record scanned, read, and written bytes | Server-side primary counters |
| BL-194 | Planned | Measure serialized result sizes | Bytes, characters, entries, and results |
| BL-195 | Planned | Measure `tools/list` size | Bytes and schema count by profile |
| BL-196 | Planned | Measure calls per reference workflow | Local-operation versus model-round-trip comparisons |
| BL-197 | Planned | Add optional schema/response token approximation | Supplemental, not primary, metric |
| BL-198 | Planned | Establish benchmark baselines | Versioned environment and reproducibility notes |
| BL-199 | Planned | Define CI regression budgets | Noise-aware thresholds and review path |
| BL-200 | Planned | Add MCP `outputSchema` | Contract foundation after cleanup |
| BL-201 | Planned | Add MCP `structuredContent` | Align runtime results and declared schema |
| BL-202 | Planned | Review MCP tool annotations | Accurate metadata, never authorization |
| BL-203 | Planned | Define normalized machine-readable errors | Stable categories without raw OS leakage |
| BL-204 | Planned | Evaluate official MCP conformance testing and add schema snapshots | Evaluate official conformance tooling; snapshot names, required fields, descriptions, input/output schemas |
| BL-205 | Planned | Add response-size regression tests | Representative success/error results |
| BL-206 | Planned | Document local deterministic work principle | Prefer local copy/edit/hash/search over model retransmission |
| BL-207 | Planned | Define supported MCP protocol-version strategy | Keep core version-independent; negotiation and compatibility tests for each advertised revision; current implementation remains `2025-11-25` |
| BL-208 | Planned | Define MCP extension-negotiation strategy | Official vendor-prefix/slash identifiers, capability negotiation, downgrade/mismatch tests, no authorization implication |
| BL-209 | Planned | Decide MCP Tasks compatibility | Evaluate `io.modelcontextprotocol/tasks` and supported client behavior before asynchronous MCP exposure |
| BL-210 | Planned | Map internal operation lifecycle to MCP Tasks | Define tested state, result, error, cancellation, TTL, and redaction mapping; internal states may be more detailed |
| BL-211 | Planned | Decide fallback when MCP Tasks is unavailable | Bounded synchronous result or explicit capability error; no ad hoc custom job-tool contract |
| BL-212 | Planned | Validate all input/output schemas as JSON Schema 2020-12 | Cover current and future schemas, dialect declarations, snapshots, and protocol-version compatibility |

### CI, release, and quality epic

| ID | Status | Task | Scope and acceptance notes |
|---|---|---|---|
| BL-213 | Planned | Add release notes or tag-based release workflow | Release maturity after technical rename |
| BL-214 | Planned | Add artifact verification | Run version/help and validate name/platform/metadata |
| BL-215 | Planned | Run benchmark suite in CI | Stable selection and artifacted results |
| BL-216 | Planned | Compare benchmark baselines in CI | Budgets from `BL-199` |
| BL-217 | Planned | Validate PowerShell and Bash scripts | Syntax/lint and smoke portability |
| BL-218 | Planned | Run race detector for stateful components | Jobs, process registry, output buffers, shutdown |
| BL-219 | Planned | Add Windows/Linux process test jobs | Observation and managed lifecycle |
| BL-220 | Planned | Add Operations/Job test jobs | Cancellation, timeout, cleanup, and race coverage |
| BL-221 | Planned | Verify FlashGate release artifact names | After Sprint 3.42, including archives and summaries |
| BL-222 | Planned | Enforce `tools/list` budget | Profile-specific size regression |
| BL-223 | Planned | Run schema snapshot checks in CI | Contract changes require explicit review |
| BL-224 | Planned | Run response-size tests in CI | Prevent unbounded contract regressions |
| BL-225 | Planned | Search repository for legacy names after Sprint 3.42 | Allow only migration/history exceptions |
| BL-226 | Planned | Keep standard test/vet/lint/build gates | Preserve existing validation on Windows and Linux |

### Sprint 3.42 technical rename

| ID | Status | Task | Scope and acceptance notes |
|---|---|---|---|
| BL-227 | Planned | Rename local folder to `flashgate-mcp` | User-coordinated move after clean review |
| BL-228 | Planned | Rename GitHub repository to `flashgate-mcp` | Manual owner action documented |
| BL-229 | Planned | Update Git remote URL | Verify fetch/push and redirects; no auth change |
| BL-230 | Planned | Update Go module and imports | `github.com/blacksheepkhan/flashgate-mcp` |
| BL-231 | Planned | Rename binary to `flashgate-mcp` | Windows/Linux build and usage |
| BL-232 | Planned | Change MCP server implementation name (`serverInfo.name`) to `flashgate` | Initialize response and tests |
| BL-233 | Planned | Review package and command paths | Avoid unrelated refactoring |
| BL-234 | Planned | Update README, changelog, and documentation names | Preserve migration/history context |
| BL-235 | Planned | Update PowerShell and Bash scripts | Paths, errors, examples, and smoke expectations |
| BL-236 | Planned | Update CI and release artifact names | Workflows only in rename sprint |
| BL-237 | Planned | Update installation and configuration examples | New folder and binary names |
| BL-238 | Planned | Update smoke tests | Implementation name (`serverInfo.name`), binary, paths, and current tool names |
| BL-239 | Planned | Search all files for legacy names | Classify every remaining occurrence |
| BL-240 | Planned | Write technical rename migration note | Old/new repository, module, binary, implementation name (`serverInfo.name`), local folder |
| BL-241 | Planned | Verify GitHub redirect behavior | Clone links, issues, releases, actions, and remote guidance |
| BL-242 | Planned | Document manual repository rename action | Exact prerequisites, action, verification, and rollback guidance |
| BL-243 | Planned | Keep rename sprint functionally neutral | No feature or tool-contract changes mixed in |

### Sprint 3.43 pre-1.0 tool contract cleanup

| ID | Status | Task | Scope and acceptance notes |
|---|---|---|---|
| BL-244 | Planned | Rename `list_files` to `list_directory` | Code, schema, tests, docs, examples, smoke |
| BL-245 | Planned | Rename `stat_path` to `get_path_info` | Clear general terminology |
| BL-246 | Planned | Rename `mkdir` to `create_directory` | Avoid Unix jargon |
| BL-247 | Planned | Remove `exists_path` | No compatibility alias for unpublished contract |
| BL-248 | Planned | Remove `rename_path` | `move_path` covers rename and movement |
| BL-249 | Planned | Define `move_path` rename/move semantics | Same/cross-volume and overwrite behavior |
| BL-250 | Planned | Define missing-path `get_path_info` result | `{ "exists": false, "path": ... }` |
| BL-251 | Planned | Define normalized filesystem error codes | Distinguish missing, access, root, invalid, type, I/O |
| BL-252 | Planned | Review input/output schemas and required fields | Remove ambiguity and accidental optionality |
| BL-253 | Planned | Optimize tool descriptions | Accurate, compact, model-useful language |
| BL-254 | Planned | Update JSON schema snapshots and unit tests | Explicit breaking-contract approval |
| BL-255 | Planned | Update smoke tests and MCP `tools/call` tests | New baseline and removed tools |
| BL-256 | Planned | Update tool docs, client examples, and catalog | One coordinated contract view |
| BL-257 | Planned | Document breaking changes in changelog | No artificial deprecation compatibility |

### Sprint 3.44 Codex read-only activation preparation

| ID | Status | Task | Scope and acceptance notes |
|---|---|---|---|
| BL-258 | Planned | Add Codex read-only configuration example | Use FlashGate binary, implementation name (`serverInfo.name`), and cleaned tool names |
| BL-259 | Planned | Add Claude Desktop configuration example | Windows/Linux with renamed artifact |
| BL-260 | Planned | Add general MCP client examples | Minimal local STDIO configuration |
| BL-261 | Planned | Add read-only troubleshooting guide | Root, permissions, JSON-RPC, binary, and profile/tool-list issues |
| BL-262 | Planned | Create activation checklist | Build, root, read-only profile, tools/list, smoke, rollback |
| BL-263 | Planned | Update read-only smoke test | FlashGate paths and cleaned read tools |
| BL-264 | Planned | Validate new documentation paths and links | No legacy installation path except migration |
| BL-265 | Planned | Document non-developer read-only validation | Script-based checks without Go |
| BL-266 | Planned | Keep activation external to preparation sprint | No live Codex/config/auth change without separate confirmation |

### Documentation and client compatibility epic

| ID | Status | Task | Scope and acceptance notes |
|---|---|---|---|
| BL-267 | Done | Link planning and history from README | Backlog, roadmap, changelog, architecture, security, ADRs |
| BL-268 | Planned | Keep CHANGELOG updated each sprint | Added/Changed/Security and breaking changes |
| BL-269 | Planned | Keep BACKLOG updated each sprint | Canonical IDs, sprint refs, migration rule, completed status |
| BL-270 | Planned | Maintain FlashGate project identity reference | Name, tagline, scope, transition, planned identifiers |
| BL-271 | Planned | Maintain architecture and ADRs | Current/target/planned/deferred separation |
| BL-272 | Planned | Document benchmark method and baselines | Environments, noise, budgets, results |
| BL-273 | Planned | Document capabilities, profiles, and named roots | Configuration and security model |
| BL-274 | Planned | Document Operations/Job Manager | Handles, states, limits, lifecycle, cleanup |
| BL-275 | Planned | Document process and execution security | Handles, PIDs, allowlists, isolation, redaction |
| BL-276 | Planned | Document open-source governance, modules/providers, and MCP extensions | Separate local provider and negotiated protocol-extension decision gates |
| BL-277 | Planned | Maintain non-developer smoke-test documentation | PowerShell/Bash and expected results |
| BL-278 | Planned | Review README, CHANGELOG, and BACKLOG each sprint | Prevent stale identity, implementation, and planning claims |

## Cross-epic rules

- Security tasks apply to their domain tasks without duplicating canonical definitions.
- Operations/jobs are optional lifecycle infrastructure; short synchronous operations may run directly and domain logic/ownership stays outside the manager.
- Benchmarks and threat models must justify separate binaries, IPC, indexes, or external adapters/providers.
- FlashGate modules/providers and MCP protocol extensions are separate concepts and contracts.
- Deprecated MCP Roots is never the foundation of named-root authorization.
- MCP annotations never replace server-side authorization.
- Planned tool cleanup and technical rename occur only in their dedicated sprints.
- Before 1.0, breaking changes are allowed but require coordinated tests, documentation, examples, smoke tests, and changelog entries.
