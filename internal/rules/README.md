# internal/rules

Resolucion de la ronda entre la accion del jugador y la del pez.

## Responsabilidad

- Convertir dos `domain.Move` en un `domain.RoundOutcome`.

## Implementacion actual

- `ClassicEvaluator`: `Blue > Red`, `Red > Yellow`, `Yellow > Blue`.

## Como extenderlo

- Crea otro evaluador con el mismo contrato del motor.
- Inyectalo al crear `game.Engine`.

## Buenas practicas

- Mantener esta capa pura y determinista.
- Evitar dependencias con UI, mazos o estado acumulado.
