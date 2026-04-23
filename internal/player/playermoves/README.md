# internal/player/playermoves

Controla las barajas de decision disponibles para cada movimiento del jugador.

## Responsabilidad

- Inicializar el estado de recursos del jugador.
- Validar si un movimiento puede jugarse en la ronda actual.
- Exponer la carta superior visible de cada color.
- Consumir la carta elegida, moverla al descarte y programar la recuperacion futura cuando la baraja se agota.

## Regla de arquitectura

- Este paquete aplica la mecanica de barajas de decision del jugador.
- Los presets y configuraciones reutilizables del jugador viven fuera de este paquete.
- No decide el resultado del combate ni el avance del pez en el track.
