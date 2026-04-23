package cards

import "pesca/internal/domain"

type Owner string

const (
	OwnerFish   Owner = "fish"
	OwnerPlayer Owner = "player"
)

type Trigger int

const (
	TriggerOnOwnerWin Trigger = iota
	TriggerOnOwnerLose
	TriggerOnRoundDraw
)

type EffectContext struct {
	Owner   Owner
	Outcome domain.RoundOutcome
}

type CardEffect struct {
	Trigger       Trigger
	DistanceShift int
	DepthShift    int
}

func (effect CardEffect) Applies(context EffectContext) bool {
	switch effect.Trigger {
	case TriggerOnOwnerWin:
		return context.Owner == OwnerFish && context.Outcome == domain.FishWin ||
			context.Owner == OwnerPlayer && context.Outcome == domain.PlayerWin
	case TriggerOnOwnerLose:
		return context.Owner == OwnerFish && context.Outcome == domain.PlayerWin ||
			context.Owner == OwnerPlayer && context.Outcome == domain.FishWin
	case TriggerOnRoundDraw:
		return context.Outcome == domain.Draw
	default:
		return false
	}
}

func FilterEffects(effects []CardEffect, context EffectContext) []CardEffect {
	filteredEffects := make([]CardEffect, 0, len(effects))
	for _, effect := range effects {
		if !effect.Applies(context) {
			continue
		}

		filteredEffects = append(filteredEffects, effect)
	}

	return filteredEffects
}

type FishCard struct {
	Move    domain.Move
	Effects []CardEffect
}

type PlayerCard struct {
	Move    domain.Move
	Effects []CardEffect
}

func NewFishCard(move domain.Move, effects ...CardEffect) FishCard {
	return FishCard{
		Move:    move,
		Effects: append([]CardEffect(nil), effects...),
	}
}

func NewPlayerCard(move domain.Move, effects ...CardEffect) PlayerCard {
	return PlayerCard{
		Move:    move,
		Effects: append([]CardEffect(nil), effects...),
	}
}
