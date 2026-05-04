# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]
**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

## Summary

[Concise description of the intended implementation and why it fits the spec]

## Technical Context

**Language/Version**: Go [NEEDS CONFIRMATION IF CHANGING]  
**Primary Dependencies**: standard library, existing repo packages, [additional dependency or `none`]  
**Storage**: in-memory runtime and repository files unless the spec says otherwise  
**Testing**: `go test ./...`, targeted package tests, `golangci-lint run`  
**Target Platform**: CLI on macOS/Linux/Windows unless the spec narrows it  
**Project Type**: modular Go application with CLI entrypoints under `cmd/`  
**Constraints**: preserve current modular boundaries; keep domain packages UI-agnostic  
**Scale/Scope**: single feature slice for the current game/runtime architecture

## Constitution Check

- [ ] Spec is the active source of truth for this feature
- [ ] Scope is explicit, with out-of-scope items recorded
- [ ] Planned changes preserve modular package boundaries
- [ ] Validation path includes tests and lint where behavior changes materially
- [ ] Risks, assumptions, and tradeoffs are documented

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

- Primary packages to change: [list package paths]
- Secondary packages to review: [list package paths or `none`]
- Files or areas intentionally avoided: [list package paths or `none`]

### Execution Strategy

1. [first implementation step]
2. [second implementation step]
3. [validation or rollout step]

## Validation Plan

- Automated: [specific `go test` packages, `go test ./...`, lint, or other checks]
- Manual: [CLI flow or scenario to exercise]
- Regression focus: [areas most likely to break]

## Risks / Tradeoffs

- [risk or tradeoff]
- [risk or tradeoff]
