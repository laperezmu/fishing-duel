# Tasks: Splash interactivo con saltos y mejoras de cana

**Input**: Design documents from `/specs/001-splash-interactivo/`
**Prerequisites**: `plan.md`, `spec.md`; use `research.md`, `data-model.md`, and `contracts/` when present

**Tests**: Include test tasks whenever behavior changes materially, not only when explicitly requested.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel
- **[Story]**: `US1`, `US2`, etc.
- Every task should name concrete file or package paths

## Phase 1: Setup

- [X] T001 Review `specs/001-splash-interactivo/spec.md`, `specs/001-splash-interactivo/plan.md`, and impacted packages listed in the plan
- [X] T002 Review current splash flow in `internal/encounter/transition.go`, `internal/progression/track.go`, `internal/app/session.go`, and `internal/cli/opening.go`

---

## Phase 2: Foundations

- [X] T003 Update shared splash domain types and validation in `internal/encounter/types.go` to represent configurable jump-based splash state instead of random-only escape chance
- [X] T004 [P] Add splash sequence state and snapshot support in `internal/match` package files and their tests under `internal/match/*_test.go`
- [X] T005 [P] Refactor splash transition and progression contracts in `internal/encounter/transition.go`, `internal/progression/track.go`, and `internal/progression/splash.go`
- [X] T006 [P] Add baseline regression tests for the new splash contract in `internal/encounter/transition_test.go`, `internal/encounter/types_test.go`, and `internal/progression/track_test.go`
- [X] T007 Extend round and status presentation contracts for pending splash state in `internal/presentation/views.go` and `internal/presentation/presenter.go`
- [X] T008 [P] Add regression tests for preserved splash triggers and round/end snapshots in `internal/progression/track_test.go`, `internal/match/snapshots_test.go`, `internal/endings/encounter_test.go`, and `internal/presentation/presenter_test.go`

**Checkpoint**: shared foundations are ready and story work can proceed safely

---

## Phase 3: User Story 1 - Resolver un splash con habilidad (Priority: P1)

**Goal**: Reemplazar la resolucion aleatoria del splash por una interaccion jugable de timing con exito/fallo y timeout, manteniendo el encounter operativo desde CLI.

**Independent Test**: Activar un splash en un encounter aislado y verificar que la UI muestra una prueba con tiempo limite cuyo resultado decide si el pez sigue sujeto o escapa.

- [X] T009 [P] [US1] Add app/session splash flow tests in `internal/app/session_test.go` for pending splash resolution, success continuation, splash escape failure, and success-vs-timeout edge handling
- [X] T010 [P] [US1] Add CLI splash interaction tests in `internal/cli/ui_test.go` and `internal/cli/render_test.go` for timeout, success, failure rendering, and partial-progress failure feedback
- [X] T011 [P] [US1] Add presentation tests for splash views, labels, and pending-state snapshots in `internal/presentation/presenter_test.go`
- [X] T012 [US1] Introduce splash resolution contracts in `internal/app/session.go` and related app interfaces under `internal/app` for detecting and resolving pending splash states
- [X] T013 [US1] Implement splash-specific view models in `internal/presentation/views.go` and presenter mapping in `internal/presentation/presenter.go`
- [X] T014 [US1] Implement CLI splash rendering in `internal/cli/render.go` using the agreed cast-style interaction language from `specs/001-splash-interactivo/research.md`
- [X] T015 [US1] Implement CLI splash timing/input flow in `internal/cli/ui.go` with success, failure, and timeout resolution
- [X] T016 [US1] Integrate splash pending/result application in `internal/game` and `internal/encounter` package files updated for the new contract
- [X] T017 [US1] Add reproducible splash sandbox support for manual validation in `cmd/fishing-duel`, `internal/app`, or related fixtures so splash scenarios can be forced consistently during testing

**Checkpoint**: User Story 1 works independently

---

## Phase 4: User Story 2 - Resolver peces con multiples saltos (Priority: P2)

**Goal**: Permitir peces con secuencias configurables de 1 a 5 saltos y mostrar progreso claro hasta completar o fallar la secuencia.

**Independent Test**: Configurar peces de 1, 3 y 5 saltos y verificar que el jugador deba completar la cantidad exacta configurada, con escape inmediato al fallar cualquiera.

