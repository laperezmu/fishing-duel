# internal/match

Define el ensamblado tactico de un duelo aislado de pesca.

## Responsabilidad

- Reunir el estado necesario para que el `engine` orqueste un encounter completo.
- Mantener separados los subestados tacticos de encuentro, mazo, jugador, ronda y ciclo de vida del duelo.
- Exponer snapshots tacticos mas estrechos para capas como `presentation` o `app` cuando no necesitan el ensamblado completo.
- Servir como frontera de runtime del combate sin absorber estado de expedicion o meta.

## Regla de arquitectura

- `match.State` representa solo un duelo aislado.
- El runtime futuro de run debe consumir este estado tactico, no mezclarse dentro de el.
- Si una capacidad no pertenece claramente al combate round a round, no debe entrar en este paquete por defecto.
