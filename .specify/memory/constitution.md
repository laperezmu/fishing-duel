# Pesca Constitution

## Core Principles

### I. Spec First, Then Code
Every non-trivial change starts in spec-kit and not in ad-hoc docs or direct implementation. A feature is ready for coding only after the spec captures the problem, scope, acceptance criteria, assumptions, and open questions with enough precision to avoid avoidable rework.

### II. Single Source of Truth
Active product and engineering workflow lives in spec-kit artifacts under `.specify/` and `specs/`. Legacy material under `docs/features/`, `docs/backlog/`, `docs/discoveries/`, and `PROJECT_CONTEXT.md` is historical context only unless a current spec explicitly references it.

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

1. Run `/speckit.constitution` only when principles need amendment.
2. Create a new feature with `/speckit.specify`.
3. Clarify ambiguity with `/speckit.clarify` before planning when needed.
4. Produce the implementation approach with `/speckit.plan`.
5. Break work down with `/speckit.tasks`.
6. Implement from the generated artifacts with `/speckit.implement` or manual execution aligned to the tasks.
7. Validate with the tests and checks declared in the plan before opening or updating a PR.

## Governance

- This constitution supersedes the deprecated planning workflow documented in legacy docs.
- New active features must not be created in `docs/features/` or `docs/backlog/`.
- Any amendment must update this file and, if needed, the spec-kit templates in `.specify/templates/overrides/`.
- Reviews must check both implementation correctness and compliance with this constitution.

**Version**: 1.0.0 | **Ratified**: 2026-05-04 | **Last Amended**: 2026-05-04
