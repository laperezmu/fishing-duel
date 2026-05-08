# Quickstart: Rediseno de triggers y efectos de cartas con prioridad de resolucion

## Objetivo

Validar rapidamente el slice de prioridad de efectos, migracion de contenido y compatibilidad CLI.

## Requisitos

1. Usar Go 1.22.x.
2. Tener dependencias instaladas con `go test` listo para ejecutarse.
3. Partir del repo con los cambios del feature aplicados.

## Validacion automatizada

```bash
go test ./internal/cards ./internal/game ./internal/progression ./internal/encounter ./internal/match ./internal/content/playerprofiles ./internal/content/fishprofiles ./internal/endings ./internal/presentation ./internal/app ./internal/cli
go test ./...
golangci-lint run
```

## Validacion manual

### 1. Compatibilidad del duel CLI

```bash
go run ./cmd/fishing-duel
```

Verificar:

1. Un round con efectos simultaneos mantiene orden determinista.
2. Si dos efectos tienen la misma prioridad, el del pez resuelve primero.
3. Un reshuffle o cambio de descarte se refleja sin romper la lectura del flujo.
4. La informacion mostrada sigue siendo comprensible sin que CLI implemente reglas propias.
5. El resumen del ultimo lance muestra la traza de efectos resueltos en el orden aplicado.

### 2. Compatibilidad del run CLI

```bash
go run ./cmd/fishing-run
```

Verificar:

1. El flujo principal de run sigue operativo.
2. Los encuentros con cartas migradas no rompen progression ni finales.
3. Los snapshots o resumentes de ronda siguen siendo consistentes.
4. La compatibilidad de presenter y CLI se mantiene cuando la ronda incluye efectos migrados y trazas resueltas.

## Focos de regresion

1. Bonuses legacy de thresholds reemplazados por efectos nuevos.
2. Triggers generales vs triggers por color.
3. Fatiga del pez vs fatiga del jugador por color.
4. Endings por captura, splash o agotamiento tras el nuevo orden de resolucion.
