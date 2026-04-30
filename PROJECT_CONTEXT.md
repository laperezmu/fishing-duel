# Project Context

## Feature Integration Workflow

All new features start from a plan before implementation begins.

Each plan must include:
- Objective
- Acceptance criteria
- Scope of the change
- A concrete implementation proposal with the intended organization or direction already stated

Implementation should not start without at least this minimum planning step.

Before implementation starts:
- Create a new git branch for the feature
- Use the implementation plan document name as the branch name
- Continue all feature work on that branch until the PR is merged

## Delivery Process

After implementation is complete:
- Run `go test ./...`
- Run `golangci-lint run`
- Create a pull request in the repository
- Include a detailed description of the integrated changes in the PR body

The pull request is the handoff point for review.

## Review Process

- The user reviews the PR manually
- Review feedback is provided through the console, not through PR comments
- Any requested changes are handled in a follow-up iteration from the console

## Completion Criteria

Integration is considered complete only when the user manually closes and merges the pull request.

That merge point is the end of the feature integration cycle.

After merge:
- Return to `main`
- Continue the next feature from `main`

## Feature Plan Template

Use this template as the default starting point for any new feature:

```md
# Plan de feature: <nombre corto>

## Objetivo
<que problema resuelve esta feature y cual es el resultado esperado>

## Criterios de aceptacion
- <criterio verificable 1>
- <criterio verificable 2>
- <criterio verificable 3>

## Scope
### Incluye
- <cambio incluido 1>
- <cambio incluido 2>

### No incluye
- <limite explicito 1>
- <limite explicito 2>

## Propuesta de implementacion
- <cambio tecnico principal>
- <paquetes/archivos que probablemente se tocaran>
- <impacto en arquitectura, flujo o UX>

## Validacion
- <test manual o automatizado 1>
- <test manual o automatizado 2>

## Riesgos o decisiones abiertas
- <riesgo, tradeoff o duda pendiente>
```

Minimum required sections before implementation starts:
- Objetivo
- Criterios de aceptacion
- Scope

## Documentation Structure

- `docs/features/` stores feature plans that enter the implementation workflow.
- `docs/backlog/` stores product backlog, prioritization and future roadmap items.
- `docs/discoveries/` stores legacy discovery documents created before the single-document rule.

## Plan Storage Convention

Every new feature plan must:
- Be stored under `docs/features/`
- Be enumerated with a sequential numeric prefix
- Use the feature plan template defined in this file
- Be self-contained: analysis, proposal and implementation direction live in the same document

Recommended naming format:
- `docs/features/001-nombre-corto.md`
- `docs/features/002-nombre-corto.md`

The numbering is part of the workflow and should advance with each new plan.

Implementation branch naming:
- `docs/features/001-nombre-corto.md` -> branch `001-nombre-corto`
- `docs/features/002-nombre-corto.md` -> branch `002-nombre-corto`

## Engineering Principles

All implementation work should stay aligned with both Effective Go and practical SOLID-style design adapted to Go.

### Effective Go baseline

- Prefer small packages with clear ownership and stable vocabulary.
- Keep APIs simple, explicit and idiomatic; avoid unnecessary abstraction layers.
- Use composition over inheritance-style patterns or deep object hierarchies.
- Keep naming short, precise and consistent with the domain.
- Let behavior live close to the types or packages that own it.
- Avoid speculative generalization; extract only when a capability is demonstrably reusable.
- Keep the zero value useful when it makes sense.
- Return concrete errors with clear messages and wrap them when crossing boundaries.

### SOLID applied pragmatically in Go

- Single Responsibility: each package, type or function should have one clear reason to change.
- Open/Closed: extend behavior through stable seams, policies, interfaces or data, instead of editing unrelated code paths repeatedly.
- Liskov Substitution: interface implementations must preserve the contract expected by their consumers.
- Interface Segregation: prefer small, focused interfaces owned by the consuming package.
- Dependency Inversion: high-level flows should depend on abstractions and contracts, not concrete edge adapters or hardcoded infrastructure.

### Code health expectations

- Avoid god objects, oversized files and packages with mixed responsibilities.
- Avoid magic strings and magic numbers when they represent domain policy, balance, UX contract or shared defaults.
- Prefer typed catalogs, named constants or explicit policy objects when a value is reused or semantically important.
- Keep runtime, content/configuration and delivery edges clearly separated.
- Refactor opportunistically when new work reveals repeated logic, unstable boundaries or rising coupling.

## Backlog Convention

- Product backlog items live in `docs/backlog/`.
- Use `docs/backlog/TEMPLATE.md` for new entries.
- Each backlog item should use an identifier `BL-###`.
- When a backlog item becomes implementation-ready, convert it into a feature plan under `docs/features/`.

## Single-Document Rule

- Do not create a separate discovery document for a feature that already has a plan.
- The feature plan itself must contain the analysis needed to justify the implementation proposal.
- A plan must not defer its core organization proposal to a later discovery document.
- `docs/discoveries/` remains as historical context only unless the user explicitly asks to preserve or migrate legacy material.

## Testing Conventions

- Test cases must use explicit descriptive titles.
- There should be one test function per production function under test.
- When a production function has multiple relevant result paths, prefer table-driven tests with titled cases.
- Use mocks to isolate the unit under test and avoid coupling the test to collaborator logic.
- Keep tests declarative and avoid conditional logic inside test flows.
