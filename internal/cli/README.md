# internal/cli

Adaptador de interfaz de linea de comandos.

## Responsabilidad

- Renderizar las vistas de `internal/presentation/`.
- Leer input del usuario.
- Resolver pasos iniciales de configuracion CLI antes del loop principal, como seleccion de presets de prueba.
- Mantener el ultimo estado visual para redibujar la pantalla.
- Limpiar la terminal entre pasos para evitar spam.
- Comunicar con claridad que acciones del jugador estan disponibles o recargando.

## Archivos principales

- `ui.go`: ciclo de entrada y salida.
- `render.go`: composicion de pantallas.
- `colors.go`: ayudas ANSI.

## Si quieres mejorar la CLI

- Ajusta el layout en `render.go`.
- Ajusta paleta y acentos en `colors.go`.
- Mantén la logica del juego fuera de este paquete.

## Si quieres otra UI

- No copies este paquete; implementa una nueva UI consumiendo `internal/presentation/`.
