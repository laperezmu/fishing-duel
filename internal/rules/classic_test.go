package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"pesca/internal/domain"
)

func TestClassicEvaluator(t *testing.T) {
	tests := []struct {
		name    string
		profile FishCombatProfile
		player  domain.Move
		fish    domain.Move
		outcome domain.RoundOutcome
	}{
		{name: "blue beats red", profile: NewFishCombatProfile(), player: domain.Blue, fish: domain.Red, outcome: domain.PlayerWin},
		{name: "red beats yellow", profile: NewFishCombatProfile(), player: domain.Red, fish: domain.Yellow, outcome: domain.PlayerWin},
		{name: "yellow beats blue", profile: NewFishCombatProfile(), player: domain.Yellow, fish: domain.Blue, outcome: domain.PlayerWin},
		{name: "same move draws without tie advantage", profile: NewFishCombatProfile(), player: domain.Blue, fish: domain.Blue, outcome: domain.Draw},
		{name: "fish wins opposite matchup", profile: NewFishCombatProfile(), player: domain.Red, fish: domain.Blue, outcome: domain.FishWin},
		{name: "fish wins blue tie with configured advantage", profile: NewFishCombatProfile(domain.Blue, domain.Red), player: domain.Blue, fish: domain.Blue, outcome: domain.FishWin},
		{name: "fish wins red tie with configured advantage", profile: NewFishCombatProfile(domain.Blue, domain.Red), player: domain.Red, fish: domain.Red, outcome: domain.FishWin},
		{name: "tie stays draw outside configured advantage", profile: NewFishCombatProfile(domain.Blue, domain.Red), player: domain.Yellow, fish: domain.Yellow, outcome: domain.Draw},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaluator := NewClassicEvaluator(test.profile)
			assert.Equal(t, test.outcome, evaluator.Evaluate(test.player, test.fish))
		})
	}
}
