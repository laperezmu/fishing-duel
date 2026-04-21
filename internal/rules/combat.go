package rules

import "pesca/internal/domain"

type FishCombatProfile struct {
	TieAdvantageMoves []domain.Move
}

func NewFishCombatProfile(tieAdvantageMoves ...domain.Move) FishCombatProfile {
	moves := append([]domain.Move(nil), tieAdvantageMoves...)
	return FishCombatProfile{TieAdvantageMoves: moves}
}

func (p FishCombatProfile) HasTieAdvantage(move domain.Move) bool {
	for _, candidate := range p.TieAdvantageMoves {
		if candidate == move {
			return true
		}
	}

	return false
}

type CombatContext struct {
	PlayerMove  domain.Move
	FishMove    domain.Move
	FishProfile FishCombatProfile
}

type CombatCondition interface {
	Apply(CombatContext) (domain.RoundOutcome, bool)
}
