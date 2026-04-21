package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"pesca/internal/combat"
)

type summaryStats struct {
	AverageCaptureByPolicy map[string]float64
	AverageForcedCapture   [3]float64
	BestForcedWins         [3]int
}

func main() {
	var (
		simulations = flag.Int("simulations", 3000, "number of simulations per scenario")
		seed        = flag.Int64("seed", 42, "base random seed")
		markdownOut = flag.String("markdown-out", "reports/combat_batch_report.md", "markdown output path")
		jsonOut     = flag.String("json-out", "reports/combat_batch_report.json", "json output path")
	)
	flag.Parse()

	report, err := combat.RunBatchReport(combat.DefaultBatchScenarios(), combat.BatchOptions{
		Simulations: *simulations,
		Seed:        *seed,
	})
	if err != nil {
		fail(err)
	}

	if err := writeJSON(*jsonOut, report); err != nil {
		fail(err)
	}
	if err := writeMarkdown(*markdownOut, report); err != nil {
		fail(err)
	}

	fmt.Printf("batch report written to %s and %s\n", *markdownOut, *jsonOut)
}

func writeJSON(path string, report combat.BatchReport) error {
	if err := ensureParentDir(path); err != nil {
		return err
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func writeMarkdown(path string, report combat.BatchReport) error {
	if err := ensureParentDir(path); err != nil {
		return err
	}

	var b strings.Builder
	b.WriteString("# Reporte batch de simulacion de combate\n\n")
	b.WriteString(fmt.Sprintf("- Escenarios: %d\n", len(report.Scenarios)))
	b.WriteString(fmt.Sprintf("- Simulaciones por escenario: %d\n", report.Options.Simulations))
	b.WriteString(fmt.Sprintf("- Seed base: %d\n\n", report.Options.Seed))

	summary := buildSummary(report)
	b.WriteString("## Resumen global\n\n")
	b.WriteString("| Metrica | Forzar | Tensar | Soltar |\n")
	b.WriteString("| --- | ---: | ---: | ---: |\n")
	b.WriteString(fmt.Sprintf("| Capture rate promedio forzado | %.4f | %.4f | %.4f |\n", summary.AverageForcedCapture[combat.Forzar], summary.AverageForcedCapture[combat.Tensar], summary.AverageForcedCapture[combat.Soltar]))
	b.WriteString(fmt.Sprintf("| Escenarios donde la accion forzada fue la mejor | %d | %d | %d |\n\n", summary.BestForcedWins[combat.Forzar], summary.BestForcedWins[combat.Tensar], summary.BestForcedWins[combat.Soltar]))

	b.WriteString("| Politica | Capture rate promedio |\n")
	b.WriteString("| --- | ---: |\n")
	policies := []string{"optimal", "heuristic", "random"}
	for _, policy := range policies {
		b.WriteString(fmt.Sprintf("| %s | %.4f |\n", policy, summary.AverageCaptureByPolicy[policy]))
	}
	b.WriteString("\n")

	grouped := groupByCategory(report.Scenarios)
	for _, category := range []string{"baseline", "loadout"} {
		scenarios := grouped[category]
		if len(scenarios) == 0 {
			continue
		}
		b.WriteString(fmt.Sprintf("## Escenarios %s\n\n", category))
		b.WriteString("| Escenario | Pez | Slot | Loadout | EV exacto F/T/S | Forced cap F/T/S | Mejor forzada | Optimal | Heuristic | Random |\n")
		b.WriteString("| --- | --- | ---: | --- | --- | --- | --- | ---: | ---: | ---: |\n")
		for _, scenario := range scenarios {
			b.WriteString(fmt.Sprintf(
				"| %s | %s | %d | %s | %.4f / %.4f / %.4f | %.4f / %.4f / %.4f | %s (%.4f) | %.4f | %.4f | %.4f |\n",
				scenario.Scenario.Name,
				scenario.Scenario.Config.Fish.String(),
				scenario.Scenario.Config.InitialTrackPos,
				loadoutLabel(scenario.Scenario.Config),
				scenario.InitialActionEV[combat.Forzar],
				scenario.InitialActionEV[combat.Tensar],
				scenario.InitialActionEV[combat.Soltar],
				scenario.ForcedActionResults[combat.Forzar].CaptureRate,
				scenario.ForcedActionResults[combat.Tensar].CaptureRate,
				scenario.ForcedActionResults[combat.Soltar].CaptureRate,
				scenario.BestForcedAction.String(),
				scenario.BestForcedRate,
				scenario.PolicyResults["optimal"].CaptureRate,
				scenario.PolicyResults["heuristic"].CaptureRate,
				scenario.PolicyResults["random"].CaptureRate,
			))
		}
		b.WriteString("\n")
	}

	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func buildSummary(report combat.BatchReport) summaryStats {
	summary := summaryStats{AverageCaptureByPolicy: map[string]float64{"optimal": 0, "heuristic": 0, "random": 0}}
	count := float64(len(report.Scenarios))
	for _, scenario := range report.Scenarios {
		for _, policy := range []string{"optimal", "heuristic", "random"} {
			summary.AverageCaptureByPolicy[policy] += scenario.PolicyResults[policy].CaptureRate / count
		}
		for _, action := range combat.AllActions() {
			summary.AverageForcedCapture[action] += scenario.ForcedActionResults[action].CaptureRate / count
		}
		summary.BestForcedWins[scenario.BestForcedAction]++
	}
	return summary
}

func groupByCategory(scenarios []combat.BatchScenarioReport) map[string][]combat.BatchScenarioReport {
	grouped := make(map[string][]combat.BatchScenarioReport)
	for _, scenario := range scenarios {
		grouped[scenario.Scenario.Category] = append(grouped[scenario.Scenario.Category], scenario)
	}
	for category := range grouped {
		sort.Slice(grouped[category], func(i, j int) bool {
			return grouped[category][i].Scenario.Name < grouped[category][j].Scenario.Name
		})
	}
	return grouped
}

func loadoutLabel(config combat.CombatConfig) string {
	parts := make([]string, 0, 3)
	if config.BaitGuardAvailable {
		parts = append(parts, "bait")
	}
	if config.RodMods.Offensive != 0 {
		parts = append(parts, "off="+maskLabel(config.RodMods.Offensive))
	}
	if len(parts) == 0 {
		return "none"
	}
	return strings.Join(parts, ",")
}

func maskLabel(mask combat.ColorMask) string {
	parts := make([]string, 0, 3)
	if mask&(1<<0) != 0 {
		parts = append(parts, "red")
	}
	if mask&(1<<1) != 0 {
		parts = append(parts, "blue")
	}
	if mask&(1<<2) != 0 {
		parts = append(parts, "yellow")
	}
	return strings.Join(parts, "+")
}

func ensureParentDir(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
