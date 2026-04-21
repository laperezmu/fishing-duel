BATCH_SIMULATIONS ?= 3000
BATCH_SEED ?= 42
BATCH_MARKDOWN ?= reports/combat_batch_report.md
BATCH_JSON ?= reports/combat_batch_report.json

.PHONY: test batch-report batch-report-fast

test:
	go test ./...

batch-report:
	go run ./cmd/combat-report --simulations $(BATCH_SIMULATIONS) --seed $(BATCH_SEED) --markdown-out $(BATCH_MARKDOWN) --json-out $(BATCH_JSON)

batch-report-fast:
	go run ./cmd/combat-report --simulations 500 --seed $(BATCH_SEED) --markdown-out $(BATCH_MARKDOWN) --json-out $(BATCH_JSON)
