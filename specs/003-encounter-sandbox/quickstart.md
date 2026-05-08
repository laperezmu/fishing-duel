# Quickstart: Formalizar fishing-duel como sandbox de encounters

## Objetivo

Validar rapidamente que el sandbox formal de encounters permite configurar partidas, inspeccionar runtime, reproducir escenarios y mantener intacto `fishing-run`.

## Requisitos

1. Usar Go 1.22.x.
2. Partir del repo con los cambios del feature aplicados.
3. Contar con `go test` y `golangci-lint` disponibles.

## Validacion automatizada

```bash
go test ./internal/app ./internal/cli ./internal/presentation ./internal/content/playerprofiles ./internal/content/fishprofiles ./internal/content/watercontexts ./internal/match ./internal/game ./internal/encounter ./internal/progression
go test ./...
golangci-lint run
go test -cover ./internal/app ./internal/cli ./internal/presentation ./internal/content/playerprofiles ./internal/content/fishprofiles ./internal/content/watercontexts ./internal/match ./internal/game ./internal/encounter ./internal/progression
```

Verificar:

1. Todos los paquetes impactados pasan sus tests.
2. La cobertura agregada de los paquetes modificados por la feature alcanza al menos 90%.
3. No aparecen warnings de lint relacionados con boundaries, complejidad accidental o codigo muerto en la superficie tocada.

## Validacion manual

### 1. Sandbox guiado

```bash
go run ./cmd/fishing-duel
```

Verificar:

1. El binario se presenta como sandbox de encounters.
2. El modo guiado sigue siendo util para una prueba simple.
3. La salida muestra trazas resueltas suficientes para inspeccionar el runtime.

### 2. Sandbox manual

```bash
go run ./cmd/fishing-duel
```

Verificar:

1. El usuario puede elegir preset de pez manualmente.
2. El usuario puede fijar cartas concretas del jugador y del pez.
3. El usuario puede fijar seed y reproducir el mismo comportamiento observable.
4. Los overrides de apertura o estado soportados se reflejan sin romper el runtime.

### 3. Compatibilidad de run

```bash
go run ./cmd/fishing-run
```

Verificar:

1. El flujo de run sigue operativo.
2. No aparecen prompts ni conceptos de sandbox dentro de la run.
3. Presentation y CLI siguen siendo consistentes con el comportamiento anterior de run.

## Focos de regresion

1. Bootstrap de encounter guiado vs manual.
2. Seleccion manual de preset de pez y cartas concretas.
3. Reproducibilidad por seed.
4. Trazas de runtime y desempates visibles.
5. Aislamiento total del flujo `fishing-run`.
