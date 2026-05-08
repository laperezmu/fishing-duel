# Research: Splash interactivo con saltos y mejoras de cana

## Decision 1: Modelar el splash como estado de runtime del encounter, no como logica embebida en CLI

- Decision: Introducir una representacion explicita del splash dentro del runtime del encounter y snapshots asociados, dejando a la CLI solo la resolucion concreta de input y timing.
- Rationale: El evento afecta estado tactico, endings y progresion, asi que su fuente de verdad debe vivir en el dominio del encounter. Esto preserva la regla constitucional de mantener paquetes de dominio UI-agnostic.
- Alternatives considered:
  - Resolver todo el splash dentro de `internal/cli`: descartado porque acopla reglas del encounter a terminal.
  - Mantener solo un evento puntual sin estado intermedio: descartado porque varios saltos y bonus por salto requieren seguimiento de progreso.

## Decision 2: Mover la decision final del splash desde aleatoriedad inmediata a una resolucion posterior orquestada por `internal/app`

- Decision: Cuando una ronda dispare splash, el runtime debe marcar una secuencia pendiente en lugar de cerrar de inmediato por RNG; luego `internal/app` pide resolverla a la UI y aplica el resultado al engine.
- Rationale: El flujo actual de `app.Session` ya coordina UI y engine. Extender esa orquestacion mantiene al engine centrado en reglas y al adaptador CLI centrado en experiencia interactiva.
- Alternatives considered:
  - Mantener un `SplashEscapeDecider` aleatorio con nueva probabilidad: descartado porque contradice el objetivo principal de habilidad.
  - Hacer que `game.Engine` lea input directamente: descartado por romper separacion de responsabilidades.

## Decision 3: Reutilizar el patron del minijuego de cast para el primer slice CLI

- Decision: Basar el primer splash interactivo CLI en un patron de barra, objetivo visible y confirmacion por `Enter`, con ventana temporal configurable y progreso por salto.
- Rationale: El repo ya tiene una interaccion de timing compatible con CLI en el cast. Reutilizar el mismo lenguaje reduce riesgo, acelera entrega y mantiene coherencia de UX.
- Alternatives considered:
  - Pedir combinaciones de teclas o input libre: descartado por mayor complejidad y peor portabilidad en terminal.
  - Crear una simulacion totalmente distinta por salto: descartado para no duplicar conceptos de UI en el primer slice.

## Decision 4: Ubicar el bonus por splash ganado en el loadout efectivo del jugador

- Decision: Incorporar la recompensa por salto ganado como capacidad de loadout/rod efectiva, accesible desde el runtime de la partida sin atar la regla a un preset concreto.
- Rationale: La spec habla de mejoras de cana. El lugar mas estable para esa capacidad hoy es `internal/player/loadout` y los presets de `rod`, evitando mezclar beneficio de build con encounter puro.
- Alternatives considered:
  - Poner el bonus en perfiles de pez: descartado porque no pertenece al jugador.
  - Hardcodear el bonus en CLI o bootstrap: descartado por no ser data-friendly ni reusable.

## Decision 5: Mantener los detonantes actuales de splash y posponer nuevas fuentes del evento

- Decision: El slice solo cambia la resolucion del splash y su configuracion, sin anadir triggers adicionales mas alla de los ya existentes por movimiento vertical sobre la superficie.
- Rationale: La spec y la issue piden preservar los detonantes actuales para acotar el refactor. Esto reduce riesgo sobre reglas de round y mantiene un vertical slice pequeno.
- Alternatives considered:
  - Revisar toda la taxonomia de triggers de splash: descartado por aumentar mucho el alcance.
  - Mover splash a una condicion terminal separada de progresion: descartado por no resolver el origen del evento actual.

## Decision 6: Mantener el contrato de validacion con tests unitarios por paquete y regression end-to-end minima

- Decision: Agregar tests unitarios en encounter, progression, app, presentation, CLI y content donde cambie el comportamiento, mas `go test ./...` y `golangci-lint run` como validacion final.
- Rationale: Sigue las practicas actuales del repo y el principio constitucional de entrega testeable.
- Alternatives considered:
  - Cubrir solo con pruebas manuales de CLI: descartado por insuficiente frente al cambio de reglas.
  - Introducir un tipo nuevo de prueba de integracion pesada: descartado para mantener el slice pequeno.
