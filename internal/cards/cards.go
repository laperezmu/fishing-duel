package cards

import "pesca/internal/domain"

type Owner string

const (
	OwnerFish   Owner = "fish"
	OwnerPlayer Owner = "player"
)

type Phase int

const (
	PhaseDraw Phase = iota
	PhaseOutcome
)

type Trigger int

const (
	TriggerOnDraw Trigger = iota
	TriggerOnOwnerWin
	TriggerOnOwnerLose
	TriggerOnRoundDraw
)

type EffectContext struct {
	Owner   Owner
	Phase   Phase
	Outcome domain.RoundOutcome
}

type CardEffect struct {
	Trigger       Trigger
	DistanceShift int
	DepthShift    int

	CaptureDistanceBonus           int
	ExhaustionCaptureDistanceBonus int
	SurfaceDepthBonus              int
}

func (effect CardEffect) Applies(context EffectContext) bool {
	switch effect.Trigger {
	case TriggerOnDraw:
		return context.Phase == PhaseDraw
	case TriggerOnOwnerWin:
		if context.Phase != PhaseOutcome {
			return false
		}
		return context.Owner == OwnerFish && context.Outcome == domain.FishWin ||
			context.Owner == OwnerPlayer && context.Outcome == domain.PlayerWin
	case TriggerOnOwnerLose:
		if context.Phase != PhaseOutcome {
			return false
		}
		return context.Owner == OwnerFish && context.Outcome == domain.PlayerWin ||
			context.Owner == OwnerPlayer && context.Outcome == domain.FishWin
	case TriggerOnRoundDraw:
		if context.Phase != PhaseOutcome {
			return false
		}
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
	Name    string
	Summary string
	Move    domain.Move
	Effects []CardEffect
}

type PlayerCard struct {
	Name    string
	Summary string
	Move    domain.Move
	Effects []CardEffect
}

func NewFishCard(move domain.Move, effects ...CardEffect) FishCard {
	return FishCard{
		Move:    move,
		Effects: append([]CardEffect(nil), effects...),
	}
}

func NewNamedFishCard(name, summary string, move domain.Move, effects ...CardEffect) FishCard {
	card := NewFishCard(move, effects...)
	card.Name = name
	card.Summary = summary
	return card
}

func NewPlayerCard(move domain.Move, effects ...CardEffect) PlayerCard {
	return PlayerCard{
		Move:    move,
		Effects: append([]CardEffect(nil), effects...),
	}
}

func NewNamedPlayerCard(name, summary string, move domain.Move, effects ...CardEffect) PlayerCard {
	card := NewPlayerCard(move, effects...)
	card.Name = name
	card.Summary = summary
	return card
}

func CloneFishCard(card FishCard) FishCard {
	clonedCard := NewFishCard(card.Move, card.Effects...)
	clonedCard.Name = card.Name
	clonedCard.Summary = card.Summary
	return clonedCard
}

func ClonePlayerCard(card PlayerCard) PlayerCard {
	clonedCard := NewPlayerCard(card.Move, card.Effects...)
	clonedCard.Name = card.Name
	clonedCard.Summary = card.Summary
	return clonedCard
}
