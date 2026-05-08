# Tasks: Formalizar fishing-duel como sandbox de encounters

**Input**: Design documents from `/specs/003-encounter-sandbox/`
**Prerequisites**: `plan.md`, `spec.md`; use `research.md`, `data-model.md`, and `contracts/` when present

**Tests**: Include test tasks whenever behavior changes materially, not only when explicitly requested.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel
- **[Story]**: `US1`, `US2`, etc.
- Every task should name concrete file or package paths

## Phase 1: Setup

- [ ] T001 Review `specs/003-encounter-sandbox/spec.md`, `specs/003-encounter-sandbox/plan.md`, `specs/003-encounter-sandbox/research.md`, and current sandbox-facing packages under `cmd/` and `internal/`
- [ ] T002 Review current encounter bootstrap seams in `internal/app/bootstrap.go`, `internal/app/opening.go`, `internal/app/spawn.go`, `cmd/fishing-duel/main.go`, and `cmd/fishing-run/main.go`

---

## Phase 2: Foundations

- [ ] T003 Define shared sandbox configuration, scenario, override, and card-selection contracts in `internal/app` and `internal/match` with minimal UI-agnostic supporting types
- [ ] T004 [P] Add or update preset lookup helpers, fixture data, and scenario seed support in `internal/content/playerprofiles`, `internal/content/fishprofiles`, and `internal/content/watercontexts`
- [ ] T005 [P] Add or update baseline tests that guard sandbox contracts, run isolation, and coverage-sensitive seams in `internal/app/*_test.go`, `internal/content/*/*_test.go`, and `internal/match/*_test.go`

**Checkpoint**: shared foundations are ready and story work can proceed safely

---

## Phase 3: User Story 1 - Configurar encounters de forma explicita (Priority: P1)

**Goal**: Formalizar el sandbox con modo guiado y manual, permitiendo elegir presets, preset de pez, cartas concretas, seed y overrides de apertura sin afectar `fishing-run`.

**Independent Test**: Iniciar el sandbox, recorrer modo guiado y modo manual, seleccionar preset de jugador, preset de pez, cartas concretas, cana, aditamentos y contexto, y verificar que el encounter arranca con esa configuracion exacta.

- [ ] T006 [P] [US1] Add focused tests for sandbox setup resolution, seed reproducibility, and manual fish/card selection in `internal/app/bootstrap_test.go`, `internal/app/opening_test.go`, `internal/app/spawn_test.go`, and `internal/content/*/*_test.go`
- [ ] T007 [P] [US1] Implement sandbox configuration and setup resolution flow in `internal/app/bootstrap.go`, `internal/app/opening.go`, `internal/app/spawn.go`, and new supporting files under `internal/app`
- [ ] T008 [US1] Implement manual preset and card-selection support in `internal/content/playerprofiles`, `internal/content/fishprofiles`, and any supporting files under `internal/content/watercontexts`
- [ ] T009 [US1] Update CLI sandbox setup prompts and command framing in `cmd/fishing-duel/main.go`, `internal/cli/ui.go`, and related CLI setup files under `internal/cli`

**Checkpoint**: User Story 1 works independently

---

## Phase 4: User Story 2 - Inspeccionar la resolucion real del runtime (Priority: P2)

**Goal**: Exponer en el sandbox trazas y modos de salida que permitan inspeccionar triggers, efectos, prioridades, desempates y estado antes/despues de la resolucion sin mover reglas a la UI.

**Independent Test**: Ejecutar una ronda con efectos simultaneos y comprobar en la salida del sandbox que la traza de runtime muestra activaciones, orden, desempates y estado derivado de forma consistente con el engine.

- [ ] T010 [P] [US2] Add focused tests for detailed traces, snapshot enrichment, and debug-mode presentation in `internal/game/*_test.go`, `internal/match/*_test.go`, `internal/presentation/*_test.go`, and `internal/cli/*_test.go`
- [ ] T011 [P] [US2] Enrich runtime trace and sandbox-oriented snapshots in `internal/game`, `internal/match`, `internal/encounter`, and `internal/progression` only where strictly required by the feature
- [ ] T012 [US2] Integrate sandbox output modes and trace rendering in `internal/presentation/presenter.go`, `internal/presentation/views.go`, `internal/cli/render.go`, and related CLI flow files under `internal/cli`

