# Coverage Baseline

Baseline coverage percentages before implementation starts.

**Date**: 2026-05-08

## Package Coverage

| Package | Current % | Target % | Notes |
|--------|----------|----------|-------|
| internal/app | TBD | >=90% | |
| internal/cli | TBD | >=90% | |
| internal/presentation | TBD | >=90% | |
| internal/content/playerprofiles | TBD | >=90% | |
| internal/content/fishprofiles | TBD | >=90% | |
| internal/content/watercontexts | TBD | >=90% | |
| internal/match | TBD | >=90% | |

## Commands

Run before starting implementation:

```bash
go test -cover ./internal/app/...
go test -cover ./internal/cli/...
go test -cover ./internal/presentation/...
go test -cover ./internal/content/playerprofiles/...
go test -cover ./internal/content/fishprofiles/...
go test -cover ./internal/content/watercontexts/...
go test -cover ./internal/match/...
```