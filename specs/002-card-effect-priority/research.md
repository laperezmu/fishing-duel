# Research: Rediseno de triggers y efectos de cartas con prioridad de resolucion

## Decision 1: Modelar prioridad a nivel de efecto y resolverla en `internal/game`

- Decision: La prioridad se define por efecto individual en el modelo de cartas y la secuencia ordenada de resolucion se construye en `internal/game`.
- Rationale: `internal/cards` ya concentra triggers y efectos, mientras `internal/game` es el punto donde hoy se combinan efectos de jugador y pez. Mantener el ordenado en engine evita mover reglas a CLI o presentation y preserva limites modulares.
- Alternatives considered:
  - Resolver prioridad en `internal/presentation` o `internal/cli`: descartado porque mezcla UI con reglas.
  - Resolver prioridad dentro de `internal/progression`: descartado porque progression hoy aplica efectos ya seleccionados y no deberia orquestar toda la secuencia inter-paquete.

## Decision 2: Sustituir bonuses legacy por un catalogo explicito de efectos soportados

- Decision: `CaptureDistanceBonus` y `SurfaceDepthBonus` se migran a un catalogo nuevo de efectos explicitos, incluyendo variantes por entidad y por color cuando el comportamiento visible sea similar pero la semantica cambie.
- Rationale: El spec exige deprecacion formal y cierre del gap entre contenido actual y lista objetivo. Mantener bonuses legacy como concepto principal perpetuaria ambiguedad entre thresholds temporales, fatiga y movimiento real.
- Alternatives considered:
  - Mantener los bonuses legacy y solo agregar prioridad: descartado porque no satisface el objetivo del issue.
  - Reemplazar todo por una unica estructura generica sin distinguir variantes: descartado porque esconderia restricciones importantes como pez vs jugador o color concreto.

## Decision 3: Mantener la aplicacion de reglas en paquetes de runtime y usar snapshots solo para observabilidad

- Decision: `internal/encounter` y `internal/progression` siguen aplicando efectos sobre estado y thresholds; `internal/match` y `internal/presentation` solo exponen informacion suficiente para inspeccion, tests y compatibilidad CLI.
- Rationale: Esto respeta la constitucion del repo: dominio UI-agnostic, app coordinando flujos y presentation traduciendo estado. Tambien minimiza el cambio sobre `cmd/` y adaptadores CLI.
- Alternatives considered:
  - Mover parte de la resolucion al presenter para mostrar prioridades: descartado porque rompe boundaries.
  - No exponer informacion nueva en snapshots: descartado porque dificultaria validar el orden determinista y el desempate.

## Decision 4: Migrar contenido en las dos fuentes actuales sin abrir un sistema nuevo de configuracion

- Decision: El slice contempla migrar tanto presets hardcoded del jugador como catalogos JSON de pez al contrato nuevo, manteniendo las fuentes actuales de contenido.
- Rationale: Hoy el contenido vive en dos formas distintas (`internal/content/playerprofiles` y `internal/content/fishprofiles`) y ambas ya alimentan el runtime. Cambiar ambas en este slice evita inconsistencias funcionales.
- Alternatives considered:
  - Migrar solo pez y dejar jugador para otro item: descartado porque el spec cubre cartas de ambos lados.
  - Introducir un nuevo sistema externo de contenido durante esta tarea: descartado por ampliar innecesariamente el alcance.

## Decision 5: Cubrir la regresion con tests unitarios de reglas, contenido y presentation adaptada

- Decision: La validacion automatizada se centra en tests de `internal/cards`, `internal/game`, `internal/progression`, `internal/encounter`, `internal/match`, `internal/content/*`, y pruebas de compatibilidad en `internal/presentation`, `internal/app` y `internal/cli`.
- Rationale: El mayor riesgo es alterar el orden real de resolucion, thresholds y finales de encounter. La cobertura distribuida permite detectar quiebres de semantica sin depender de smoke tests manuales solamente.
- Alternatives considered:
  - Confiar en `go test ./...` sin ampliar casos dirigidos: descartado por riesgo de dejar huecos en el nuevo orden de prioridad.
  - Validar solo por CLI manualmente: descartado porque no ofrece suficiente aislamiento de reglas.
