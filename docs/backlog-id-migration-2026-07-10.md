# Backlog ID Migration - 2026-07-10

## Reason

Sprint 3.41 consolidated planning into one authoritative task catalog. The previous backlog contained gaps, reused summary references, duplicated completed summaries, and a separate `BL-D` history series. Continuous IDs make sprint references machine-checkable and prevent tasks from being fully defined in multiple places.

## New rule

Canonical tasks start at `BL-001`, use three digits, have no gaps, and are defined exactly once in `BACKLOG.md`. Sprint tables reference IDs only. Historical summaries do not receive separate task IDs.

> Neue Aufgaben werden an der fachlich richtigen Position eingefügt. Danach werden alle nachfolgenden BL-IDs repositoryweit fortlaufend umnummeriert.

This file records the still-uncommitted Sprint 3.41 migration and may be corrected until Sprint 3.41 is merged. After merge, every dated migration file is immutable history. Any later full renumbering creates a new dated migration file; earlier tables are not overwritten. A small migration index may be added if multiple migrations exist.

## Complete old-to-new mapping

An em dash means that the old value was a non-canonical undocumented summary reference rather than a task definition. Ranges mean the old broad task was split into explicit canonical work. “Merged” identifies duplicate or overlapping old entries consolidated into one current task.

| Old ID | New ID | Disposition |
|---|---|---|
| BL-001 | — | Removed undocumented Sprint 3.40 summary reference; not a canonical task |
| BL-002 | — | Removed undocumented Sprint 3.40 summary reference; not a canonical task |
| BL-004 | BL-036 | Renumbered |
| BL-005 | BL-189–BL-198 | Expanded into measurable benchmark tasks |
| BL-007 | BL-055 | Renumbered |
| BL-008 | BL-044 | Renumbered |
| BL-009 | BL-043 | Renumbered |
| BL-010 | BL-048 | Expanded to batch hashing; single-item behavior remains covered |
| BL-011 | BL-068 | Expanded with search threat model |
| BL-012 | BL-070 | Renumbered |
| BL-013 | BL-072 | Renumbered |
| BL-014 | BL-080 | Renumbered |
| BL-015 | BL-113 | Renumbered and expanded |
| BL-016 | BL-114 | Renumbered and paginated |
| BL-017 | BL-116 | Renumbered |
| BL-018 | BL-126 | Renumbered; managed handles are the default |
| BL-019 | BL-122 | Renumbered; uses Managed Process Engine |
| BL-020 | BL-117 | Renumbered |
| BL-021 | BL-136 | Renumbered and expanded threat model |
| BL-022 | BL-137 | Renumbered; executable IDs and absolute resolution added |
| BL-023 | BL-140 | Renumbered |
| BL-024 | BL-141 | Renumbered |
| BL-025 | BL-142 | Renumbered |
| BL-026 | BL-139 | Renumbered and tied to named roots |
| BL-027 | BL-150 | Renumbered; explicit high-risk decision gate |
| BL-028 | BL-153 | Renumbered |
| BL-029 | BL-062, BL-154 | Split into scoped operation and system exposure |
| BL-030 | BL-155 | Renumbered and allowlisted |
| BL-031 | BL-158 | Renumbered as privacy-sensitive gate |
| BL-032 | BL-100 | Renumbered and expanded capability model |
| BL-033 | BL-103, BL-110 | Split into profile configuration and dynamic registration |
| BL-034 | BL-021 | Merged with old BL-058, BL-059, and BL-D019 |
| BL-035 | BL-022 | Renumbered |
| BL-036 | BL-172 | Renumbered |
| BL-037 | BL-259 | Renumbered and moved to Sprint 3.44 |
| BL-038 | BL-260 | Renumbered and moved to Sprint 3.44 |
| BL-039 | BL-261 | Renumbered and moved to Sprint 3.44 |
| BL-040 | BL-023 | Renumbered, status remains Done |
| BL-041 | BL-213 | Renumbered |
| BL-042 | BL-214 | Renumbered |
| BL-043 | BL-217 | Renumbered |
| BL-044 | BL-267 | Renumbered, status now Done in Sprint 3.41 |
| BL-045 | BL-268 | Renumbered |
| BL-046 | BL-269 | Renumbered |
| BL-047 | BL-265 | Renumbered and moved to Sprint 3.44 |
| BL-048 | BL-275 | Renumbered and expanded |
| BL-049 | BL-015 | Merged with BL-D015; status remains Done |
| BL-050 | BL-017 | Merged with old BL-051, BL-052, and BL-D017 |
| BL-051 | BL-017 | Merged with old BL-050, BL-052, and BL-D017 |
| BL-052 | BL-017 | Merged with old BL-050, BL-051, and BL-D017 |
| BL-053 | BL-016 | Merged with old BL-054 and BL-D016 |
| BL-054 | BL-016 | Merged with old BL-053 and BL-D016 |
| BL-055 | BL-018 | Merged with BL-D018 |
| BL-056 | BL-019 | Merged into the completed limit baseline |
| BL-057 | BL-020 | Merged into the completed filesystem/response limit baseline |
| BL-058 | BL-021 | Merged with old BL-034, BL-059, and BL-D019 |
| BL-059 | BL-021 | Merged with old BL-034, BL-058, and BL-D019 |
| BL-060 | BL-024 | Merged with BL-D020; status remains Done |
| BL-061 | BL-025 | Split from BL-D020 and retained as its own completed task |
| BL-D001 | BL-001 | `BL-D` history number retired; canonical Done task retained |
| BL-D002 | BL-002 | `BL-D` history number retired; canonical Done task retained |
| BL-D003 | BL-003 | `BL-D` history number retired; canonical Done task retained |
| BL-D004 | BL-004 | `BL-D` history number retired; canonical Done task retained |
| BL-D005 | BL-005 | `BL-D` history number retired; canonical Done task retained |
| BL-D006 | BL-006 | `BL-D` history number retired; canonical Done task retained |
| BL-D007 | BL-007 | `BL-D` history number retired; canonical Done task retained |
| BL-D008 | BL-008 | `BL-D` history number retired; canonical Done task retained |
| BL-D009 | BL-009 | `BL-D` history number retired; canonical Done task retained |
| BL-D010 | BL-010 | `BL-D` history number retired; canonical Done task retained |
| BL-D011 | BL-011 | `BL-D` history number retired; canonical Done task retained |
| BL-D012 | BL-012 | `BL-D` history number retired; canonical Done task retained |
| BL-D013 | BL-013 | `BL-D` history number retired; canonical Done task retained |
| BL-D014 | BL-014 | `BL-D` history number retired; canonical Done task retained |
| BL-D015 | BL-015 | `BL-D` history number retired; merged with old BL-049 |
| BL-D016 | BL-016 | `BL-D` history number retired; merged with old BL-053 and BL-054 |
| BL-D017 | BL-017 | `BL-D` history number retired; merged with old BL-050–BL-052 |
| BL-D018 | BL-018 | `BL-D` history number retired; merged with old BL-055 |
| BL-D019 | BL-019–BL-022 | `BL-D` history number retired; broad summary split/merged into completed limits, diagnostics, and defaults |
| BL-D020 | BL-023–BL-025 | `BL-D` history number retired; broad smoke summary split/merged into three completed tasks |

