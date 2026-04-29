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

type DiscardVisibility string

const (
	DiscardVisibilityFull     DiscardVisibility = "full"
	DiscardVisibilityMoveOnly DiscardVisibility = "move_only"
	DiscardVisibilityMasked   DiscardVisibility = "masked"
	DiscardVisibilityHidden   DiscardVisibility = "hidden"
)

func (effect CardEffect) Applies(context EffectContext) bool {
	switch effect.Trigger {
	case TriggerOnDraw:
		return context.Phase == PhaseDraw
	case TriggerOnOwnerWin:
		return effect.appliesOnOwnerWin(context)
	case TriggerOnOwnerLose:
		return effect.appliesOnOwnerLose(context)
	case TriggerOnRoundDraw:
		return effect.appliesOnRoundDraw(context)
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
	Name              string
	Summary           string
	Move              domain.Move
	Effects           []CardEffect
	DiscardVisibility DiscardVisibility
}

type PlayerCard struct {
	Name    string
	Summary string
	Move    domain.Move
	Effects []CardEffect
}

func NewFishCard(move domain.Move, effects ...CardEffect) FishCard {
	return FishCard{
		Move:              move,
		Effects:           append([]CardEffect(nil), effects...),
		DiscardVisibility: DiscardVisibilityFull,
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
	clonedCard.DiscardVisibility = card.effectiveDiscardVisibility()
	return clonedCard
}

func ClonePlayerCard(card PlayerCard) PlayerCard {
	clonedCard := NewPlayerCard(card.Move, card.Effects...)
	clonedCard.Name = card.Name
	clonedCard.Summary = card.Summary
	return clonedCard
}

func (effect CardEffect) appliesOnOwnerWin(context EffectContext) bool {
	if context.Phase != PhaseOutcome {
		return false
	}

	if context.Owner == OwnerFish {
		return context.Outcome == domain.FishWin
	}

	return context.Owner == OwnerPlayer && context.Outcome == domain.PlayerWin
}

func (effect CardEffect) appliesOnOwnerLose(context EffectContext) bool {
	if context.Phase != PhaseOutcome {
		return false
	}

	if context.Owner == OwnerFish {
		return context.Outcome == domain.PlayerWin
	}

	return context.Owner == OwnerPlayer && context.Outcome == domain.FishWin
}

func (effect CardEffect) appliesOnRoundDraw(context EffectContext) bool {
	return context.Phase == PhaseOutcome && context.Outcome == domain.Draw
}

func (card FishCard) effectiveDiscardVisibility() DiscardVisibility {
	if card.DiscardVisibility == "" {
		return DiscardVisibilityFull
	}

	return card.DiscardVisibility
}

func (card FishCard) EffectiveDiscardVisibility() DiscardVisibility {
	return card.effectiveDiscardVisibility()
}
