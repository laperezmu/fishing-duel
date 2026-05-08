# Implementation Plan: Rediseno de triggers y efectos de cartas con prioridad de resolucion

**Branch**: `002-card-effect-priority` | **Date**: 2026-05-07 | **Spec**: `specs/002-card-effect-priority/spec.md`
**Input**: Feature specification from `specs/002-card-effect-priority/spec.md`

## Summary

Implementar un rediseño del sistema de efectos de cartas que sustituya el orden implicito actual por una resolucion determinista basada en prioridad por efecto, manteniendo la arquitectura modular del proyecto y la compatibilidad con CLI. El enfoque concentra el cambio en el modelo de cartas, la construccion del plan de resolucion en engine, la aplicacion de efectos en encounter/progression, la migracion del contenido de pez y jugador, y la exposicion de resultados suficientes para tests y presentacion sin mover reglas de negocio a la UI.

## Out of Scope

- Crear sistemas nuevos de economia, meta-progresion o progression global fuera del dominio de cartas y encounters.
- Reemplazar la arquitectura modular actual o mover logica de runtime a `cmd/`, `internal/cli` o `internal/presentation`.
- Redisenar los flujos principales de CLI mas alla de la compatibilidad necesaria para mostrar el resultado del contrato nuevo.

## Technical Context

**Language/Version**: Go 1.22.0  
**Primary Dependencies**: standard library, existing repo packages, `testify` ya presente en el repo  
**Storage**: in-memory runtime and repository files unless the spec says otherwise  
**Testing**: `go test ./...`, targeted package tests, `golangci-lint run`  
**Target Platform**: CLI on macOS/Linux/Windows  
**Project Type**: modular Go application with CLI entrypoints under `cmd/`  
**Constraints**: preserve current modular boundaries; keep domain packages UI-agnostic; mantener compatibilidad con `cmd/fishing-duel` y `cmd/fishing-run`; no mover reglas de prioridad a presentation o CLI; agregar tests unitarios para reglas y funcionalidades nuevas  
**Scale/Scope**: single feature slice for the current game/runtime architecture

## Constitution Check

- [x] Spec is the active source of truth for this feature
- [x] Scope is explicit, with out-of-scope items recorded
- [x] Planned changes preserve modular package boundaries
- [x] Validation path includes tests and lint where behavior changes materially
- [x] Risks, assumptions, and tradeoffs are documented

Post-design check:

- [x] La prioridad y resolucion de efectos permanecen en runtime/dominio, no en CLI ni presentation.
- [x] El plan conserva la separacion entre contenido, engine, progression, encounter y adaptadores de UI.
- [x] El cambio se mantiene como un slice focalizado sobre cartas y encounters, sin abrir sistemas meta ni refactors horizontales fuera del alcance del spec.

## Project Structure

### Documentation (this feature)

```text
specs/002-card-effect-priority/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
└── tasks.md
```

### Source Code (repository root)

```text
cmd/
├── fishing-duel/
└── fishing-run/

internal/
├── app/
├── cards/
├── cli/
├── content/
├── deck/
├── domain/
├── encounter/
├── endings/
├── game/
├── match/
├── player/
├── presentation/
├── progression/
└── rules/
```

**Structure Decision**: Prefer touching the smallest number of packages that can own the feature cleanly. Keep runtime, content, presentation, and bootstrap concerns separated.

## Implementation Approach

### Package Impact

- Primary packages to change: `internal/cards`, `internal/game`, `internal/progression`, `internal/encounter`, `internal/match`, `internal/content/playerprofiles`, `internal/content/fishprofiles`
- Secondary packages to review: `internal/endings`, `internal/presentation`, `internal/app`, `internal/cli`, `cmd/fishing-duel`, `cmd/fishing-run`
- Files or areas intentionally avoided: `internal/rules`, `internal/deck`, economia/meta, backlog docs fuera de `specs/`

### Execution Strategy

1. Redefinir el contrato de efectos y triggers en `internal/cards`, incorporando prioridad por efecto, variantes por entidad o color y el catalogo nuevo con deprecacion formal de bonuses legacy.
2. Rehacer la resolucion del runtime en `internal/game`, `internal/progression`, `internal/encounter` y `internal/match` para construir y aplicar una secuencia determinista de efectos sin filtrar reglas a presentation o CLI.
3. Migrar contenido y validaciones en `internal/content/*`, ajustar snapshots o presenter solo donde haga falta para mantener compatibilidad CLI, y cerrar el slice con tests unitarios y regresion focalizada.
4. Generar y verificar una matriz de cobertura de migracion que clasifique el 100% de triggers y efectos legacy como mantener, ajustar, reemplazar o retirar antes del cierre de la feature.

## Validation Plan

- Automated: `go test ./internal/cards ./internal/game ./internal/progression ./internal/encounter ./internal/match ./internal/content/playerprofiles ./internal/content/fishprofiles ./internal/endings ./internal/presentation ./internal/app ./internal/cli`, `go test ./...`, `golangci-lint run`
- Manual: ejecutar `go run ./cmd/fishing-duel` y validar rounds con efectos simultaneos, desempates a favor del pez, discard oculto o reordenado, reshuffle y fatiga; ejecutar `go run ./cmd/fishing-run` para comprobar que el flujo principal sigue operativo.
- Regression focus: orden actual implicito de efectos, thresholds legacy, condiciones terminales por captura o splash, compatibilidad de presets hardcoded y catalogos JSON, hints y snapshots mostrados por CLI.

## Risks / Tradeoffs

- El runtime actual separa bonuses de thresholds y shifts de posicion; unificar o coordinar ambas rutas puede tocar varios paquetes a la vez, pero evita dejar la prioridad repartida en reglas implicitas.
- Mantener compatibilidad con contenido actual puede requerir una capa temporal de migracion o normalizacion de efectos legacy, lo que aumenta complejidad a corto plazo a cambio de reducir quiebres de preset y catalogo.
- Exponer trazabilidad suficiente para tests y presenter puede agrandar snapshots o tipos de ronda, pero ayuda a verificar orden de resolucion sin acoplar la UI al engine.
