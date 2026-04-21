package rules

import (
	"pesca/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClassicEvaluatorEvaluate(t *testing.T) {
	tests := []struct {
		title   string
		profile FishCombatProfile
		player  domain.Move
		fish    domain.Move
		outcome domain.RoundOutcome
	}{
		{title: "returns player win when the player uses blue against fish red", profile: NewFishCombatProfile(), player: domain.Blue, fish: domain.Red, outcome: domain.PlayerWin},
		{title: "returns player win when the player uses red against fish yellow", profile: NewFishCombatProfile(), player: domain.Red, fish: domain.Yellow, outcome: domain.PlayerWin},
		{title: "returns player win when the player uses yellow against fish blue", profile: NewFishCombatProfile(), player: domain.Yellow, fish: domain.Blue, outcome: domain.PlayerWin},
		{title: "returns draw when both moves match without tie advantage", profile: NewFishCombatProfile(), player: domain.Blue, fish: domain.Blue, outcome: domain.Draw},
		{title: "returns fish win when the fish has the stronger matchup", profile: NewFishCombatProfile(), player: domain.Red, fish: domain.Blue, outcome: domain.FishWin},
		{title: "returns fish win when blue tie advantage is configured", profile: NewFishCombatProfile(domain.Blue, domain.Red), player: domain.Blue, fish: domain.Blue, outcome: domain.FishWin},
		{title: "returns fish win when red tie advantage is configured", profile: NewFishCombatProfile(domain.Blue, domain.Red), player: domain.Red, fish: domain.Red, outcome: domain.FishWin},
		{title: "returns draw when the tie color is outside the configured advantage set", profile: NewFishCombatProfile(domain.Blue, domain.Red), player: domain.Yellow, fish: domain.Yellow, outcome: domain.Draw},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			evaluator := NewClassicEvaluator(test.profile)
			assert.Equal(t, test.outcome, evaluator.Evaluate(test.player, test.fish))
		})
	}
}
