# internal/match

Define el ensamblado tactico de un duelo aislado de pesca.

## Responsabilidad

- Reunir el estado necesario para que el `engine` orqueste un encounter completo.
- Mantener separados los subestados tacticos de encuentro, mazo, jugador, ronda y ciclo de vida del duelo.
- Exponer snapshots tacticos mas estrechos para capas como `presentation` o `app` cuando no necesitan el ensamblado completo.
- Exponer contratos de mutacion mas estrechos como `ProgressionState`, `EndingState` y `PlayerMoveRuntime` para que los colaboradores del duelo no dependan del ensamblado completo.
- Servir como frontera de runtime del combate sin absorber estado de expedicion o meta.

## Contratos principales

- `State`: ensamblado tactico completo del duelo para orquestacion.
- `StatusSnapshot`, `RoundSnapshot`, `SummarySnapshot`: lecturas listas para presentacion y app.
- `RoundResult`: resultado de ronda con snapshots tacticos ya derivados, sin reexportar el runtime completo.
- `ProgressionState`: acceso minimo para actualizar track, stats y encounter durante la progresion.
- `EndingState`: acceso minimo para resolver condiciones terminales del duelo.
- `PlayerMoveRuntime`: acceso minimo para preparar, validar y consumir recursos de movimiento del jugador.

## Regla de arquitectura

- `match.State` representa solo un duelo aislado.
- El runtime futuro de run debe consumir este estado tactico, no mezclarse dentro de el.
- Si una capacidad no pertenece claramente al combate round a round, no debe entrar en este paquete por defecto.
