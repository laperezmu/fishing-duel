# Plan de feature: presets-de-baraja-para-cli

## Objetivo

Introducir una seleccion inicial de presets de baraja del pez en la experiencia CLI para poder probar de forma rapida las mecanicas nuevas del encounter sin modificar codigo cada vez.

El resultado esperado es un flujo de arranque donde, antes de crear la partida, el jugador elige uno de varios presets de baraja del pez, confirma la seleccion y luego entra al loop habitual del juego.

## Criterios de aceptacion

- Al iniciar el ejecutable CLI se muestra una lista de presets de baraja del pez.
- El usuario puede elegir un preset por indice y confirmar la seleccion.
- Si el usuario cancela la confirmacion, vuelve a la lista de presets sin romper el flujo.
- Despues de confirmar un preset, la partida continua con el flujo actual de intro, rondas y game over.
- Existen varios presets pensados para probar mecanicas recientes como `on_draw`, modificadores post-outcome, profundidad y thresholds temporales.
- `go test ./...` pasa.
- `golangci-lint run` pasa.

## Scope

### Incluye
- Definir una lista estable de presets de baraja del pez para testing manual desde CLI.
- Modelar metadatos de preset suficientes para mostrarlos en la seleccion inicial.
- Añadir una pantalla CLI para elegir y confirmar preset.
- Integrar el preset elegido en el bootstrap de `cmd/fishing-duel/main.go`.
- Cubrir con tests la seleccion, la confirmacion y la construccion de la baraja elegida.

### No incluye
- Crear aun una UI avanzada de inspeccion completa de cartas con ilustraciones.
- Introducir configuracion persistente del ultimo preset usado.
- Convertir esta seleccion inicial en un sistema general de mods o contenido externo.

## Propuesta de implementacion

- Crear en `internal/deck/` una lista de presets reutilizable por el CLI.
- Mantener un preset clasico sin efectos y varios presets orientados a probar features recientes.
- Reutilizar `cli.UI` para resolver la seleccion inicial y la confirmacion antes de instanciar `game.Engine`.
- Mantener el resto del flujo de `app.Session` sin cambios para que la seleccion viva solo en el bootstrap.
- Asegurar que los presets copian sus cartas al construir la baraja para no compartir slices mutables.

## Validacion

- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.
- Cubrir con tests al menos:
  - seleccion valida de preset
  - cancelacion de confirmacion y vuelta a la lista
  - rechazo de entradas no validas
  - construccion correcta de presets de prueba

## Riesgos o decisiones abiertas

- Si los presets mezclan demasiadas mecanicas a la vez, pueden dejar de ser utiles como herramientas de prueba.
- Habra que equilibrar entre presets deterministas para debugging y presets mas cercanos al flujo normal del juego.
