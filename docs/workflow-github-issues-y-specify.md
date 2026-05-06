# Workflow de GitHub Issues y Specify

Este documento define el flujo activo del proyecto. GitHub Issues gestiona el backlog operativo y `specify` define, planifica y desglosa el trabajo antes de implementarlo.

## Fuente de verdad

- El backlog vivo se gestiona en GitHub Issues.
- La definicion funcional y tecnica vive en `specs/` y `.specify/`.
- El historial previo del backlog queda en GitHub y en el historial de git, no en documentos paralelos dentro del repo.

## Estructura de tracking

### Milestones

Cada milestone representa una iniciativa activa del roadmap:

- `Run MVP Foundation`
- `Zone Progression`
- `Fish Data-Driven Expansion`
- `Build and Services`
- `Economy and Meta Boundary`
- `Architecture and Delivery`
- `Graphical Client Prototype`

### Numeracion

Los items nuevos usan una numeracion corta y ordenada por milestone:

- `RMF-##` para `Run MVP Foundation`
- `ZP-##` para `Zone Progression`
- `FDE-##` para `Fish Data-Driven Expansion`
- `BS-##` para `Build and Services`
- `EMB-##` para `Economy and Meta Boundary`
- `AD-##` para `Architecture and Delivery`
- `GCP-##` para `Graphical Client Prototype`

### Labels minimas

Cada issue nueva debe tener como minimo:

- una label `stream:*`
- una label `area:*`
- una label `type:*`
- una label `priority:*`
- una label `status:*`

### Estados recomendados

- `status:ready-for-spec`: item curado y listo para crear spec
- `status:in-spec`: la spec se esta escribiendo o aclarando
- `status:planned`: existe `plan.md`
- `status:ready-for-delivery`: existe `tasks.md`
- `status:blocked`: depende de otra decision o item

## Plantilla de issue

La plantilla oficial vive en `.github/ISSUE_TEMPLATE/backlog-item.yml`.

Cada issue debe capturar:

- ID ordenado por milestone
- stream
- area
- tipo de trabajo
- objetivo
- resultado esperado
- dependencias
- prioridad
- estado de `specify`
- contexto relevante del repo

## Flujo operativo

### 1. Captura

- Toda idea no trivial nace como GitHub Issue.
- Se asignan milestone, labels y ID nuevo, por ejemplo `[RMF-06] ...`.
- El estado inicial normal es `status:ready-for-spec`.

### 2. Curado

- Se valida que el item este bien ubicado por milestone, area, tipo y prioridad.
- Si la idea es ambigua, primero se aclara en la issue antes de empezar la spec.

### 3. Spec

- Cuando el item entra en foco, se ejecuta `/speckit.specify`.
- La issue pasa a `status:in-spec`.
- La issue queda como vista ejecutiva; el detalle vive en `specs/.../spec.md`.

### 4. Clarify

- Si quedan dudas relevantes, se ejecuta `/speckit.clarify` antes de planificar.
- No se planifica trabajo no trivial con ambiguedades materiales abiertas.

### 5. Plan

- Con la spec clara, se ejecuta `/speckit.plan`.
- La issue pasa a `status:planned`.
- El enfoque de implementacion, validacion y contratos queda en `plan.md` y artefactos asociados.

### 6. Tasks

- Luego se ejecuta `/speckit.tasks`.
- Cuando existe `tasks.md`, la issue pasa a `status:ready-for-delivery`.
- Si hace falta granularidad extra, se pueden crear sub-issues tecnicas, pero `tasks.md` sigue siendo la referencia primaria de ejecucion.

### 7. Implementacion

- Se implementa desde `spec.md`, `plan.md` y `tasks.md`.
- Los cambios tecnicos se integran via branch, commits y PR.

### 8. Validacion

- Antes de abrir o actualizar PR, se ejecutan las validaciones definidas en el plan.
- Como base, cuando aplique: `go test ./...` y `golangci-lint run`.

### 9. Cierre

- Cuando el alcance queda mergeado, se cierra la issue.
- El historial queda en la issue, la spec, la PR y los commits, sin duplicar backlog dentro de `docs/`.

## Reglas de mantenimiento

- Ningun item no trivial debe pasar directo a implementacion sin issue.
- Ningun item no trivial debe pasar a implementacion sin artefactos de `specify`.
- No se deben volver a crear backlog docs activos dentro de `docs/`.
- GitHub Issues es el tablero operativo; `specs/` y `.specify/` son la fuente activa de definicion.

## Comandos utiles con `gh`

```bash
gh issue list --limit 100
gh issue view <numero>
gh issue create
gh issue edit <numero> --add-label "status:planned"
gh issue edit <numero> --remove-label "status:ready-for-spec"
gh issue edit <numero> --milestone "Run MVP Foundation"
gh issue list --milestone "Run MVP Foundation"
```

## Primer item recomendado

Para inaugurar el flujo completo conviene arrancar con una issue ya prioritaria y cercana al MVP, por ejemplo:

- `RMF-01` si quieren cerrar primero reglas globales de run
- `RMF-02` si quieren cerrar primero la economia minima

El recorrido esperado para ese primer item es:

`GitHub Issue -> /speckit.specify -> /speckit.clarify (si aplica) -> /speckit.plan -> /speckit.tasks -> implementacion -> validacion -> PR -> merge`
