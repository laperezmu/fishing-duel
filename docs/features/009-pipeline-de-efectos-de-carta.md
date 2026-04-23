# Plan de feature: pipeline-de-efectos-de-carta

## Objetivo

Completar la siguiente iteracion de `BL-005` llevando el sistema de efectos de carta desde contratos compartidos a un pipeline real de resolucion por fases, donde `on_draw` exista como fase funcional del motor y no solo como una idea de la taxonomia.

El resultado esperado es un flujo de round donde una carta pueda emitir efectos al revelarse y despues del outcome, manteniendo el sistema agnostico al owner de la carta, sin mutar el `rig` y sin dejar efectos temporales vivos mas alla del round en que nacen.

Para esta feature se asume ya una restriccion de producto: las `playerCards` no se eligen individualmente durante el combate. El jugador parte con una baraja por decision, cada una asociada a un unico color, y de base dispone de 3 cartas iguales por color sin efectos.

## Criterios de aceptacion

- El motor ejecuta una fase real de `on_draw` cuando entra la carta activa del pez al round.
- Los efectos de carta se resuelven mediante un pipeline explicito por fases en lugar de depender de filtrados ad hoc dentro de progresion.
- El pipeline sigue siendo agnostico al owner de la carta y reutilizable por `FishCard` y futuras `playerCards`.
- Existe al menos un efecto real de carta que use `on_draw` y altere el round de forma verificable.
- Los efectos temporales de thresholds siguen viviendo solo en estado temporal del round y nunca mutan el `rig`.
- El pipeline deja preparada la integracion futura de `playerCards` respetando que el jugador roba o consume desde barajas por color, no desde una mano de seleccion libre durante combate.
- `go test ./...` pasa.
- `golangci-lint run` pasa.

## Scope

### Incluye
- Introducir una fase real de `on_draw` en el flujo del motor.
- Definir un pipeline minimo de fases para efectos de carta, al menos draw y post-outcome.
- Ajustar el contexto de round para que pueda transportar efectos habilitados por fase y owner.
- Respetar que las futuras `playerCards` nacen desde barajas por color y no desde una seleccion libre carta a carta durante el combate.
- Mantener el desacople entre `rules`, `progression`, `encounter` y `endings` segun lo definido en el discovery `008`.
- Implementar al menos un caso real de efecto `on_draw`, preferentemente sobre estado temporal del round.
- Preservar los efectos actuales ligados al outcome y su compatibilidad con `FishCard`.
- Cubrir con tests el nuevo orden de resolucion y la expiracion de estado temporal al cerrar el round.

### No incluye
- Implementar aun el sistema completo de barajas de `playerCards` por color dentro del loop jugable.
- Introducir multiples familias nuevas de efectos de carta en la misma iteracion.
- Reemplazar `rules` por una capa nueva anterior mientras los hooks actuales escalen correctamente.
- Redisenar la UI o presentar visualmente cartas completas con ilustracion.

## Propuesta de implementacion

- Extender `internal/cards/` para que los efectos puedan declararse por fase de forma explicita, manteniendo triggers relativos al owner.
- Evolucionar `internal/game/engine.go` para separar la resolucion del round en pasos claros: draw/revelacion, efectos `on_draw`, resolucion base de combate, hooks de outcome en `rules`, efectos post-outcome, progresion, eventos derivados y finales.
- Ajustar `internal/match/` para que el contexto temporal del round pueda transportar la carta emisora, los efectos habilitados y thresholds temporales solo mientras dura la ronda.
- Mantener compatibilidad con un modelo futuro donde el jugador tenga 3 cartas base por color sin efectos y donde las mejoras introduzcan variaciones sobre esas barajas por color.
- Mantener `internal/progression/` centrado en aplicar impactos espaciales y otros payloads compatibles con esa capa, en lugar de decidir por si mismo que trigger corre.
- Mantener `internal/endings/` como consumidor del estado ya resuelto por el pipeline, respetando que los bonuses temporales del round expiran al finalizar.
- Añadir un ejemplo real de carta de pez con efecto `on_draw` que permita validar la fase nueva sin requerir aun `playerCards` jugables.
- Documentar cualquier decision adicional derivada de la implementacion en `docs/discoveries/008-taxonomia-de-efectos-de-carta.md` si cambia o se precisa algun punto del discovery.

## Validacion

- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.
- Cubrir con tests al menos:
  - activacion real de `on_draw`
  - orden de fases del pipeline
  - preservacion de efectos `on_owner_win` y `on_owner_lose`
  - expiracion del estado temporal del round tras cerrar la ronda
  - preservacion de la captura y escape actuales cuando no hay efectos `on_draw`

## Riesgos o decisiones abiertas

- Si el pipeline introduce demasiada infraestructura antes de tener `playerCards` reales, puede quedar mas abstracto de lo necesario.
- Habra que decidir que forma concreta toma el ejemplo inicial de `on_draw` para que sea util y no distorsione el balance del encounter.
- Si aparecen efectos que necesitan mutar outcome base antes de `rules`, habra que reevaluar si los hooks actuales siguen siendo suficientes o si la separacion por fases necesita crecer.
