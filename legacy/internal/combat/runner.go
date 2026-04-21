package combat

import (
	"fmt"
	"math/rand"
	"strings"
)

type Scenario struct {
	Config      CombatConfig
	PolicyName  string
	Simulations int
	Seed        int64
}

type ScenarioResult struct {
	Scenario                Scenario
	Policy                  string
	CaptureRate             float64
	EscapeRate              float64
	AvgRounds               float64
	AvgFatigues             float64
	ActionUsage             [actionCount]int
	InitialActionEV         [actionCount]float64
	TerminalReasonBreakdown map[string]int
	BaitSaveRate            float64
	OffensiveTriggerRate    float64
}

func NewPolicyByName(name string, rng *rand.Rand) (PlayerPolicy, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "", "optimal":
		return NewOptimalPolicy(), nil
	case "heuristic":
		return NewHeuristicPolicy(), nil
	case "random":
		return NewRandomPolicy(rng), nil
	default:
		return nil, fmt.Errorf("unknown policy %q", name)
	}
}

func RunScenario(s Scenario) (ScenarioResult, error) {
	if s.Simulations <= 0 {
		return ScenarioResult{}, fmt.Errorf("simulations must be > 0")
	}
	if s.Config.InitialTrackPos < 1 || s.Config.InitialTrackPos > 5 {
		return ScenarioResult{}, fmt.Errorf("initial track position must be between 1 and 5")
	}

	master := newDeterministicRand(s.Seed)
	policy, err := NewPolicyByName(s.PolicyName, master.Rand())
	if err != nil {
		return ScenarioResult{}, err
	}

	return runScenarioWithRand(s, policy, master), nil
}

type deterministicRandom interface {
	Rand() *rand.Rand
}

type randomWrapper struct {
	rng *rand.Rand
}

func newDeterministicRand(seed int64) *randomWrapper {
	return &randomWrapper{rng: rand.New(rand.NewSource(seed))}
}

func (r *randomWrapper) Rand() *rand.Rand {
	return r.rng
}

func runScenarioWithRand(s Scenario, policy PlayerPolicy, master deterministicRandom) ScenarioResult {

	result := ScenarioResult{
		Scenario:                s,
		Policy:                  policy.Name(),
		TerminalReasonBreakdown: make(map[string]int),
	}

	optimal := NewOptimalPolicy()
	initial := BeliefState{
		TrackPos:           uint8(s.Config.InitialTrackPos),
		Fish:               s.Config.Fish,
		Draw:               NewFamilyCounts(3, 3, 3),
		Discard:            FamilyCounts{},
		FatigueCount:       0,
		BaitGuardAvailable: s.Config.BaitGuardAvailable,
		RodMods:            s.Config.RodMods,
	}
	result.InitialActionEV = optimal.ActionValues(initial).Values

	var captures, escapes int
	var roundsTotal, fatiguesTotal int
	var baitSaves, offensiveTriggers int

	for i := 0; i < s.Simulations; i++ {
		state := NewCombatState(s.Config, master.Rand())
		combatResult := RunCombat(state, policy, master.Rand())
		if combatResult.Outcome == OutcomeCapture {
			captures++
		}
		if combatResult.Outcome == OutcomeEscape {
			escapes++
		}
		roundsTotal += combatResult.Rounds
		fatiguesTotal += combatResult.Fatigues
		baitSaves += combatResult.BaitSaves
		offensiveTriggers += combatResult.OffensiveModTriggers
		for idx, uses := range combatResult.ActionUsage {
			result.ActionUsage[idx] += uses
		}
		result.TerminalReasonBreakdown[combatResult.TerminalReason.String()]++
	}

	result.CaptureRate = float64(captures) / float64(s.Simulations)
	result.EscapeRate = float64(escapes) / float64(s.Simulations)
	result.AvgRounds = float64(roundsTotal) / float64(s.Simulations)
	result.AvgFatigues = float64(fatiguesTotal) / float64(s.Simulations)
	result.BaitSaveRate = float64(baitSaves) / float64(s.Simulations)
	result.OffensiveTriggerRate = float64(offensiveTriggers) / float64(s.Simulations)

	return result
}
