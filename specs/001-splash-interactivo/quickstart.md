# Quickstart: Splash interactivo con saltos y mejoras de cana

## Objetivo

Validar rapidamente el slice de splash interactivo en los dos ejecutables CLI existentes.

## Preparacion

1. Tener Go instalado.
2. Trabajar desde la raiz del repo.
3. Usar la branch `001-splash-interactivo-cana`.

## Validacion automatizada

```bash
go test ./internal/encounter ./internal/progression ./internal/game ./internal/app ./internal/presentation ./internal/cli ./internal/player/loadout ./internal/content/rodpresets ./internal/content/anglerprofiles
go test ./...
$(go env GOPATH)/bin/golangci-lint run
```

## Validacion manual del sandbox

```bash
go run ./cmd/fishing-duel
```

Escenarios a comprobar:

1. Un splash de un solo salto puede resolverse con exito y el encounter continua.
2. Un splash falla por timeout o fallo de input y termina en `splash_escape`.
3. Un pez configurado con varios saltos exige completar toda la secuencia.
4. Una cana con bonus de splash acerca al pez tras cada salto ganado.

## Validacion manual de la run

```bash
go run ./cmd/fishing-run
```

Escenarios a comprobar:

1. La run sigue entrando y saliendo de encounters sin romper el flujo general.
2. Los resumenes de ronda y fin de encounter siguen mostrando correctamente un escape por chapoteo.
3. El juego sigue siendo completamente operable desde CLI.
