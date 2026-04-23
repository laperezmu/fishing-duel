package rules

import (
	"pesca/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClassicEvaluatorEvaluate(t *testing.T) {
	tests := []struct {
		title             string
		fishCombatProfile FishCombatProfile
		outcomeHooks      []OutcomeHook
		playerMove        domain.Move
		fishMove          domain.Move
		wantOutcome       domain.RoundOutcome
	}{
		{title: "returns player win when the player uses blue against fish red", fishCombatProfile: NewFishCombatProfile(), playerMove: domain.Blue, fishMove: domain.Red, wantOutcome: domain.PlayerWin},
		{title: "returns player win when the player uses red against fish yellow", fishCombatProfile: NewFishCombatProfile(), playerMove: domain.Red, fishMove: domain.Yellow, wantOutcome: domain.PlayerWin},
		{title: "returns player win when the player uses yellow against fish blue", fishCombatProfile: NewFishCombatProfile(), playerMove: domain.Yellow, fishMove: domain.Blue, wantOutcome: domain.PlayerWin},
		{title: "returns draw when both moves match without tie advantage", fishCombatProfile: NewFishCombatProfile(), playerMove: domain.Blue, fishMove: domain.Blue, wantOutcome: domain.Draw},
		{title: "returns fish win when the fish has the stronger matchup", fishCombatProfile: NewFishCombatProfile(), playerMove: domain.Red, fishMove: domain.Blue, wantOutcome: domain.FishWin},
		{title: "returns fish win when blue tie advantage is configured", fishCombatProfile: NewFishCombatProfile(domain.Blue, domain.Red), playerMove: domain.Blue, fishMove: domain.Blue, wantOutcome: domain.FishWin},
		{title: "returns fish win when red tie advantage is configured", fishCombatProfile: NewFishCombatProfile(domain.Blue, domain.Red), playerMove: domain.Red, fishMove: domain.Red, wantOutcome: domain.FishWin},
		{title: "returns draw when the tie color is outside the configured advantage set", fishCombatProfile: NewFishCombatProfile(domain.Blue, domain.Red), playerMove: domain.Yellow, fishMove: domain.Yellow, wantOutcome: domain.Draw},
		{
			title:             "applies outcome hooks after resolving the base combat outcome",
			fishCombatProfile: NewFishCombatProfile(),
			outcomeHooks: []OutcomeHook{outcomeHookFunc(func(context CombatContext, currentOutcome domain.RoundOutcome) (domain.RoundOutcome, bool) {
				assert.Equal(t, domain.Blue, context.PlayerMove)
				assert.Equal(t, domain.Red, context.FishMove)
				assert.Equal(t, domain.PlayerWin, currentOutcome)
				return domain.Draw, true
			})},
			playerMove:  domain.Blue,
			fishMove:    domain.Red,
			wantOutcome: domain.Draw,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			evaluator := NewClassicEvaluator(test.fishCombatProfile).WithOutcomeHooks(test.outcomeHooks...)
			assert.Equal(t, test.wantOutcome, evaluator.Evaluate(test.playerMove, test.fishMove))
		})
	}
}

func TestClassicEvaluatorWithOutcomeHooks(t *testing.T) {
	hookA := outcomeHookFunc(func(CombatContext, domain.RoundOutcome) (domain.RoundOutcome, bool) {
		return domain.Draw, true
	})
	hookB := outcomeHookFunc(func(CombatContext, domain.RoundOutcome) (domain.RoundOutcome, bool) {
		return domain.FishWin, true
	})

	evaluator := NewClassicEvaluator(NewFishCombatProfile()).WithOutcomeHooks(hookA, hookB)

	assert.Len(t, evaluator.OutcomeHooks, 2)
	assert.Empty(t, NewClassicEvaluator(NewFishCombatProfile()).OutcomeHooks)
	assert.Equal(t, domain.FishWin, evaluator.Evaluate(domain.Blue, domain.Red))
}

type outcomeHookFunc func(CombatContext, domain.RoundOutcome) (domain.RoundOutcome, bool)

func (hook outcomeHookFunc) Apply(context CombatContext, currentOutcome domain.RoundOutcome) (domain.RoundOutcome, bool) {
	return hook(context, currentOutcome)
}
