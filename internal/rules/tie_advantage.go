package rules

import "pesca/internal/domain"

type TieAdvantageCondition struct{}

func (TieAdvantageCondition) Apply(context CombatContext) (domain.RoundOutcome, bool) {
	if context.PlayerMove != context.FishMove {
		return 0, false
	}

	if !context.Fish.HasTieAdvantage(context.FishMove) {
		return 0, false
	}

	return domain.FishWin, true
}
