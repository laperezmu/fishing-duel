# internal/content/playerprofiles

Define presets y configuraciones reutilizables de barajas del jugador fuera del runtime de `playermoves`.

## Responsabilidad

- Exponer presets de barajas del jugador para bootstrap y testing manual.
- Mantener configuracion reusable separada de la mecanica de consumo y recuperacion.
- Servir como primera capa data-friendly para futuras expansiones del jugador.

## Limites

- No valida ni consume cartas durante la ronda.
- No resuelve recuperacion, recarga ni disponibilidad en tiempo de ejecucion.
