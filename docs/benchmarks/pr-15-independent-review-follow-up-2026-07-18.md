# PR #15 independent-review follow-up (2026-07-18)

The independent review of PR #15 identified one baseline-output protection
blocker and eight additional findings. The blocker is addressed in the PR by the
handle-bound non-authoritative output writer and its Windows/Linux policy tests.
This note records the remaining findings without claiming that they are fixed.

`BACKLOG.md` is the canonical steering source. The findings are scheduled as
separate work after PR #15 is merged:

| Severity | Canonical task | Finding | Primary components |
|---|---|---|---|
| Major | BL-314 | Validator recomputes budgets and applies strict schema/unknown-field checks | artifact validator, schema, budgets |
| Major | BL-315 | Semantic workflow gates enforce expected bytes and entries | workflows, runner, budgets |
| Major | BL-316 | Authoritative runs require an isolated explicit corpus parent | controllers, workspace validation |
| Major | BL-317 | Baselines carry verifiable binary/build/controller/workspace/host provenance | schema, controllers, host gate |
| Major | BL-318 | Native Linux race plus Windows/Linux policy/window coverage runs in CI | workflows and platform scripts |
| Minor | BL-319 | Linux clock ticks are obtained rather than assumed to be 100 Hz | Linux resource collector |
| Minor | BL-320 | BL-190 and BL-198 status is reconciled without rewriting history | BACKLOG and Sprint 3.45d reporting |
| Minor | BL-321 | Testing documentation matches actual copy/search benchmark coverage | `docs/testing.md`, benchmark inventory |

The full risk and acceptance criteria for every item are defined once in the
canonical task rows. Existing broader tasks such as BL-247, BL-248, BL-250,
BL-258, BL-260, and BL-307 do not independently cover the complete review
finding, so the new tasks link the precise acceptance boundaries without
changing those older tasks. None of BL-314 through BL-321 is implemented by the
baseline-output blocker fix.

## Backlog ID correction - 2026-07-20

This dated review record preserves the identifiers valid when it was written.
The active canonical IDs changed through `docs/backlog-id-migration-2026-07-20.md`.

Current mappings for this document:

| Historical ID | Current ID |
|---|---|
| BL-247 | BL-249 |
| BL-248 | BL-250 |
| BL-250 | BL-252 |
| BL-258 | BL-260 |
| BL-260 | BL-262 |
| BL-307 | BL-309 |
| BL-314 | BL-316 |
| BL-315 | BL-317 |
| BL-316 | BL-318 |
| BL-317 | BL-319 |
| BL-318 | BL-320 |
| BL-319 | BL-321 |
| BL-320 | BL-322 |
| BL-321 | BL-323 |

The active PR #15 follow-up range is therefore `BL-316` through `BL-323`.
The original review text above remains unchanged as historical evidence.