**Checkpoint**: User Stories 1 and 2 both work independently

---

## Phase 5: User Story 3 - Reproducir escenarios y ejecutar pruebas de QA con baja friccion (Priority: P3)

**Goal**: Permitir ejecutar escenarios guardados y configuraciones no interactivas o semi-reproducibles para QA y regresion, manteniendo aislado el flujo de `fishing-run`.

**Independent Test**: Ejecutar un escenario QA predefinido o una configuracion reproducible, validar que el sandbox aplica el setup completo y confirmar que `fishing-run` sigue operativo sin prompts ni conceptos del sandbox.

- [ ] T013 [P] [US3] Add focused tests for reusable scenarios, non-interactive setup, and run-flow isolation in `internal/app/*_test.go`, `internal/cli/*_test.go`, and `internal/run/*_test.go` where applicable
- [ ] T014 [P] [US3] Implement reusable scenario loading, scenario-driven card selection, and seed-aware replay setup in `internal/app`, `internal/content/fishprofiles`, `internal/content/playerprofiles`, and new scenario-supporting files under `internal/content` or `internal/app`
- [ ] T015 [US3] Integrate non-interactive or semi-reproducible sandbox execution and verify `fishing-run` remains untouched in `cmd/fishing-duel/main.go`, `internal/cli`, `internal/app/run_bootstrap.go`, and `cmd/fishing-run/main.go`

**Checkpoint**: User Stories 1, 2, and 3 all work independently

---

## Final Phase: Validation and Finish

- [ ] T016 [P] Update sandbox documentation and scenario references in `specs/003-encounter-sandbox/quickstart.md`, `specs/003-encounter-sandbox/contracts/sandbox-cli-contract.md`, and any touched CLI README or command help text
- [ ] T017 Run agreed automated validation commands, including coverage verification, from `specs/003-encounter-sandbox/quickstart.md`
- [ ] T018 Run the manual verification path from `specs/003-encounter-sandbox/quickstart.md` for guided sandbox, manual sandbox, and `fishing-run`
- [ ] T019 Prepare PR notes linked to `specs/003-encounter-sandbox/spec.md`, `specs/003-encounter-sandbox/plan.md`, and completed sandbox scope in `specs/003-encounter-sandbox/tasks.md`

---

## Dependencies & Execution Order

- Setup first: `T001-T002`
- Foundations block all story work: `T003-T005`
- `US1` starts after foundations: `T006-T009`
- `US2` depends on `US1` because trace visibility and richer output need the final sandbox setup and stable manual configuration flow: `T010-T012`
- `US3` depends on `US1` and `US2` because scenarios and reproducible execution rely on the final setup contract and visible runtime evidence: `T013-T015`
- Validation happens after the desired stories are complete: `T016-T019`

## Parallel Execution Examples

- **US1**: Run `T006` and `T007` in parallel after foundations, then integrate with `T008-T009`
- **US2**: Run `T010` and `T011` in parallel before `T012`
- **US3**: Run `T013` and `T014` in parallel before `T015`
- **Foundations**: Run `T004` and `T005` in parallel after `T003`

## Implementation Strategy

- MVP first: complete Phase 3 (`US1`) to formalize the sandbox setup surface, manual fish/card selection, and seed support while keeping `fishing-run` isolated
- Next, complete Phase 4 (`US2`) to expose trustworthy runtime traces and richer output without moving rules into CLI or presentation
- Finally, complete Phase 5 (`US3`) to add reusable scenarios, reproducible execution paths, and final run-isolation validation

## Notes

- Favor small vertical slices over broad refactors
- Keep tasks mapped to a user story or to shared foundations
- If a task changes behavior, include the verification path in the same story phase
