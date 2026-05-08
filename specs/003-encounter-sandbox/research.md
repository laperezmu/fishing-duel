# Research: Formalizar fishing-duel como sandbox de encounters

## Decision 1: Introducir un `SandboxConfig` y un flujo de app separado del bootstrap de run

- Decision: Crear una configuracion y un flujo propios del sandbox en `internal/app`, en vez de expandir `EncounterBootstrapConfig` o reutilizar directamente el setup de `fishing-run`.
- Rationale: Esto preserva la separacion de responsabilidades del repositorio, evita contaminar `fishing-run` con conceptos de testing y mantiene ownership claro entre bootstrap de sandbox y loop progresivo.
- Alternatives considered:
  - Reutilizar `EncounterBootstrapConfig` y seguir agregandole flags: descartado por acoplar demasiadas responsabilidades a un bootstrap pensado para guided flow.
  - Mover toda la configuracion a `cmd/`: descartado porque duplicaria logica de coordinacion fuera de `internal/app`.

## Decision 2: Tratar la seleccion manual de pez y cartas como overlays sobre presets existentes

- Decision: La seleccion manual de preset de pez y cartas concretas se modela como overlay sobre presets y perfiles ya existentes, sin reemplazar los catalogos base.
- Rationale: Minimiza cambios, reutiliza validaciones actuales, evita sistemas paralelos de contenido y mantiene compatibilidad con el flujo guiado y con escenarios QA.
- Alternatives considered:
  - Crear un editor libre de cartas desde cero: descartado por ampliar demasiado el alcance inicial.
  - Duplicar presets especiales solo para sandbox: descartado por introducir deuda y divergencia de contenido.

## Decision 3: Mantener engine y paquetes de dominio UI-agnostic, enriqueciendo snapshots en vez de renderers ad hoc

- Decision: La evidencia adicional del sandbox debe salir de estado estructurado en `internal/match` y del runtime existente, no de logica ensamblada en CLI o presenter.
- Rationale: Alinea la feature con Effective Go, SOLID pragmatico y la constitucion; ademas mejora testabilidad y reduce riesgo de discrepancias entre runtime y salida visible.
- Alternatives considered:
  - Construir trazas textuales solo en CLI: descartado porque haria mas fragil la validacion y mezclaría UI con reglas.
  - Añadir un sistema de logging paralelo al runtime: descartado por complejidad innecesaria para el alcance.

## Decision 4: Priorizar first slice con guided/manual setup, seed y escenarios antes que overrides profundos de runtime

- Decision: El alcance inicial cubre setup formal, preset de pez manual, cartas concretas, seed reproducible, salida mas rica y escenarios reutilizables; los overrides avanzados de estado se limitan a los que puedan validarse sin romper encapsulacion.
- Rationale: Cumple el objetivo principal del issue con cambios minimos y reduce el riesgo sobre deck, progression y encounter internals.
- Alternatives considered:
  - Implementar todos los overrides de estado pedidos desde el primer slice: descartado por alto riesgo de acoplamiento transversal.
  - Omitir overrides por completo: descartado porque el spec requiere capacidad real de testing fino.

## Decision 5: Medir cobertura sobre los paquetes impactados por la feature

- Decision: El objetivo de cobertura minima del 90% se controla sobre los paquetes tocados por el sandbox, apoyado por `go test -cover` y tests unitarios nuevos o actualizados.
- Rationale: Es consistente con el pedido del usuario y evita usar un porcentaje global del repo como proxy poco util para esta feature.
- Alternatives considered:
  - Exigir 90% global del repositorio: descartado por ser desproporcionado para una feature localizada.
  - No medir cobertura formalmente: descartado por ir contra el requerimiento explicito.
