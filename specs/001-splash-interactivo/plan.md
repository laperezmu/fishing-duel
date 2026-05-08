# Implementation Plan: Splash interactivo con saltos y mejoras de cana

**Branch**: `001-splash-interactivo-cana` | **Date**: 2026-05-07 | **Spec**: `specs/001-splash-interactivo/spec.md`
**Input**: Feature specification from `specs/001-splash-interactivo/spec.md`

## Summary

Implementar un flujo de splash interactivo que sustituya la resolucion aleatoria actual por una secuencia de saltos resuelta con input del jugador, manteniendo el motor modular y la compatibilidad con CLI. El diseno separa claramente: configuracion y estado tactico del splash en runtime UI-agnostic, orquestacion del flujo interactivo en `internal/app`, presentacion de vistas de splash en `internal/presentation`, y resolucion concreta de input/timing en `internal/cli`, con tests unitarios en encounter, app, presentation y CLI donde cambie el comportamiento.

## Technical Context

**Language/Version**: Go (version definida por el repo actual)  
**Primary Dependencies**: standard library, existing repo packages, `testify` ya presente en el repo, sin nuevas dependencias externas  
**Storage**: in-memory runtime and repository files unless the spec says otherwise  
**Testing**: `go test ./...`, targeted package tests, `golangci-lint run`  
**Target Platform**: CLI on macOS/Linux/Windows unless the spec narrows it  
**Project Type**: modular Go application with CLI entrypoints under `cmd/`  
**Constraints**: preserve current modular boundaries; keep domain packages UI-agnostic; mantener compatibilidad con `cmd/fishing-duel` y `cmd/fishing-run`; no reintroducir aleatoriedad directa en la resolucion final del splash interactivo  
**Scale/Scope**: single feature slice for the current game/runtime architecture

## Constitution Check

- [x] Spec is the active source of truth for this feature
- [x] Scope is explicit, with out-of-scope items recorded
- [x] Planned changes preserve modular package boundaries
- [x] Validation path includes tests and lint where behavior changes materially
- [x] Risks, assumptions, and tradeoffs are documented

Post-design check:

- [x] El runtime de encounter sigue UI-agnostic y no conoce terminal ni timing concreto.
- [x] La orquestacion interactiva queda en `internal/app` con implementaciones de borde en `internal/cli`.
- [x] El plan mantiene un slice pequeno: reemplaza splash aleatorio, agrega configuracion de saltos y bonus de cana, sin abrir economia ni nuevos sistemas meta.

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature-name]/
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

- Primary packages to change: `internal/encounter`, `internal/match`, `internal/game`, `internal/app`, `internal/presentation`, `internal/cli`, `internal/player/loadout`, `internal/content/rodpresets`, `internal/content/anglerprofiles`
- Secondary packages to review: `internal/endings`, `internal/progression`, `cmd/fishing-duel`, `cmd/fishing-run`
- Files or areas intentionally avoided: `internal/rules`, `internal/deck`, `internal/cards`, economia/meta, backlog docs legacy

### Execution Strategy

1. Refactorizar el runtime de splash para modelar secuencia, configuracion y resultado sin depender de UI ni RNG directo, y adaptar progresion/endings al nuevo contrato.
2. Introducir la orquestacion interactiva en app + presenter + CLI, incluyendo vista de saltos, input con tiempo limite y aplicacion de bonus de cana por salto ganado.
3. Actualizar presets/configuracion inicial, cubrir la regresion con tests unitarios y validar el flujo completo en ambos ejecutables CLI.

## Validation Plan

- Automated: `go test ./internal/encounter ./internal/progression ./internal/game ./internal/app ./internal/presentation ./internal/cli ./internal/player/loadout ./internal/content/rodpresets ./internal/content/anglerprofiles`, `go test ./...`, `golangci-lint run`
- Manual: ejecutar `go run ./cmd/fishing-duel` y forzar un splash; verificar casos de 1 salto, varios saltos, timeout/fallo y bonus de cana; ejecutar `go run ./cmd/fishing-run` para comprobar compatibilidad del flujo principal
- Regression focus: triggers actuales de splash, snapshots de ronda, causas terminales del encounter, lectura CLI del ultimo round, configuracion de loadout y rods existentes

## Risks / Tradeoffs

- Repartir la responsabilidad entre engine y session puede introducir mas pasos de orquestacion; a cambio se mantiene la separacion entre dominio y UI.
- Mantener compatibilidad con CLI puede empujar a un minijuego inicial mas simple de lo deseado; se acepta para preservar un primer slice pequeno y testeable.
- Incorporar bonus de cana en este slice amplía el impacto sobre loadout y contenido, pero evita dejar el splash como sistema aislado sin integracion con build.
