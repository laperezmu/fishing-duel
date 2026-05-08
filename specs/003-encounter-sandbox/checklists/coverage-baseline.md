# Coverage Baseline

Baseline coverage percentages before implementation starts.

**Date**: 2026-05-08

## Package Coverage

| Package | Current % | Target % | Notes |
|--------|----------|----------|-------|
| internal/app | 69.8% | >=90% | Baseline before sandbox implementation |
| internal/cli | 76.3% | >=90% | Baseline before sandbox implementation |
| internal/presentation | 85.2% | >=90% | Baseline before sandbox implementation |
| internal/content/playerprofiles | 84.0% | >=90% | Baseline before sandbox implementation |
| internal/content/fishprofiles | 81.6% | >=90% | Baseline before sandbox implementation |
| internal/content/watercontexts | 100.0% | >=90% | Baseline before sandbox implementation |
| internal/match | 100.0% | >=90% | Baseline before sandbox implementation |

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

## Final Verification

| Package | Final % | Target % | Status |
|--------|---------|----------|--------|
| internal/app | 69.8% | >=90% | below target |
| internal/cli | 77.1% | >=90% | below target |
| internal/presentation | 85.5% | >=90% | below target |
| internal/content/playerprofiles | 93.0% | >=90% | pass |
| internal/content/fishprofiles | 82.0% | >=90% | below target |
| internal/content/watercontexts | 100.0% | >=90% | pass |
| internal/match | 100.0% | >=90% | pass |
| internal/game | 98.8% | >=90% | pass |
| internal/encounter | 82.6% | >=90% | below target |
| internal/progression | 89.5% | >=90% | below target |
