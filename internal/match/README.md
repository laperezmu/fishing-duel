# internal/match

Modelo compartido del estado acumulado de una partida.

## Responsabilidad

- Definir el snapshot persistente del encuentro.
- Exponer el resultado de una ronda como dato compartido entre capas.
- Servir como contrato comun entre motor, progresion, finales, presentacion y app.

## Piezas principales

- `State`: estado completo de la partida.
- `RoundState`: estado temporal del round para efectos que no deben persistir mas alla de una ronda.
- `DeckState`: snapshot tecnico del mazo del pez.
- `PlayerRig`: capacidades operativas actuales del jugador.
- `PlayerMoveResources`: estado acumulado de las barajas de decision del jugador.
- `Stats`: victorias, derrotas y empates acumulados.
- `ResolvedRound`: resultado base de una ronda antes de derivar la vista final.
- `ResolvedRound.PlayerCard`: carta concreta del jugador usada en la ronda.
- `ResolvedRound.DrawEffects`: efectos activados al revelar la carta activa del round.
- `ResolvedRound.OutcomeEffects`: efectos activados despues de resolver el outcome base.
- `RoundResult`: resultado de una ronda ya ejecutada.

## Regla de arquitectura

- Este paquete solo contiene datos compartidos.
- No debe contener logica de orquestacion, UI o reglas de combate.
