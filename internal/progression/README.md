# internal/progression

Efectos de una ronda sobre el estado persistente del juego.

## Responsabilidad

- Aplicar el resultado de la ronda al `game.State`.
- Mover la distancia del pez.
- Actualizar estadisticas de victorias y empates.

## Implementacion actual

- `TrackPolicy`: acerca al pez cuando gana el jugador y lo aleja cuando gana el pez.

## Como extenderlo

- Crea otra politica para encuentros especiales.
- Ejemplos: peces que retroceden dos pasos, peces que ignoran empates, bonus por combos.

## Regla de arquitectura

- Esta capa no debe decidir por si sola si la partida termino; eso le corresponde a `internal/endings/`.
