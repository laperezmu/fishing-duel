package rules

import "pesca/internal/domain"

type ClassicEvaluator struct {
	FishCombatProfile FishCombatProfile
	CombatConditions  []CombatCondition
}

func NewClassicEvaluator(fishCombatProfile FishCombatProfile, combatConditions ...CombatCondition) ClassicEvaluator {
	configuredConditions := append([]CombatCondition{TieAdvantageCondition{}}, combatConditions...)
	return ClassicEvaluator{
		FishCombatProfile: fishCombatProfile,
		CombatConditions:  configuredConditions,
	}
}

func (e ClassicEvaluator) Evaluate(playerMove, fishMove domain.Move) domain.RoundOutcome {
	context := CombatContext{
		PlayerMove:  playerMove,
		FishMove:    fishMove,
		FishProfile: e.FishCombatProfile,
	}

	for _, combatCondition := range e.CombatConditions {
		if outcome, ok := combatCondition.Apply(context); ok {
			return outcome
		}
	}

	return evaluateClassicBase(playerMove, fishMove)
}

func evaluateClassicBase(playerMove, fishMove domain.Move) domain.RoundOutcome {
	if playerMove == fishMove {
		return domain.Draw
	}

	if (playerMove == domain.Blue && fishMove == domain.Red) ||
		(playerMove == domain.Red && fishMove == domain.Yellow) ||
		(playerMove == domain.Yellow && fishMove == domain.Blue) {
		return domain.PlayerWin
	}

	return domain.FishWin
}
