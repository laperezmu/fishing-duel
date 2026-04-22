# internal/progression

Efectos de una ronda sobre el estado persistente del juego.

## Responsabilidad

- Aplicar el resultado de la ronda al `match.State`.
- Mover la distancia del pez.
- Aplicar movimiento vertical y modificadores de carta del encounter.
- Actualizar estadisticas de victorias y empates.

## Implementacion actual

- `TrackPolicy`: acerca o aleja al pez, aplica modificadores de carta y resuelve eventos de superficie.

## Como extenderlo

- Crea otra politica para encuentros especiales.
- Ejemplos: peces que retroceden dos pasos, peces que ignoran empates, bonus por combos, cambios de profundidad o eventos acuaticos especiales.

## Regla de arquitectura

- Esta capa no debe decidir por si sola si la partida termino; eso le corresponde a `internal/endings/`.