Old `BL-003` and `BL-006` did not exist in the pre-migration backlog and therefore have no migration row.

## New tasks introduced during Sprint 3.41 review

These tasks have no pre-migration ID:

| New ID | Task |
|---|---|
| BL-174 | Fail closed when no root is explicitly configured |
| BL-176 | Review license and distribution compatibility before external module contract |
| BL-207 | Define supported MCP protocol-version strategy |
| BL-208 | Define MCP extension-negotiation strategy |
| BL-209 | Decide MCP Tasks compatibility |
| BL-210 | Map internal operation lifecycle to MCP Tasks |
| BL-211 | Decide fallback when MCP Tasks is unavailable |
| BL-212 | Validate all input/output schemas as JSON Schema 2020-12 |

## Historical references

Old IDs may still appear in commits, pull requests, exported reports, or prior discussions. They identify the pre-2026-07-10 backlog and must be interpreted through this table. Outside this migration document, old IDs should remain only when quoting necessary historical records.

## Validation summary

- Canonical task count: **278**
- Lowest canonical ID: **BL-001**
- Highest canonical ID: **BL-278**
- Missing canonical numbers: **none**
- Duplicate canonical definitions: **none**
- Canonical `BL-D` definitions: **none**

The canonical sequence is confirmed continuous with no gaps.
