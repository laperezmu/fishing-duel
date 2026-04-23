# internal/playermoves

Controla las barajas de decision disponibles para cada movimiento del jugador.

## Responsabilidad

- Definir la configuracion base de barajas por color y su recuperacion.
- Inicializar el estado de recursos del jugador.
- Validar si un movimiento puede jugarse en la ronda actual.
- Exponer la carta superior visible de cada color.
- Consumir la carta elegida, moverla al descarte y programar la recuperacion futura cuando la baraja se agota.

## Regla de arquitectura

- Este paquete aplica la mecanica de barajas de decision del jugador.
- No decide el resultado del combate ni el avance del pez en el track.
