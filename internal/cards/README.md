# internal/cards

Define cartas y modificadores de encounter reutilizables por combatientes.

## Responsabilidad

- Describir cartas del pez y, en el futuro, cartas del jugador.
- Representar efectos configurables ligados a cartas.
- Mantener los triggers de activacion desacoplados del motor.

## Regla de arquitectura

- Este paquete solo describe datos de cartas y modificadores.
- No resuelve por si mismo reglas de ronda ni condiciones terminales.
