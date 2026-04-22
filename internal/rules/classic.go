package rules

import "pesca/internal/domain"

type ClassicEvaluator struct {
	FishCombatProfile FishCombatProfile
	CombatConditions  []CombatCondition
	OutcomeHooks      []OutcomeHook
}

func NewClassicEvaluator(fishCombatProfile FishCombatProfile, combatConditions ...CombatCondition) ClassicEvaluator {
	configuredConditions := append([]CombatCondition{TieAdvantageCondition{}}, combatConditions...)
	return ClassicEvaluator{
		FishCombatProfile: fishCombatProfile,
		CombatConditions:  configuredConditions,
	}
}

func (e ClassicEvaluator) WithOutcomeHooks(outcomeHooks ...OutcomeHook) ClassicEvaluator {
	e.OutcomeHooks = append([]OutcomeHook(nil), outcomeHooks...)
	return e
}

func (e ClassicEvaluator) Evaluate(playerMove, fishMove domain.Move) domain.RoundOutcome {
	context := CombatContext{
		PlayerMove:  playerMove,
		FishMove:    fishMove,
		FishProfile: e.FishCombatProfile,
	}

	for _, combatCondition := range e.CombatConditions {
		if outcome, ok := combatCondition.Apply(context); ok {
			return e.applyOutcomeHooks(context, outcome)
		}
	}

	return e.applyOutcomeHooks(context, evaluateClassicBase(playerMove, fishMove))
}

func (e ClassicEvaluator) applyOutcomeHooks(context CombatContext, outcome domain.RoundOutcome) domain.RoundOutcome {
	updatedOutcome := outcome
	for _, outcomeHook := range e.OutcomeHooks {
		nextOutcome, ok := outcomeHook.Apply(context, updatedOutcome)
		if !ok {
			continue
		}

		updatedOutcome = nextOutcome
	}

	return updatedOutcome
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
