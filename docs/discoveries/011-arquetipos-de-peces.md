# Discovery: arquetipos de peces

## Objetivo

Definir una primera capa de arquetipos de pez con foco mecanico, configurables y poco dependientes de una fantasia fija.

La idea es usar los sistemas ya construidos en combate para describir perfiles faciles de mover a datos y faciles de reconfigurar si cambian la tematica o el balance.

## Principios

- Los arquetipos describen comportamiento, no lore.
- Los ids deben ser estables y tecnicos.
- Los presets CLI pueden usar nombres mas legibles, pero la configuracion base nace desde perfiles reutilizables.
- El primer slice debe ser data-friendly antes que exhaustivo.

## Arquetipos iniciales

### `baseline_cycle`
- Referencia sin presion especializada.
- Sirve como control para comparar cambios de sistema.

### `draw_tempo`
- Concentra su valor en efectos al revelar la carta.
- Ideal para probar `on_draw` y thresholds temporales.

### `horizontal_pressure`
- Empuja el encuentro hacia el escape por distancia.
- Prioriza `DistanceShift` y resultados favorables o neutros.

### `vertical_escape`
- Presiona desde la profundidad.
- Prioriza `DepthShift` y alterna hundimiento con respiracion segun el outcome.

### `surface_control`
- Juega cerca de la superficie.
- Usa efectos que reorganizan la lectura vertical del cierre.

### `deck_exhaustion`
- Busca crear ventanas cortas de resolucion alrededor del agotamiento del mazo.
- Se apoya en bonuses temporales para cierres por baraja.

### `hybrid_pressure`
- Mezcla varios ejes sin especializarse en uno solo.
- Sirve como perfil compuesto para validar el pipeline completo.

## Parametros minimos de perfil

- `id`
- `archetype_id`
- `name`
- `description`
- `details`
- `cards`
- `cards_to_remove`
- `shuffle`

## Estructura inicial implementada

- `internal/fishprofiles/profiles.go` define `ArchetypeID`, `CardPattern` y `Profile`.
- `internal/deck/presets.go` deriva los presets jugables desde esos perfiles.
- Esta capa aun vive en Go, pero ya separa configuracion de wiring manual.

## Siguiente expansion natural

- Permitir multiples perfiles por arquetipo.
- Enriquecer metadata para rareza, recompensas o tags de encounter.
- Mover la configuracion a formatos externos cuando el volumen de contenido lo justifique.
