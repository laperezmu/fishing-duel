package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"pesca/internal/combat"
)

//go run ./cmd/combat-sim --slot 3 --fish blue-red --policy optimal --simulations 10000 --seed 42

func main() {
	var (
		slot        = flag.Int("slot", 4, "initial track slot (1-5)")
		fish        = flag.String("fish", "blue-red", "fish profile: black, red, blue, yellow, red-yellow, blue-red, blue-yellow")
		policy      = flag.String("policy", "optimal", "policy: optimal, heuristic, random")
		simulations = flag.Int("simulations", 10000, "number of combats to simulate")
		seed        = flag.Int64("seed", 42, "random seed")
		bait        = flag.Bool("bait", false, "whether bait guard is active")
		offensive   = flag.String("offensive", "", "comma-separated offensive rod colors")
		defensive   = flag.String("defensive", "", "comma-separated defensive rod colors")
		jsonOut     = flag.Bool("json", false, "print result as json")
	)
	flag.Parse()

	profile, err := combat.ParseFishProfile(*fish)
	if err != nil {
		fail(err)
	}

	rodMods, err := parseRodModifiers(*offensive, *defensive)
	if err != nil {
		fail(err)
	}

	scenario := combat.Scenario{
		Config: combat.CombatConfig{
			InitialTrackPos:    *slot,
			Fish:               profile,
			BaitGuardAvailable: *bait,
			RodMods:            rodMods,
		},
		PolicyName:  *policy,
		Simulations: *simulations,
		Seed:        *seed,
	}

	result, err := combat.RunScenario(scenario)
	if err != nil {
		fail(err)
	}

	if *jsonOut {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(result); err != nil {
			fail(err)
		}
		return
	}

	fmt.Printf("Combat simulation\n")
	fmt.Printf("scenario: slot=%d fish=%s policy=%s sims=%d seed=%d\n", result.Scenario.Config.InitialTrackPos, result.Scenario.Config.Fish, result.Policy, result.Scenario.Simulations, result.Scenario.Seed)
	fmt.Printf("capture_rate: %.4f\n", result.CaptureRate)
	fmt.Printf("escape_rate: %.4f\n", result.EscapeRate)
	fmt.Printf("avg_rounds: %.4f\n", result.AvgRounds)
	fmt.Printf("avg_fatigues: %.4f\n", result.AvgFatigues)
	fmt.Printf("bait_save_rate: %.4f\n", result.BaitSaveRate)
	fmt.Printf("offensive_trigger_rate: %.4f\n", result.OffensiveTriggerRate)
	fmt.Printf("defensive_trigger_rate: %.4f\n", result.DefensiveTriggerRate)
	fmt.Printf("action_usage: forzar=%d tensar=%d soltar=%d\n", result.ActionUsage[combat.Forzar], result.ActionUsage[combat.Tensar], result.ActionUsage[combat.Soltar])
	fmt.Printf("initial_action_ev: forzar=%.4f tensar=%.4f soltar=%.4f\n", result.InitialActionEV[combat.Forzar], result.InitialActionEV[combat.Tensar], result.InitialActionEV[combat.Soltar])

	keys := make([]string, 0, len(result.TerminalReasonBreakdown))
	for key := range result.TerminalReasonBreakdown {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	fmt.Println("terminal_reason_breakdown:")
	for _, key := range keys {
		fmt.Printf("  %s=%d\n", key, result.TerminalReasonBreakdown[key])
	}
}

func parseRodModifiers(offensiveRaw, defensiveRaw string) (combat.RodModifiers, error) {
	var mods combat.RodModifiers
	for _, part := range strings.Split(offensiveRaw, ",") {
		if strings.TrimSpace(part) == "" {
			continue
		}
		color, err := combat.ParseColor(part)
		if err != nil {
			return combat.RodModifiers{}, err
		}
		mods = mods.WithOffensive(color)
	}
	for _, part := range strings.Split(defensiveRaw, ",") {
		if strings.TrimSpace(part) == "" {
			continue
		}
		color, err := combat.ParseColor(part)
		if err != nil {
			return combat.RodModifiers{}, err
		}
		mods = mods.WithDefensive(color)
	}
	return mods, nil
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
