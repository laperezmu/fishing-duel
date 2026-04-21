package rules

import "pesca/internal/domain"

type ClassicEvaluator struct {
	Fish       FishCombatProfile
	Conditions []Condition
}

func NewClassicEvaluator(fish FishCombatProfile, conditions ...Condition) ClassicEvaluator {
	configuredConditions := append([]Condition{TieAdvantageCondition{}}, conditions...)
	return ClassicEvaluator{
		Fish:       fish,
		Conditions: configuredConditions,
	}
}

func (e ClassicEvaluator) Evaluate(player, fish domain.Move) domain.RoundOutcome {
	context := CombatContext{
		PlayerMove: player,
		FishMove:   fish,
		Fish:       e.Fish,
	}

	for _, condition := range e.Conditions {
		if outcome, ok := condition.Apply(context); ok {
			return outcome
		}
	}

	return evaluateClassicBase(player, fish)
}

func evaluateClassicBase(player, fish domain.Move) domain.RoundOutcome {
	if player == fish {
		return domain.Draw
	}

	if (player == domain.Blue && fish == domain.Red) ||
		(player == domain.Red && fish == domain.Yellow) ||
		(player == domain.Yellow && fish == domain.Blue) {
		return domain.PlayerWin
	}

	return domain.FishWin
}
