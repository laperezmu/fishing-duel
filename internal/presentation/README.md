# internal/presentation

Traduccion de estado tecnico a contenido mostrable.

## Responsabilidad

- Convertir snapshots tacticos estrechos del duelo en view models.
- Resolver titulos, etiquetas de acciones y textos de resultado.
- Permitir cambiar idioma, tono o nomenclatura sin tocar el motor.

## Piezas principales

- `Catalog`: diccionario de textos.
- `Presenter`: transforma estado a vistas.
- `match.StatusSnapshot`, `match.RoundSnapshot`, `match.SummarySnapshot`: contratos de lectura tactica consumidos por presentacion.
- `IntroView`, `StatusView`, `RoundView`, `SummaryView`: contratos de presentacion.
- `MoveOption`: representa una accion del jugador junto con sus usos y su estado de recarga.
- Las vistas del encounter muestran tanto distancia como profundidad y eventos de superficie.

## Como extenderlo

- Crea otro `Catalog` para otro idioma o fantasia.
- Crea otro `Presenter` si necesitas una estructura de vistas distinta.

## Regla de arquitectura

- Este paquete puede conocer el dominio y snapshots tacticos del juego.
- No debe depender de una UI concreta.
