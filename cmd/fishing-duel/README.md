# cmd/fishing-duel

Punto de entrada del ejecutable.

## Responsabilidad

- Crear el mazo inicial del pez.
- Crear la configuracion del encuentro.
- Construir el motor (`game.Engine`).
- Elegir la capa de presentacion.
- Elegir la UI concreta.

## Que no debe hacer

- No debe contener reglas del juego.
- No debe contener logica de mazo.
- No debe formatear estados de juego complejos.

## Si quieres cambiar la experiencia inicial

- Sustituye el `Catalog` del paquete `internal/presentation/`.
- Sustituye la UI de `internal/cli/` por otra interfaz.
- Sustituye las politicas de `rules`, `progression` o `endings` al crear el motor.
