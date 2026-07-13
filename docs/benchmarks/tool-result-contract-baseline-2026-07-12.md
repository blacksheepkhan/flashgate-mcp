# Tool-result contract benchmark baseline – 2026-07-12

## Scope

Sprint 3.45a compares serialization of the historical direct domain object, MCP TextContent-only (variant A), and the selected TextContent plus `structuredContent` contract (variant B). The benchmark measures only result construction and JSON serialization; it does not include filesystem I/O, STDIO process startup, client decoding, or model execution.

## Environment

| Field | Value |
|---|---|
| Date | 2026-07-12 |
| Go | `go1.26.4` |
| OS/architecture | `windows/amd64` |
| CPU | AMD Ryzen 7 3700U with Radeon Vega Mobile Gfx |
| Logical benchmark parallelism | 8 |
| Command | `go test -run '^$' -bench '^BenchmarkCallToolResultSerialization$' -benchmem -benchtime=500ms -count=1 ./internal/mcp/tools/...` |
| Payload command | same benchmark with `-benchtime=200ms` after moving the non-timed payload metric behind `ResetTimer` |

Each benchmark case runs for at least the requested benchtime; the `N` column is the actual number of serialization repetitions selected by Go. `B/op` and `allocs/op` are reported by `b.ReportAllocs`. `payload bytes` is the compact UTF-8 JSON length of the JSON-RPC `result` value only, not the complete JSON-RPC response envelope.

## Fixtures

| Fixture | Shape |
|---|---|
| `path_info_missing` | `{path, exists:false}` |
| `path_info_existing` | existing file metadata |
| `directory_small` | two entries |
| `directory_500_entries` | 500 deterministic entries |
| `text_file_small` | 26-byte text domain payload |
| `text_file_64kib` | 65,536-byte text domain payload |

## Results

| Fixture | Variant | N | ns/op | B/op | allocs/op | payload bytes |
|---|---|---:|---:|---:|---:|---:|
| Missing info | historical | 1,312,872 | 423.9 | 48 | 1 | 44 |
| Missing info | text-only | 368,496 | 1,628 | 248 | 5 | 89 |
| Missing info | text + structured | 164,035 | 3,963 | 416 | 6 | 154 |
| Existing info | historical | 744,351 | 804.4 | 96 | 1 | 85 |
| Existing info | text-only | 271,147 | 2,421 | 392 | 5 | 138 |
| Existing info | text + structured | 100,330 | 7,003 | 640 | 6 | 244 |
| Small directory | historical | 464,942 | 1,081 | 96 | 1 | 91 |
| Small directory | text-only | 187,477 | 2,925 | 408 | 5 | 148 |
| Small directory | text + structured | 78,147 | 7,821 | 672 | 6 | 260 |
| 500-entry directory | historical | 2,734 | 214,770 | 27,272 | 1 | 25,897 |
| 500-entry directory | text-only | 1,144 | 439,024 | 87,546 | 5 | 29,938 |
| 500-entry directory | text + structured | 511 | 1,311,745 | 143,297 | 6 | 55,856 |
| Small text | historical | 1,000,000 | 577.9 | 64 | 1 | 51 |
| Small text | text-only | 326,277 | 1,864 | 296 | 5 | 97 |
| Small text | text + structured | 126,340 | 4,241 | 464 | 6 | 169 |
| 64-KiB text | historical | 2,870 | 193,870 | 73,754 | 1 | 65,563 |
| 64-KiB text | text-only | 1,314 | 450,054 | 222,785 | 5 | 65,608 |
| 64-KiB text | text + structured | 271 | 2,150,997 | 370,628 | 6 | 131,192 |

Timing/allocation values are from the final 500-ms review run after constructor object-validation and ownership hardening. Payload bytes remained identical to the initial 200-ms measurement.

## Interpretation

Variant B costs one additional wrapper representation over variant A and duplicates the domain JSON on the wire. Small results therefore have high relative envelope overhead; the 64-KiB text response is almost exactly doubled. The selected implementation also uses six allocations per call in these fixtures, compared with one historically and five for text-only.

Variant B remains selected because MCP requires `content[]`, Codex consumes that content, and `structuredContent` preserves a directly machine-readable typed object for current and future clients. Protocol conformance and client compatibility are hard requirements; payload reduction cannot retain the invalid historical form.

No CI threshold is set from this single Windows machine and single recorded run. Later BL-194/198/199/205 work must collect repeated Windows/Ubuntu baselines, quantify noise, and decide reviewable budgets. A rough bytes/4 token heuristic may be used for orientation only; it is not an exact tokenizer measurement and is intentionally not presented as a benchmark result here.
