---
description: "Task list template for Pesca feature delivery"
---

# Tasks: [FEATURE NAME]

**Input**: Design documents from `/specs/[###-feature-name]/`
**Prerequisites**: `plan.md`, `spec.md`; use `research.md`, `data-model.md`, and `contracts/` when present

**Tests**: Include test tasks whenever behavior changes materially, not only when explicitly requested.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel
- **[Story]**: `US1`, `US2`, etc.
- Every task should name concrete file or package paths

## Phase 1: Setup

- [ ] T001 Review current spec, plan, and impacted packages
- [ ] T002 Create or update feature scaffolding only if the implementation requires it

---

## Phase 2: Foundations

- [ ] T003 Update shared types, policies, or contracts required by the feature
- [ ] T004 [P] Add or adjust supporting fixtures, catalogs, or configuration data
- [ ] T005 [P] Add or update baseline tests that guard the changed contract

**Checkpoint**: shared foundations are ready and story work can proceed safely

---

## Phase 3: User Story 1 - [Title] (Priority: P1)

**Goal**: [brief outcome]

**Independent Test**: [how to verify this slice]

- [ ] T006 [P] [US1] Add or update focused tests in [path]
- [ ] T007 [P] [US1] Implement primary behavior in [path]
- [ ] T008 [US1] Integrate the slice through the owning package flow in [path]
- [ ] T009 [US1] Update presentation or bootstrap wiring in [path] if required

**Checkpoint**: User Story 1 works independently

---

## Phase 4: User Story 2 - [Title] (Priority: P2)

**Goal**: [brief outcome]

**Independent Test**: [how to verify this slice]

- [ ] T010 [P] [US2] Add or update focused tests in [path]
- [ ] T011 [P] [US2] Implement secondary behavior in [path]
- [ ] T012 [US2] Integrate the slice with existing runtime or content flow in [path]

**Checkpoint**: User Stories 1 and 2 both work independently

---

## Final Phase: Validation and Finish

- [ ] T013 [P] Update relevant documentation referenced by the spec
- [ ] T014 Run agreed automated validation commands
- [ ] T015 Run the manual verification path from the plan
- [ ] T016 Prepare PR notes linked to the spec and completed scope

---

## Dependencies & Execution Order

- Setup first
- Foundations block story work
- Stories should follow priority order unless the plan explicitly allows parallel execution
- Validation happens after the desired stories are complete

## Notes

- Favor small vertical slices over broad refactors
- Keep tasks mapped to a user story or to shared foundations
- If a task changes behavior, include the verification path in the same story phase
