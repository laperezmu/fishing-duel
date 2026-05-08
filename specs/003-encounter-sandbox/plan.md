# Implementation Plan: Formalizar fishing-duel como sandbox de encounters

**Branch**: `003-encounter-sandbox` | **Date**: 2026-05-08 | **Spec**: `specs/003-encounter-sandbox/spec.md`
**Input**: Feature specification from `specs/003-encounter-sandbox/spec.md`

## Summary

Formalizar `fishing-duel` como un sandbox de encounters con setup guiado y manual, seleccion directa de preset de pez y cartas concretas, seed reproducible, trazas de runtime mas ricas y escenarios reutilizables, manteniendo la separacion de responsabilidades del proyecto, el flujo existente de `fishing-run` y el principio de cambios minimos sobre el runtime. El enfoque prioriza seams en `internal/app`, `internal/cli`, `internal/presentation` y `internal/content/*`, reutilizando engine y snapshots actuales en vez de abrir un rediseño transversal.

## Technical Context

**Language/Version**: Go 1.22.x  
**Primary Dependencies**: standard library, existing repo packages, `testify`, `golangci-lint`  
**Storage**: in-memory runtime and repository files for reusable scenarios unless the spec says otherwise  
**Testing**: targeted `go test` package runs, `go test ./...`, `golangci-lint run`, package-level coverage reports, manual CLI verification  
**Target Platform**: CLI on macOS/Linux/Windows  
**Project Type**: modular Go application with CLI entrypoints under `cmd/`  
**Constraints**: respetar division de responsabilidades y paquetes; seguir Effective Go y principios SOLID pragmaticos; mantener `internal/game`, `internal/encounter`, `internal/progression` y `internal/content/*` UI-agnostic; no afectar el flujo de `fishing-run`; agregar o actualizar tests unitarios cuando cambie comportamiento; alcanzar cobertura minima del 90% sobre los paquetes tocados por la feature; realizar cambios minimos necesarios para entregar el sandbox  
**Scale/Scope**: single feature slice over the encounter sandbox surface, with guided/manual setup, richer inspection, and reproducible scenarios

## Constitution Check

- [x] Spec is the active source of truth for this feature
- [x] Scope is explicit, with out-of-scope items recorded
- [x] Planned changes preserve modular package boundaries
- [x] Validation path includes tests and lint where behavior changes materially
- [x] Risks, assumptions, and tradeoffs are documented

Post-design check:

- [x] El plan conserva runtime, contenido, presentation, CLI y bootstrap en paquetes de ownership clara.
- [x] La propuesta evita contaminar `fishing-run` con conceptos de sandbox y mantiene su flujo actual diferenciado.
- [x] El trabajo se divide en slices finos sobre setup, trazas y escenarios, evitando un rewrite amplio del engine.
- [x] La validacion incluye actualizacion de tests unitarios y objetivo de cobertura minima del 90% en el area modificada.

## Project Structure

### Documentation (this feature)

```text
specs/003-encounter-sandbox/
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

**Structure Decision**: Introducir el sandbox formal como un flujo de composicion separado sobre piezas ya existentes. Extender `internal/app`, `internal/cli`, `internal/presentation` y `internal/content/*` solo donde la feature lo requiera, dejando `internal/run` y el contrato central de gameplay aislados salvo consumo legitimo de snapshots o configuracion.

## Implementation Approach

### Package Impact

- Primary packages to change: `cmd/fishing-duel`, `internal/app`, `internal/cli`, `internal/presentation`, `internal/content/playerprofiles`, `internal/content/fishprofiles`, `internal/content/watercontexts`, `internal/match`
- Secondary packages to review: `internal/game`, `internal/encounter`, `internal/progression`, `internal/deck`, `internal/run`, `cmd/fishing-run`
- Files or areas intentionally avoided: `internal/run` behavior, `cmd/fishing-run` user flow, global gameplay rules unrelated to sandbox configuration, broad refactors of cards/effects outside seams required for trace visibility or state seeding

### Execution Strategy

1. Extraer un flujo de sandbox formal en `internal/app` con configuracion propia, manteniendo un modo guiado compatible y agregando un modo manual para elegir presets, preset de pez, cartas concretas, seed y overrides de apertura sin tocar el loop de `fishing-run`.
2. Extender `internal/content/*`, `internal/presentation` y `internal/cli` para soportar seleccion manual de cartas, escenarios reutilizables, modos de salida y trazas de runtime mas ricas, preservando boundaries UI-agnostic y minimizando cambios sobre engine.
3. Incorporar solo las seams estrictamente necesarias en `internal/game`, `internal/match`, `internal/encounter` y `internal/progression` para aceptar configuracion semilla del sandbox y exponer evidencia estructurada, cerrando con tests unitarios actualizados y cobertura >=90% en paquetes impactados.

## Validation Plan

- Automated: `go test ./internal/app ./internal/cli ./internal/presentation ./internal/content/playerprofiles ./internal/content/fishprofiles ./internal/content/watercontexts ./internal/match ./internal/game ./internal/encounter ./internal/progression`, `go test ./...`, `golangci-lint run`, `go test -cover ./internal/app ./internal/cli ./internal/presentation ./internal/content/playerprofiles ./internal/content/fishprofiles ./internal/content/watercontexts ./internal/match ./internal/game ./internal/encounter ./internal/progression`
- Manual: ejecutar el sandbox en modo guiado y manual, validar seleccion de preset de pez, cartas concretas, seed fija, escenario reutilizable y traza de efectos; ejecutar `go run ./cmd/fishing-run` para confirmar que el flujo de run sigue intacto.
- Regression focus: bootstrap de encounters, spawn de pez derivado vs manual, trazas de resolucion, visibilidad de descarte, estados de splash, overrides de apertura, y ausencia de impacto funcional en `fishing-run`.

## Risks / Tradeoffs

- El override de estado avanzado puede forzar seams nuevas sobre deck o encounter; para mantener cambios minimos conviene limitar el primer alcance a overrides de apertura y estado previo estrictamente validado.
- Exponer demasiada configuracion en un solo flujo CLI puede degradar UX; separar guiado y manual reduce ese riesgo sin duplicar reglas.
- El objetivo de cobertura >=90% en paquetes tocados puede requerir redistribuir logica fuera de handlers de CLI para hacerla mas testeable, lo cual es positivo pero debe mantenerse acotado.
- Reusar `cmd/fishing-duel` como base reduce costo, pero el renombre y el nuevo framing exigen cuidar rutas, textos y referencias sin romper scripts o habitos existentes durante la transicion.
