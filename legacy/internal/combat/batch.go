package combat

import "fmt"

type BatchScenario struct {
	Name     string
	Category string
	Config   CombatConfig
}

type BatchOptions struct {
	Simulations int
	Seed        int64
}

type BatchScenarioReport struct {
	Scenario            BatchScenario
	InitialActionEV     [actionCount]float64
	PolicyResults       map[string]ScenarioResult
	ForcedActionResults [actionCount]ScenarioResult
	BestForcedAction    Action
	BestForcedRate      float64
}

type BatchReport struct {
	Options   BatchOptions
	Scenarios []BatchScenarioReport
}

func DefaultBatchScenarios() []BatchScenario {
	fishes := []FishProfile{
		NewBlackFish(),
		NewMonoFish(ColorRed),
		NewMonoFish(ColorBlue),
		NewMonoFish(ColorYellow),
		NewBiColorFish(ColorRed, ColorYellow),
		NewBiColorFish(ColorBlue, ColorRed),
		NewBiColorFish(ColorBlue, ColorYellow),
	}

	var scenarios []BatchScenario
	for _, slot := range []int{3, 4, 5} {
		for _, fish := range fishes {
			scenarios = append(scenarios, BatchScenario{
				Name:     fmt.Sprintf("baseline-slot-%d-%s", slot, fish.String()),
				Category: "baseline",
				Config: CombatConfig{
					InitialTrackPos: slot,
					Fish:            fish,
				},
			})
		}
	}

	scenarios = append(scenarios,
		BatchScenario{
			Name:     "bait-slot-5-blue-red",
			Category: "loadout",
			Config: CombatConfig{
				InitialTrackPos:    5,
				Fish:               NewBiColorFish(ColorBlue, ColorRed),
				BaitGuardAvailable: true,
			},
		},
		BatchScenario{
			Name:     "off-red-slot-4-red",
			Category: "loadout",
			Config: CombatConfig{
				InitialTrackPos: 4,
				Fish:            NewMonoFish(ColorRed),
				RodMods:         RodModifiers{}.WithOffensive(ColorRed),
			},
		},
	)

	return scenarios
}

func RunBatchReport(scenarios []BatchScenario, opts BatchOptions) (BatchReport, error) {
	if opts.Simulations <= 0 {
		return BatchReport{}, fmt.Errorf("simulations must be > 0")
	}
	if len(scenarios) == 0 {
		return BatchReport{}, fmt.Errorf("at least one scenario is required")
	}

	report := BatchReport{Options: opts, Scenarios: make([]BatchScenarioReport, 0, len(scenarios))}
	for i, scenario := range scenarios {
		entry, err := runBatchScenario(scenario, opts, i)
		if err != nil {
			return BatchReport{}, err
		}
		report.Scenarios = append(report.Scenarios, entry)
	}
	return report, nil
}

func runBatchScenario(scenario BatchScenario, opts BatchOptions, scenarioIndex int) (BatchScenarioReport, error) {
	initial := BeliefState{
		TrackPos:           uint8(scenario.Config.InitialTrackPos),
		Fish:               scenario.Config.Fish,
		Draw:               NewFamilyCounts(3, 3, 3),
		Discard:            FamilyCounts{},
		FatigueCount:       0,
		BaitGuardAvailable: scenario.Config.BaitGuardAvailable,
		RodMods:            scenario.Config.RodMods,
	}

	optimal := NewOptimalPolicy()
	actionEV := optimal.ActionValues(initial).Values
	report := BatchScenarioReport{
		Scenario:        scenario,
		InitialActionEV: actionEV,
		PolicyResults:   make(map[string]ScenarioResult),
	}

	policyNames := []string{"optimal", "heuristic", "random"}
	for idx, policyName := range policyNames {
		result, err := RunScenario(Scenario{
			Config:      scenario.Config,
			PolicyName:  policyName,
			Simulations: opts.Simulations,
			Seed:        deriveSeed(opts.Seed, scenarioIndex, idx),
		})
		if err != nil {
			return BatchScenarioReport{}, err
		}
		report.PolicyResults[policyName] = result
	}

	bestRate := -1.0
	bestAction := Forzar
	for idx, action := range allActions {
		result, err := runScenarioWithPolicy(Scenario{
			Config:      scenario.Config,
			Simulations: opts.Simulations,
			Seed:        deriveSeed(opts.Seed, scenarioIndex, 100+idx),
		}, NewFixedPolicy(action))
		if err != nil {
			return BatchScenarioReport{}, err
		}
		report.ForcedActionResults[action] = result
		if result.CaptureRate > bestRate {
			bestRate = result.CaptureRate
			bestAction = action
		}
	}

	report.BestForcedAction = bestAction
	report.BestForcedRate = bestRate
	return report, nil
}

func runScenarioWithPolicy(s Scenario, policy PlayerPolicy) (ScenarioResult, error) {
	if s.Simulations <= 0 {
		return ScenarioResult{}, fmt.Errorf("simulations must be > 0")
	}
	if s.Config.InitialTrackPos < 1 || s.Config.InitialTrackPos > 5 {
		return ScenarioResult{}, fmt.Errorf("initial track position must be between 1 and 5")
	}

	master := newDeterministicRand(s.Seed)
	return runScenarioWithRand(s, policy, master), nil
}

func deriveSeed(base int64, scenarioIndex, offset int) int64 {
	return base + int64((scenarioIndex+1)*1000+offset*17)
}
