package rules

import "pesca/internal/domain"

type ClassicEvaluator struct{}

func (ClassicEvaluator) Evaluate(player, fish domain.Move) domain.RoundOutcome {
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
