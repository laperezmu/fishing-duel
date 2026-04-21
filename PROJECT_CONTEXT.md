# Project Context

## Feature Integration Workflow

All new features start from a plan before implementation begins.

Each plan must include:
- Objective
- Acceptance criteria
- Scope of the change

Implementation should not start without at least this minimum planning step.

## Delivery Process

After implementation is complete:
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
