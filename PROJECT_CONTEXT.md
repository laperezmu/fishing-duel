# Project Context

This file is a lightweight pointer to the active project workflow.

## Active Workflow

- GitHub Issues is the operational backlog.
- `specs/` and `.specify/` are the active source of truth for feature definition and planning.
- The official working agreement lives in `.specify/memory/constitution.md`.
- The day-to-day operating guide lives in `docs/workflow-github-issues-y-specify.md`.

## Delivery Path

For non-trivial work, use this sequence:

1. Create or refine a GitHub Issue.
2. Run `/speckit.specify`.
3. Run `/speckit.clarify` when ambiguity remains.
4. Run `/speckit.plan`.
5. Run `/speckit.tasks`.
6. Implement from the generated artifacts.
7. Validate with the checks declared in the plan.
8. Open or update the PR.

## Engineering Baseline

- Primary language is Go.
- Domain packages remain UI-agnostic.
- `cmd/` composes dependencies, `internal/app/` coordinates flows, and presentation adapters translate technical state for each UI.
- Favor small vertical slices, explicit boundaries, and verifiable acceptance criteria.
