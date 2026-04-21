# internal/rules

Resolucion de la ronda entre la accion del jugador y la del pez.

## Responsabilidad

- Convertir dos `domain.Move` en un `domain.RoundOutcome`.

## Implementacion actual

- `ClassicEvaluator`: resuelve la regla base `Blue > Red`, `Red > Yellow`, `Yellow > Blue`.
- `Condition`: permite anadir condiciones de combate antes de aplicar la regla base.
- `TieAdvantageCondition`: hace que el pez gane un empate si el color pertenece a su perfil de combate.

## Como extenderlo

- Si la mecanica reemplaza por completo la resolucion, crea otro evaluador con el mismo contrato del motor.
- Si la mecanica complementa la resolucion actual, implementa otra `Condition` y anadela al construir `ClassicEvaluator`.
- Inyecta el evaluador resultante al crear `game.Engine`.

## Buenas practicas

- Mantener esta capa pura y determinista.
- Evitar dependencias con UI, mazos o estado acumulado.
