# internal/cards

Define cartas y modificadores de encounter reutilizables por combatientes.

## Responsabilidad

- Describir cartas del pez y cartas del jugador.
- Representar efectos configurables ligados a cartas mediante contratos compartidos.
- Mantener los triggers de activacion desacoplados del motor.

## Regla de arquitectura

- Este paquete solo describe datos de cartas y modificadores.
- Puede tener tipos de carta distintos por owner mientras comparte contratos comunes de efectos.
- No resuelve por si mismo reglas de ronda ni condiciones terminales.
