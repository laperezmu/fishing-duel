# internal/domain

Tipos canonicos del juego.

## Contiene

- `Move`: las tres jugadas internas (`Blue`, `Red`, `Yellow`).
- `RoundOutcome`: resultado abstracto de una ronda.

## Regla de uso

- Este paquete no debe conocer textos tematicos ni detalles de UI.
- Usa estos tipos como lenguaje comun entre paquetes.

## Cuando extenderlo

- Si agregas nuevos colores o nuevas categorias de resultado.
- Si necesitas nuevos enums compartidos por varias capas.

## Cuando no tocarlo

- Si solo cambias nombres visibles como `Tirar` o `Embestir`; eso pertenece a `internal/presentation/`.
