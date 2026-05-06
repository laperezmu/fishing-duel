# Pesca Constitution

## Core Principles

### I. Spec First, Then Code
Every non-trivial change starts in spec-kit and not in ad-hoc docs or direct implementation. A feature is ready for coding only after the spec captures the problem, scope, acceptance criteria, assumptions, and open questions with enough precision to avoid avoidable rework.

### II. Single Source of Truth
Active product and engineering workflow lives in GitHub Issues plus spec-kit artifacts under `.specify/` and `specs/`. The issue tracker owns backlog state, prioritization, milestones, and labels; `specs/` and `.specify/` own feature definition, planning, and task breakdown. `PROJECT_CONTEXT.md` and supporting docs should only point to this workflow, not recreate parallel backlog systems.

### III. Modular Go Boundaries
Code stays aligned with Effective Go and pragmatic SOLID in Go. Runtime, content/configuration, presentation, and bootstrap concerns stay clearly separated. Packages should keep focused ownership, explicit APIs, stable vocabulary, and avoid speculative abstractions.

### IV. Testable Delivery
Acceptance criteria must be verifiable. Production code changes should add or update automated tests when behavior changes materially, and every implementation plan must state the intended validation path, including `go test ./...` and `golangci-lint run` when applicable.

### V. Small, Intentional Evolution
Prefer thin vertical slices over broad rewrites. New work should minimize coupling, preserve existing gameplay intent, and record meaningful tradeoffs in the spec or plan instead of hiding them inside code.

## Technical Guardrails

- Primary language is Go.
- Main delivery surfaces today are CLI executables under `cmd/` and supporting packages under `internal/`.
- Domain packages must remain UI-agnostic.
- `cmd/` composes dependencies, `internal/app/` coordinates flows, and presentation adapters translate technical state into user-facing output.
- Avoid magic values when they represent domain policy, UX contract, or shared defaults; prefer named constants or typed catalogs.
- Refactor opportunistically when new work reveals unstable boundaries or repeated logic, but keep refactors scoped to the active spec.

## Workflow

1. Capture non-trivial work as a GitHub Issue using the repository issue template and assign the correct milestone, labels, and ordered ID.
2. Run `/speckit.constitution` only when principles need amendment.
3. Create or refine the feature with `/speckit.specify` once the issue is prioritized.
4. Clarify ambiguity with `/speckit.clarify` before planning when needed.
5. Produce the implementation approach with `/speckit.plan`.
6. Break work down with `/speckit.tasks`.
7. Implement from the generated artifacts with `/speckit.implement` or manual execution aligned to the tasks.
8. Validate with the tests and checks declared in the plan before opening or updating a PR.
9. Close the issue only after the scoped work is merged or explicitly cancelled.

## Backlog and Tracking Rules

- Every non-trivial initiative must have a GitHub Issue before implementation starts.
- The issue title should use the ordered ID for its milestone, for example `[RMF-06] ...`.
- Issues should carry at least one `stream:*`, `area:*`, `type:*`, `priority:*`, and `status:*` label.
- Suggested issue lifecycle is `status:ready-for-spec` -> `status:in-spec` -> `status:planned` -> `status:ready-for-delivery`.
- Do not recreate active backlog lists in `docs/`, ad-hoc markdown files, or `PROJECT_CONTEXT.md`.

## Governance

- This constitution supersedes the deprecated planning workflow documented in legacy docs.
- New active features must not be created in `docs/features/`, `docs/backlog/`, or any other repo-local backlog document.
- Any amendment must update this file and, if needed, the spec-kit templates in `.specify/templates/overrides/`.
- Reviews must check both implementation correctness and compliance with this constitution.

**Version**: 1.1.0 | **Ratified**: 2026-05-04 | **Last Amended**: 2026-05-06
