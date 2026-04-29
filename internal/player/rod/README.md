# internal/player/rod

Define la `rod` base del jugador y sus limites estructurales de apertura y track.

## Responsabilidad

- Describir hasta donde puede abrir la pesca la `rod` del jugador.
- Describir hasta donde puede sostenerse el track del duelo antes del escape.
- Servir como base configurable para futuros aditamentos y loadouts.

## Regla de arquitectura

- Este paquete modela la `rod` base del jugador, no el comportamiento del pez ni el setup completo.
- Sus valores pueden ampliarse en futuras iteraciones por aditamentos, items o build.
