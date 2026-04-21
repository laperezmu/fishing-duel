# internal/presentation

Traduccion de estado tecnico a contenido mostrable.

## Responsabilidad

- Convertir `match.State` y `match.RoundResult` en view models.
- Resolver titulos, etiquetas de acciones y textos de resultado.
- Permitir cambiar idioma, tono o nomenclatura sin tocar el motor.

## Piezas principales

- `Catalog`: diccionario de textos.
- `Presenter`: transforma estado a vistas.
- `IntroView`, `StatusView`, `RoundView`, `SummaryView`: contratos de presentacion.

## Como extenderlo

- Crea otro `Catalog` para otro idioma o fantasia.
- Crea otro `Presenter` si necesitas una estructura de vistas distinta.

## Regla de arquitectura

- Este paquete puede conocer el dominio y el estado del juego.
- No debe depender de una UI concreta.
