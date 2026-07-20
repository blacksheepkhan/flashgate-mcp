# Backlog ID Migration - 2026-07-20

## Status

Completed documentation migration in PR #19. The new IDs become authoritative when the PR is merged.

## Reason

Two operating-system-specific release metadata tasks were inserted at the technically appropriate position before the existing artifact-verification task:

- `BL-246`: native Windows file and product metadata
- `BL-247`: native Linux binary and package metadata

A separate dependency-maintenance task was appended as `BL-324`:

- `BL-324`: Dependabot security and version updates

The project rule requires all active canonical IDs to remain continuous after inserting tasks.

## Canonical range

The canonical catalog now remains continuous through:

```text
BL-001 through BL-324
```

## Renumbering rule

All previously active canonical IDs from `BL-246` through `BL-321` move forward by two positions:

```text
BL-246 -> BL-248
...
BL-321 -> BL-323
```

The new `BL-324` task is appended and therefore does not replace a previous ID.

Historical migration files remain immutable. In particular, `docs/backlog-id-migration-2026-07-17.md` and earlier dated records retain the identifiers valid at their respective dates.

## Mapping

| Previous ID | New ID | Note |
|---|---|---|
| BL-246 | BL-248 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-247 | BL-249 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-248 | BL-250 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-249 | BL-251 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-250 | BL-252 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-251 | BL-253 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-252 | BL-254 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-253 | BL-255 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-254 | BL-256 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-255 | BL-257 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-256 | BL-258 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-257 | BL-259 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-258 | BL-260 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-259 | BL-261 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-260 | BL-262 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-261 | BL-263 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-262 | BL-264 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-263 | BL-265 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-264 | BL-266 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-265 | BL-267 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-266 | BL-268 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-267 | BL-269 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-268 | BL-270 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-269 | BL-271 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-270 | BL-272 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-271 | BL-273 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-272 | BL-274 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-273 | BL-275 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-274 | BL-276 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-275 | BL-277 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-276 | BL-278 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-277 | BL-279 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-278 | BL-280 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-279 | BL-281 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-280 | BL-282 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-281 | BL-283 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-282 | BL-284 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-283 | BL-285 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-284 | BL-286 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-285 | BL-287 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-286 | BL-288 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-287 | BL-289 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-288 | BL-290 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-289 | BL-291 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-290 | BL-292 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-291 | BL-293 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-292 | BL-294 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-293 | BL-295 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-294 | BL-296 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-295 | BL-297 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-296 | BL-298 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-297 | BL-299 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-298 | BL-300 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-299 | BL-301 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-300 | BL-302 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-301 | BL-303 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-302 | BL-304 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-303 | BL-305 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-304 | BL-306 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-305 | BL-307 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-306 | BL-308 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-307 | BL-309 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-308 | BL-310 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-309 | BL-311 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-310 | BL-312 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-311 | BL-313 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-312 | BL-314 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-313 | BL-315 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-314 | BL-316 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-315 | BL-317 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-316 | BL-318 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-317 | BL-319 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-318 | BL-320 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-319 | BL-321 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-320 | BL-322 | Shifted by the insertion of the Windows and Linux metadata tasks |
| BL-321 | BL-323 | Shifted by the insertion of the Windows and Linux metadata tasks |

## Sprint and steering reference changes

| Reference | Previous IDs | New IDs |
|---|---|---|
| Sprint 3.42 | BL-262–BL-278 | BL-264–BL-280 |
| Sprint 3.43 | BL-279–BL-292 | BL-281–BL-294 |
| Sprint 3.44 | BL-174, BL-293–BL-301 | BL-174, BL-295–BL-303 |
| Sprint 3.53 CI/race subset | BL-250–BL-252 | BL-252–BL-254 |
| Sprint 3.58 | BL-241–BL-261, BL-303–BL-310, BL-312–BL-321 | BL-241–BL-263, BL-305–BL-312, BL-314–BL-324 |
| Version 1.0 release gate | BL-261 | BL-263 |
| Provider documentation | BL-311 | BL-313 |
| Documentation epic | BL-302–BL-313 | BL-304–BL-315 |
| PR #15 review follow-up | BL-314–BL-321 | BL-316–BL-323 |

## New tasks

| ID | Task | Position |
|---|---|---|
| BL-246 | Embed native Windows file and product metadata | Inserted before artifact verification |
| BL-247 | Define and publish native Linux binary and package metadata | Inserted before artifact verification |
| BL-324 | Configure Dependabot security and version updates | Appended after the PR #15 follow-up tasks |

## Active reference treatment

Active normative documents were updated directly to the new canonical IDs:

- `benchmarks/README.md`
- `docs/adr/0012-resource-token-efficiency-and-pre-1-0-contracts.md`
- `docs/architecture.md`
- `docs/native-multi-mode-runtime-and-service-plan.md`
- `docs/specification.md`
- `docs/version-1-scope-and-release-boundary.md`

Dated evidence documents retain their original text and carry an appended 2026-07-20 correction section:

- `docs/benchmarks/pr-15-independent-review-follow-up-2026-07-18.md`
- `docs/benchmarks/sprint-045d-resource-latency-baseline.md`
- `docs/comparative-mcp-review-2026-07-17.md`

Earlier backlog migration records, including `docs/backlog-id-migration-2026-07-17-version1-efficiency-hybrid.md`, remain byte-for-byte unchanged.

## Validation result

- Canonical IDs are continuous from `BL-001` through `BL-324`.
- Every canonical ID occurs exactly once in `BACKLOG.md`.
- The three new tasks occur exactly once.
- All previously active IDs from `BL-246` through `BL-321` are mapped forward by two positions.
- Active sprint, release-gate, documentation, provider, benchmark-CI, and PR #15 references use the new IDs.
- Earlier dated migration files remain unchanged as historical evidence.
