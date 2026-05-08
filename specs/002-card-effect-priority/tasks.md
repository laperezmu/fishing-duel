# Tasks: Rediseno de triggers y efectos de cartas con prioridad de resolucion

**Input**: Design documents from `/specs/002-card-effect-priority/`
**Prerequisites**: `plan.md`, `spec.md`; use `research.md`, `data-model.md`, and `contracts/` when present

**Tests**: Include test tasks whenever behavior changes materially, not only when explicitly requested.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel
- **[Story]**: `US1`, `US2`, etc.
- Every task should name concrete file or package paths

## Phase 1: Setup

- [X] T001 Review `specs/002-card-effect-priority/spec.md`, `specs/002-card-effect-priority/plan.md`, `specs/002-card-effect-priority/research.md`, and impacted packages under `internal/`
- [X] T002 Review current card effect flow in `internal/cards/cards.go`, `internal/game/engine.go`, `internal/progression/track.go`, `internal/encounter/thresholds.go`, and `internal/match/state.go`

---

## Phase 2: Foundations

- [X] T003 Define the new shared trigger and effect catalog contracts in `internal/cards/cards.go` and related files under `internal/cards`
- [X] T004 [P] Extend fish content effect parsing and validation for the new contract in `internal/content/fishprofiles/catalog_json.go`, `internal/content/fishprofiles/catalog.go`, and related tests under `internal/content/fishprofiles/*_test.go`
- [X] T005 [P] Extend player preset effect declarations for the new contract in `internal/content/playerprofiles/presets.go` and related tests under `internal/content/playerprofiles/*_test.go`
- [X] T006 [P] Add baseline contract regression tests for trigger filtering, priority ordering, and deprecated effect handling in `internal/cards/*_test.go`, `internal/game/engine_test.go`, and `internal/encounter/*_test.go`

**Checkpoint**: shared foundations are ready and story work can proceed safely

---

## Phase 3: User Story 1 - Definir un catalogo jugable de efectos y triggers (Priority: P1)

**Goal**: Cerrar el catalogo funcional de triggers y efectos soportados para cartas de pez y jugador, con bindings de un trigger por efecto y deprecacion explicita de efectos legacy.

**Independent Test**: Revisar una carta de pez y una del jugador con multiples efectos y verificar mediante tests de contenido y de dominio que cada efecto pertenece a un trigger permitido, usa el catalogo nuevo y no depende de bonuses legacy ambiguos.

- [X] T007 [P] [US1] Add focused catalog and validation tests in `internal/cards/*_test.go`, `internal/content/playerprofiles/presets_test.go`, and `internal/content/fishprofiles/catalog_test.go`
- [X] T008 [P] [US1] Implement trigger, effect, and entity-variant modeling in `internal/cards/cards.go` and new supporting files under `internal/cards`
- [X] T009 [US1] Integrate the new card effect contract into fish catalog loading in `internal/content/fishprofiles/catalog_json.go` and `internal/content/fishprofiles/profiles.go`
- [X] T010 [US1] Integrate the new card effect contract into player presets in `internal/content/playerprofiles/presets.go`
- [X] T011 [US1] Update effect hint rendering for the new catalog vocabulary in `internal/presentation/presenter.go` and `internal/presentation/presenter_test.go`

**Checkpoint**: User Story 1 works independently

---

## Phase 4: User Story 2 - Resolver efectos simultaneos de forma determinista (Priority: P2)

**Goal**: Reemplazar el orden implicito actual por una secuencia determinista por prioridad, con desempate a favor del pez y sin mover reglas a la UI.

**Independent Test**: Ejecutar tests de engine y progression donde jugador y pez disparan efectos simultaneos y comprobar que el orden resuelto es estable, priorizado y con desempate a favor del pez en todos los casos equivalentes.

- [X] T012 [P] [US2] Add priority-resolution tests in `internal/game/engine_test.go`, `internal/progression/track_test.go`, `internal/encounter/transition_test.go`, and `internal/endings/encounter_test.go`
- [X] T013 [P] [US2] Extend round and snapshot types for ordered effect resolution in `internal/match/state.go`, `internal/match/snapshots.go`, and `internal/match/*_test.go`
- [X] T014 [US2] Implement deterministic effect selection and ordering in `internal/game/engine.go` and supporting files under `internal/game`
- [X] T015 [US2] Refactor effect application over thresholds, movement, splash, and terminal checks in `internal/progression/track.go`, `internal/encounter/thresholds.go`, `internal/encounter/transition.go`, and related files under `internal/encounter`
- [X] T016 [US2] Update app and CLI-compatible presentation flow to consume resolved effect traces in `internal/app/session.go`, `internal/presentation/presenter.go`, `internal/cli/render.go`, and related tests under `internal/app`, `internal/presentation`, and `internal/cli`

**Checkpoint**: User Stories 1 and 2 both work independently

---

## Phase 5: User Story 3 - Migrar contenido y reglas sin perder cobertura funcional (Priority: P3)

**Goal**: Migrar el contenido existente y cerrar la cobertura funcional del contrato nuevo sin romper el runtime ni los flujos CLI.

**Independent Test**: Tomar el contenido actual del repo y verificar con tests de catalogo, presets y runtime que cada trigger o efecto existente queda migrado, reemplazado o retirado con comportamiento consistente y sin quiebres en los ejecutables CLI.

