package rules

import (
	"testing"

	"pesca/internal/domain"
)

func TestClassicEvaluator(t *testing.T) {
	tests := []struct {
		name    string
		player  domain.Move
		fish    domain.Move
		outcome domain.RoundOutcome
	}{
		{name: "blue beats red", player: domain.Blue, fish: domain.Red, outcome: domain.PlayerWin},
		{name: "red beats yellow", player: domain.Red, fish: domain.Yellow, outcome: domain.PlayerWin},
		{name: "yellow beats blue", player: domain.Yellow, fish: domain.Blue, outcome: domain.PlayerWin},
		{name: "same move draws", player: domain.Blue, fish: domain.Blue, outcome: domain.Draw},
		{name: "fish wins opposite matchup", player: domain.Red, fish: domain.Blue, outcome: domain.FishWin},
	}

	evaluator := NewClassicEvaluator(NewFishCombatProfile())
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := evaluator.Evaluate(test.player, test.fish); got != test.outcome {
				t.Fatalf("Evaluate(%v, %v) = %v, want %v", test.player, test.fish, got, test.outcome)
			}
		})
	}
}

func TestClassicEvaluatorFishWinsTieWithConfiguredColorAdvantage(t *testing.T) {
	evaluator := NewClassicEvaluator(NewFishCombatProfile(domain.Blue, domain.Red))

	if got := evaluator.Evaluate(domain.Blue, domain.Blue); got != domain.FishWin {
		t.Fatalf("Evaluate(blue, blue) = %v, want %v", got, domain.FishWin)
	}

	if got := evaluator.Evaluate(domain.Red, domain.Red); got != domain.FishWin {
		t.Fatalf("Evaluate(red, red) = %v, want %v", got, domain.FishWin)
	}
}

func TestClassicEvaluatorKeepsDrawForTieOutsideConfiguredColorAdvantage(t *testing.T) {
	evaluator := NewClassicEvaluator(NewFishCombatProfile(domain.Blue, domain.Red))

	if got := evaluator.Evaluate(domain.Yellow, domain.Yellow); got != domain.Draw {
		t.Fatalf("Evaluate(yellow, yellow) = %v, want %v", got, domain.Draw)
	}
}