- [X] T018 [P] [US2] Add focused multi-jump tests in `internal/encounter/transition_test.go`, `internal/app/session_test.go`, and `internal/cli/ui_test.go` for 1, 3, and 5 jump sequences
- [X] T019 [P] [US2] Add content/configuration tests for splash profiles and per-profile time limits in `internal/content/fishprofiles` package tests
- [X] T020 [P] [US2] Add edge-case tests for splash progress, partial success followed by failure, and capture-or-end precedence in `internal/encounter/*_test.go`, `internal/app/session_test.go`, and `internal/cli/ui_test.go`
- [X] T021 [US2] Implement multi-jump splash sequencing and progress tracking in `internal/encounter` and `internal/match` package files
- [X] T022 [US2] Model and validate splash profile configuration fields, including jump count and per-profile time limit, in `internal/encounter/types.go` and `internal/content/fishprofiles`
- [X] T023 [US2] Wire splash profile source-of-truth through `internal/content/fishprofiles`, `internal/app/bootstrap.go`, and related spawn/bootstrap paths
- [X] T024 [US2] Update splash progress presentation in `internal/presentation/presenter.go`, `internal/presentation/views.go`, and `internal/cli/render.go`
- [X] T025 [US2] Add comparative content coverage for at least two fish splash profiles in `internal/content/fishprofiles` package tests to validate differentiated difficulty by configuration

**Checkpoint**: User Stories 1 and 2 both work independently

---

## Phase 5: User Story 3 - Aprovechar mejoras de cana durante el splash (Priority: P3)

**Goal**: Aplicar recompensas de acercamiento por salto ganado desde la build del jugador, sin romper la separacion entre loadout, encounter y CLI.

**Independent Test**: Resolver un splash con y sin bonus de cana y comprobar que cada salto ganado aplica el acercamiento adicional configurado sin cruzar el limite minimo valido.

- [X] T026 [P] [US3] Add loadout and content tests for splash bonuses in `internal/player/loadout/*_test.go`, `internal/content/rodpresets/*_test.go`, and `internal/content/anglerprofiles/*_test.go`
- [X] T027 [P] [US3] Add focused runtime tests for reward application, per-jump immediacy, clamping, and capture-boundary interaction in `internal/encounter/*_test.go` and `internal/app/session_test.go`
- [X] T028 [US3] Implement rod/loadout splash bonus capability in `internal/player/loadout/state.go`, `internal/content/rodpresets/presets.go`, and `internal/content/anglerprofiles/profiles.go`
- [X] T029 [US3] Apply per-jump reward handling in `internal/app/session.go`, `internal/encounter`, and `internal/game` package files responsible for splash resolution
- [X] T030 [US3] Expose splash reward messaging in `internal/presentation/presenter.go`, `internal/presentation/views.go`, and `internal/cli/render.go`

**Checkpoint**: User Stories 1, 2, and 3 all work independently

---

## Final Phase: Validation and Finish

- [X] T031 [P] Update feature documentation references in `specs/001-splash-interactivo/quickstart.md` and any touched package README files if the final behavior differs from the current description
- [X] T032 Run automated validation commands from `specs/001-splash-interactivo/quickstart.md`
- [X] T033 Run manual verification with `go run ./cmd/fishing-duel` and `go run ./cmd/fishing-run` following `specs/001-splash-interactivo/quickstart.md`
- [X] T034 Prepare PR notes linked to `specs/001-splash-interactivo/spec.md` and the completed splash scope in `specs/001-splash-interactivo/tasks.md`

---

## Dependencies & Execution Order

- Setup first: `T001-T002`
- Foundations block story work: `T003-T008`
- `US1` starts after foundations and defines the MVP slice: `T009-T017`
- `US2` depends on `US1` because multi-jump support extends the interactive splash flow: `T018-T025`
- `US3` depends on `US1` and `US2` because reward application builds on completed jump sequencing: `T026-T030`
- Validation happens after the desired stories are complete: `T031-T034`

## Parallel Execution Examples

- **US1**: Run `T009`, `T010`, and `T011` in parallel before `T012-T017`
- **US2**: Run `T018`, `T019`, and `T020` in parallel before `T021-T025`
- **US3**: Run `T026` and `T027` in parallel before `T028-T030`
- **Foundations**: Run `T004`, `T005`, `T006`, and `T008` in parallel after `T003`

## Implementation Strategy

- MVP first: complete Phase 3 (`US1`) to replace splash RNG with a playable interactive flow
- Incremental delivery: add configurable multi-jump support in Phase 4 and rod/loadout rewards in Phase 5
- Keep each story independently verifiable before moving to the next one
- Favor small vertical slices over broad refactors and keep runtime, app, presentation, and CLI concerns separated
- Use `internal/content/fishprofiles` as the source of truth for splash profile configuration and flow it into runtime through bootstrap wiring

## Notes

- Favor small vertical slices over broad refactors
- Keep tasks mapped to a user story or to shared foundations
- If a task changes behavior, include the verification path in the same story phase
- All tasks follow the required checklist format with checkbox, task ID, optional `[P]`, required story label for story phases, and concrete file paths
