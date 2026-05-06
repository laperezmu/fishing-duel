# Workflow de GitHub Issues y Specify

Este documento define como gestionar el backlog activo del proyecto usando GitHub Issues como capa operativa y `specify` como sistema de definicion y planificacion.

## Fuente de verdad

- El backlog vivo se gestiona en GitHub Issues.
- La definicion detallada del trabajo vive en `specs/` y `.specify/`.
- `docs/backlog/` queda como contexto historico y mapa de migracion, no como tablero operativo diario.

## Estructura de tracking

### Milestones

Cada milestone representa una iniciativa relevante del roadmap:

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

La numeracion legacy `BL-*` y `GUI-*` no debe reutilizarse para trabajo nuevo.

### Labels

Las issues deben usar, como minimo:

- una label `stream:*`
- una label `area:*`
- una label `type:*`
- una label `priority:*`
- una label `status:*`

## Estados recomendados

- `status:ready-for-spec`: item curado y listo para pasar a `/speckit.specify`
- `status:in-spec`: la spec se esta escribiendo o aclarando
- `status:planned`: existe `plan.md`
- `status:ready-for-delivery`: existe `tasks.md` y el item ya puede implementarse
- `status:blocked`: depende de otra decision, spec o issue

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
- contexto historico o rutas del repo relevantes

## Flujo operativo

### 1. Crear o migrar una issue

- Crear la issue usando la plantilla `Backlog item`
- Asignar milestone y labels
- Usar el prefijo de titulo con el nuevo ID, por ejemplo: `[RMF-01] Definir recursos globales de run`

### 2. Convertirla en spec

- Cuando la issue quede priorizada, iniciar `/speckit.specify`
- Referenciar la issue y su ID en la spec cuando sea util
- Mantener la issue como vista ejecutiva; no duplicar toda la spec dentro del body

### 3. Planificar

- Ejecutar `/speckit.plan`
- Cambiar la label a `status:planned`

### 4. Desglosar a tareas

- Ejecutar `/speckit.tasks`
- Cuando exista `tasks.md`, mover la issue a `status:ready-for-delivery`

### 5. Implementar

- Implementar desde la spec y las tasks
- Si el trabajo requiere seguimiento mas fino, crear sub-issues tecnicas o usar el propio `tasks.md`

## Criterios de migracion desde backlog legacy

- Migrar solo items abiertos
- No migrar `done` ni `cancelled` como issues activas
- Conservar el ID legacy solo como referencia historica en `Source context` o en notas de migracion
- Usar el mapa definido en `docs/backlog/004-mapa-migracion-github-issues.md`

## Convenciones con `gh`

Comandos utiles:

```bash
gh issue list --limit 100
gh issue view <numero>
gh issue create
gh issue edit <numero> --add-label "status:planned"
gh issue edit <numero> --milestone "Run MVP Foundation"
gh issue list --milestone "Run MVP Foundation"
```

## Regla de mantenimiento

- Ningun item no trivial deberia pasar directo a implementacion sin issue y sin artefactos de `specify`
- Si una idea aun no esta suficientemente clara, se mantiene como issue en `status:ready-for-spec`
- Si cambia el roadmap, se reordena por milestone y se asigna un nuevo ID solo si el item realmente cambia de iniciativa
