# internal/playermoves

Controla los usos disponibles de cada movimiento del jugador.

## Responsabilidad

- Definir la configuracion base de usos y recarga.
- Inicializar el estado de recursos del jugador.
- Validar si un movimiento puede jugarse en la ronda actual.
- Consumir usos y programar la recarga futura cuando un movimiento se agota.

## Regla de arquitectura

- Este paquete aplica la mecanica de recursos del jugador.
- No decide el resultado del combate ni el avance del pez en el track.
