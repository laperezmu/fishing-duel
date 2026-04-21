# internal/app

Flujo de sesion desacoplado de la UI concreta.

## Responsabilidad

- Arrancar una partida.
- Pedir una jugada a la UI.
- Ejecutar una ronda en el motor.
- Entregar a la UI vistas listas para mostrar.

## Dependencias

- `game.Engine` para ejecutar la partida.
- `match.State` y `match.RoundResult` como datos compartidos del flujo.
- `presentation.Presenter` para traducir estado a vistas.
- Una implementacion de `UI` para mostrar y capturar interaccion.

## Por que existe

- Para que una UI nueva reutilice el flujo de partida sin duplicar logica.

## Como extenderlo

- Normalmente no deberias tocarlo para crear otra interfaz.
- Solo amplialo si aparece un flujo nuevo: menus previos, pausa, seleccion de pez o tutorial.
