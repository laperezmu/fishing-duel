# Project Context

## Feature Integration Workflow

All new features start from a plan before implementation begins.

Each plan must include:
- Objective
- Acceptance criteria
- Scope of the change

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
- `docs/discoveries/` stores active discovery documents and research outputs tied to future work.
- `docs/backlog/` stores product backlog, prioritization and future roadmap items.

## Plan Storage Convention

Every new feature plan must:
- Be stored under `docs/features/`
- Be enumerated with a sequential numeric prefix
- Use the feature plan template defined in this file

Recommended naming format:
- `docs/features/001-nombre-corto.md`
- `docs/features/002-nombre-corto.md`

The numbering is part of the workflow and should advance with each new plan.

Implementation branch naming:
- `docs/features/001-nombre-corto.md` -> branch `001-nombre-corto`
- `docs/features/002-nombre-corto.md` -> branch `002-nombre-corto`

## Backlog Convention

- Product backlog items live in `docs/backlog/`.
- Active discovery documents live in `docs/discoveries/`.
- Use `docs/backlog/TEMPLATE.md` for new entries.
- Each backlog item should use an identifier `BL-###`.
- When a backlog item moves into active discovery, document that work under `docs/discoveries/`.
- When a backlog item becomes implementation-ready, convert it into a feature plan under `docs/features/`.

## Testing Conventions

- Test cases must use explicit descriptive titles.
- There should be one test function per production function under test.
- When a production function has multiple relevant result paths, prefer table-driven tests with titled cases.
- Use mocks to isolate the unit under test and avoid coupling the test to collaborator logic.
- Keep tests declarative and avoid conditional logic inside test flows.
