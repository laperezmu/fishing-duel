package combat

import (
	"math/rand"
	"testing"
)

func TestBaseResultMatrix(t *testing.T) {
	tests := []struct {
		fish   FishFamily
		action Action
		want   int
	}{
		{Embiste, Forzar, 0},
		{Embiste, Tensar, -1},
		{Embiste, Soltar, 1},
		{Aguante, Forzar, 1},
		{Aguante, Tensar, 0},
		{Aguante, Soltar, -1},
		{Quiebre, Forzar, -1},
		{Quiebre, Tensar, 1},
		{Quiebre, Soltar, 0},
	}

	for _, tt := range tests {
		if got := BaseResult(tt.fish, tt.action); got != tt.want {
			t.Fatalf("BaseResult(%s,%s)=%d want %d", tt.fish, tt.action, got, tt.want)
		}
	}
}

func TestApplyArchetype(t *testing.T) {
	mono := NewMonoFish(ColorRed)
	if got := ApplyArchetype(0, mono, Embiste, Forzar); got != -1 {
		t.Fatalf("mono archetype got %d want -1", got)
	}

	bicolor := NewBiColorFish(ColorBlue, ColorRed)
	if got := ApplyArchetype(0, bicolor, Aguante, Tensar); got != -1 {
		t.Fatalf("bicolor archetype got %d want -1", got)
	}

	if got := ApplyArchetype(0, bicolor, Quiebre, Soltar); got != 0 {
		t.Fatalf("unexpected archetype penalty got %d want 0", got)
	}
}

func TestApplyFatigueModifier(t *testing.T) {
	if got := ApplyFatigueModifier(1, 0); got != 0 {
		t.Fatalf("fatigue inactive got %d want 0", got)
	}
	if got := ApplyFatigueModifier(1, 1); got != 1 {
		t.Fatalf("fatigue active got %d want 1", got)
	}
	if got := ApplyFatigueModifier(-1, 3); got != -1 {
		t.Fatalf("fatigue cap got %d want -1", got)
	}
}

func TestBaitGuardAvoidsFirstEscapeOnly(t *testing.T) {
	round := EvaluateRound(5, NewBlackFish(), 0, true, RodModifiers{}, Aguante, Soltar)
	if round.Outcome != OutcomeOngoing || !round.BaitGuardConsumed || round.TrackPos != 5 {
		t.Fatalf("bait guard did not save correctly: %+v", round)
	}

	round = EvaluateRound(5, NewBlackFish(), 0, false, RodModifiers{}, Aguante, Soltar)
	if round.Outcome != OutcomeEscape {
		t.Fatalf("expected escape after bait consumed, got %+v", round)
	}
}

func TestOffensiveModifierTurnsTieIntoProgress(t *testing.T) {
	mods := RodModifiers{}.WithOffensive(ColorRed)
	round := EvaluateRound(3, NewBlackFish(), 0, false, mods, Embiste, Forzar)
	if round.Result != 1 || !round.OffensiveApplied || round.TrackPos != 2 {
		t.Fatalf("unexpected offensive modifier result: %+v", round)
	}
}

func TestOffensiveModifierUsesBaseResultNotArchetype(t *testing.T) {
	mods := RodModifiers{}.WithOffensive(ColorRed)
	round := EvaluateRound(3, NewMonoFish(ColorRed), 0, false, mods, Embiste, Forzar)
	if round.Result != 0 || !round.OffensiveApplied || round.TrackPos != 3 {
		t.Fatalf("unexpected offensive/archetype interaction: %+v", round)
	}
}

func TestTotalFatigueCapturesFish(t *testing.T) {
	state := NewCombatStateWithDeck(CombatConfig{
		InitialTrackPos: 3,
		Fish:            NewBlackFish(),
	}, NewFishDeckFromDraw([]FishFamily{Embiste, Embiste, Embiste}))

	result := RunCombat(state, fixedPolicy{action: Forzar}, rand.New(rand.NewSource(7)))
	if result.Outcome != OutcomeCapture {
		t.Fatalf("expected capture, got %+v", result)
	}
	if result.TerminalReason != TerminalReasonTotalFatigue {
		t.Fatalf("expected total fatigue, got %+v", result)
	}
	if result.Fatigues != 1 {
		t.Fatalf("expected one fatigue, got %+v", result)
	}
}

func TestOptimalPolicyFindsWinningAction(t *testing.T) {
	policy := NewOptimalPolicy()
	state := BeliefState{
		TrackPos:           5,
		Fish:               NewMonoFish(ColorRed),
		Draw:               NewFamilyCounts(0, 3, 0),
		Discard:            FamilyCounts{},
		BaitGuardAvailable: false,
	}

	values := policy.ActionValues(state)
	if values.BestAction != Forzar {
		t.Fatalf("expected Forzar best action, got %s with values %+v", values.BestAction, values.Values)
	}
	if values.Values[Soltar] >= values.Values[Forzar] {
		t.Fatalf("expected Soltar worse than Forzar, got %+v", values.Values)
	}
}

func TestRunBatchReportProducesPolicyAndActionComparisons(t *testing.T) {
	report, err := RunBatchReport([]BatchScenario{{
		Name:     "test-scenario",
		Category: "baseline",
		Config: CombatConfig{
			InitialTrackPos: 4,
			Fish:            NewBlackFish(),
		},
	}}, BatchOptions{Simulations: 20, Seed: 11})
	if err != nil {
		t.Fatalf("RunBatchReport() error = %v", err)
	}
	if len(report.Scenarios) != 1 {
		t.Fatalf("expected one scenario report, got %d", len(report.Scenarios))
	}
	entry := report.Scenarios[0]
	if _, ok := entry.PolicyResults["optimal"]; !ok {
		t.Fatalf("expected optimal policy result")
	}
	if _, ok := entry.PolicyResults["heuristic"]; !ok {
		t.Fatalf("expected heuristic policy result")
	}
	if _, ok := entry.PolicyResults["random"]; !ok {
		t.Fatalf("expected random policy result")
	}
	if entry.ForcedActionResults[Forzar].Policy != "forzar" {
		t.Fatalf("expected forced forzar policy, got %q", entry.ForcedActionResults[Forzar].Policy)
	}
	if entry.BestForcedRate < 0 {
		t.Fatalf("expected non-negative best forced rate, got %f", entry.BestForcedRate)
	}
}

type fixedPolicy struct {
	action Action
}

func (p fixedPolicy) ChooseAction(BeliefState) Action {
	return p.action
}

func (p fixedPolicy) Name() string {
	return "fixed"
}