- [X] T017 [P] [US3] Add migration coverage tests for legacy-to-new mappings in `internal/content/fishprofiles/*_test.go`, `internal/content/playerprofiles/*_test.go`, and `internal/cards/*_test.go`
- [X] T018 [P] [US3] Add compatibility tests for round summaries and CLI-facing output in `internal/presentation/presenter_test.go`, `internal/app/session_test.go`, and `internal/cli/*_test.go`
- [X] T019 [US3] Produce and verify the migration coverage matrix in `specs/002-card-effect-priority/contracts/card-effect-resolution.md` and supporting assertions in `internal/content/fishprofiles/*_test.go` and `internal/content/playerprofiles/*_test.go`
- [X] T020 [US3] Migrate legacy effect usage and content fixtures in `internal/content/playerprofiles/presets.go`, `internal/content/fishprofiles/data/default_profiles.json`, and related files under `internal/content/fishprofiles`
- [X] T021 [US3] Update migration-sensitive runtime logic and ending checks in `internal/game/engine.go`, `internal/progression/track.go`, and `internal/endings/encounter.go`
- [X] T022 [US3] Finalize snapshot and presenter compatibility for migrated content in `internal/match/snapshots.go`, `internal/presentation/presenter.go`, and `internal/cli/render.go`

**Checkpoint**: User Stories 1, 2, and 3 all work independently

---

## Final Phase: Validation and Finish

- [X] T023 [P] Update feature documentation in `specs/002-card-effect-priority/quickstart.md`, `specs/002-card-effect-priority/contracts/card-effect-resolution.md`, and any touched README files if behavior differs from the current description
- [X] T024 Run agreed automated validation commands from `specs/002-card-effect-priority/quickstart.md` and record pass/fail outcomes in `specs/002-card-effect-priority/tasks.md`
- [X] T025 Run the manual verification path from `specs/002-card-effect-priority/quickstart.md` with `go run ./cmd/fishing-duel` and `go run ./cmd/fishing-run`, and record pass/fail notes in `specs/002-card-effect-priority/tasks.md`
- [X] T026 Prepare PR notes linked to `specs/002-card-effect-priority/spec.md`, `specs/002-card-effect-priority/plan.md`, and completed migration scope in `specs/002-card-effect-priority/tasks.md`

Validation notes:

- `T024` PASS - `go test ./internal/cards ./internal/game ./internal/progression ./internal/encounter ./internal/match ./internal/content/playerprofiles ./internal/content/fishprofiles ./internal/endings ./internal/presentation ./internal/app ./internal/cli`
- `T024` PASS - `go test ./...`
- `T024` PASS - `golangci-lint run`
- `T025` PASS - `go run ./cmd/fishing-duel` completo con input guiado; mostro historial del pez, ultimo lance y traza ordenada de efectos resueltos.
- `T025` PASS WITH NOTE - `go run ./cmd/fishing-run` avanzo hasta gameplay real y mostro resumenes/trazas de efectos; el proceso termino por EOF al usar input scriptado no interactivo, sin fallos funcionales previos al corte.

PR notes:

- Scope: implementa el nuevo contrato de triggers/efectos, prioridad determinista por efecto, desempate a favor del pez y migracion del contenido legacy para jugador y pez.
- Spec refs: `specs/002-card-effect-priority/spec.md`, `specs/002-card-effect-priority/plan.md`, `specs/002-card-effect-priority/contracts/card-effect-resolution.md`
- Runtime refs: `internal/cards/cards.go`, `internal/game/engine.go`, `internal/progression/track.go`, `internal/encounter/transition.go`, `internal/endings/encounter.go`
- Content refs: `internal/content/playerprofiles/presets.go`, `internal/content/fishprofiles/catalog_json.go`, `internal/content/fishprofiles/data/default_profiles.json`
- UX refs: `internal/match/snapshots.go`, `internal/presentation/presenter.go`, `internal/cli/render.go`
- Validation: `go test ./...` y `golangci-lint run` en verde; CLI duel validado con flujo guiado y run validado hasta gameplay con entrada scriptada.

---

## Dependencies & Execution Order

- Setup first: `T001-T002`
- Foundations block all story work: `T003-T006`
- `US1` starts after foundations: `T007-T011`
- `US2` depends on `US1` because deterministic resolution needs the final trigger/effect catalog and migrated shared card contracts: `T012-T016`
- `US3` depends on `US1` and `US2` because migration must target the final contract and runtime ordering behavior: `T017-T021`
- Validation happens after the desired stories are complete: `T023-T026`

## Parallel Execution Examples

- **US1**: Run `T007` and `T008` in parallel after foundations complete, then integrate with `T009-T010`
- **US2**: Run `T012` and `T013` in parallel before `T014`, then follow with `T015` and `T016`
- **US3**: Run `T017` and `T018` in parallel before `T019`, then finish with `T020-T022`

## Implementation Strategy

- MVP first: complete Phase 3 (`US1`) to lock the catalog contract and deprecations before touching runtime ordering
- Next, complete Phase 4 (`US2`) to deliver deterministic resolution with clear tests around priority and fish-first ties
- Finally, complete Phase 5 (`US3`) to migrate all active content and tighten compatibility across snapshots and CLI

## Notes

- Favor small vertical slices over broad refactors
- Keep tasks mapped to a user story or to shared foundations
- If a task changes behavior, include the verification path in the same story phase
