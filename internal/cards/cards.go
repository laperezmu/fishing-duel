package cards

import (
	"fmt"
	"pesca/internal/domain"
	"sort"
)

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
	TriggerOnCardUsed
	TriggerOnFishSplash
	TriggerOnDiscard
	TriggerOnFishReshuffle
	TriggerOnFishExhausted
	TriggerOnColorDraw
	TriggerOnOwnerColorWin
	TriggerOnOwnerColorLose
)

type EffectType string

const (
	EffectTypeUnknown                  EffectType = ""
	EffectTypeAdvanceHorizontal        EffectType = "advance_horizontal"
	EffectTypeAdvanceVertical          EffectType = "advance_vertical"
	EffectTypeReshuffleCurrentDeck     EffectType = "reshuffle_current_deck"
	EffectTypeForceDiscardColor        EffectType = "force_discard_color"
	EffectTypeInstantExhaustFish       EffectType = "instant_exhaust_fish"
	EffectTypeExhaustPlayerColor       EffectType = "exhaust_player_color"
	EffectTypeReorderPlayerDiscard     EffectType = "reorder_player_discard"
	EffectTypeHideDiscardTemporary     EffectType = "hide_discard_temporary"
	EffectTypeSuccessfulSplashApproach EffectType = "successful_splash_approach"
	EffectTypeLegacyCaptureWindow      EffectType = "legacy_capture_window"
	EffectTypeLegacySurfaceWindow      EffectType = "legacy_surface_window"
	EffectTypeLegacyExhaustionWindow   EffectType = "legacy_exhaustion_window"
)

type EffectContext struct {
	Owner          Owner
	Phase          Phase
	Outcome        domain.RoundOutcome
	CardMove       domain.Move
	TriggerMove    domain.Move
	HasTriggerMove bool
	InDiscard      bool
	FishDidSplash  bool
	FishReshuffled bool
	FishExhausted  bool
}

type CardEffect struct {
	Trigger    Trigger
	Type       EffectType
	Priority   int
	TargetMove domain.Move

	DistanceShift int
	DepthShift    int

	CaptureDistanceBonus           int
	ExhaustionCaptureDistanceBonus int
	SurfaceDepthBonus              int
}

type OwnedEffect struct {
	Owner  Owner
	Effect CardEffect
}

type DiscardVisibility string

const (
	DiscardVisibilityFull     DiscardVisibility = "full"
	DiscardVisibilityMoveOnly DiscardVisibility = "move_only"
	DiscardVisibilityMasked   DiscardVisibility = "masked"
	DiscardVisibilityHidden   DiscardVisibility = "hidden"
)

func (effect CardEffect) Applies(context EffectContext) bool {
	if effect.appliesInDrawPhase(context) {
		return true
	}
	if effect.appliesInOutcomePhase(context) {
		return true
	}
	return effect.appliesInStatePhase(context)
}

func (effect CardEffect) appliesInDrawPhase(context EffectContext) bool {
	switch effect.Trigger {
	case TriggerOnDraw:
		return context.Phase == PhaseDraw
	case TriggerOnCardUsed:
		return context.Phase == PhaseDraw
	default:
		return false
	}
}

func (effect CardEffect) appliesInOutcomePhase(context EffectContext) bool {
	resolver, ok := outcomeTriggerMatchers[effect.Trigger]
	if !ok {
		return false
	}
	return resolver(effect, context)
}

func (effect CardEffect) appliesInStatePhase(context EffectContext) bool {
	switch effect.Trigger {
	case TriggerOnDiscard:
		return context.InDiscard
	default:
		return false
	}
}

var outcomeTriggerMatchers = map[Trigger]func(CardEffect, EffectContext) bool{
	TriggerOnOwnerWin: func(effect CardEffect, context EffectContext) bool {
		return effect.appliesOnOwnerWin(context)
	},
	TriggerOnOwnerLose: func(effect CardEffect, context EffectContext) bool {
		return effect.appliesOnOwnerLose(context)
	},
	TriggerOnRoundDraw: func(effect CardEffect, context EffectContext) bool {
		return effect.appliesOnRoundDraw(context)
	},
	TriggerOnFishSplash: func(_ CardEffect, context EffectContext) bool {
		return context.Phase == PhaseOutcome && context.FishDidSplash
	},
	TriggerOnFishReshuffle: func(_ CardEffect, context EffectContext) bool {
		return context.Phase == PhaseOutcome && context.FishReshuffled
	},
	TriggerOnFishExhausted: func(_ CardEffect, context EffectContext) bool {
		return context.Phase == PhaseOutcome && context.FishExhausted
	},
	TriggerOnColorDraw: func(effect CardEffect, context EffectContext) bool {
		return effect.appliesOnColorDraw(context)
	},
	TriggerOnOwnerColorWin: func(effect CardEffect, context EffectContext) bool {
		return effect.appliesOnOwnerWin(context) && effect.matchesTriggerMove(context)
	},
	TriggerOnOwnerColorLose: func(effect CardEffect, context EffectContext) bool {
		return effect.appliesOnOwnerLose(context) && effect.matchesTriggerMove(context)
	},
}

func FilterEffects(effects []CardEffect, context EffectContext) []CardEffect {
	filteredEffects := make([]CardEffect, 0, len(effects))
	for _, effect := range effects {
		effect = effect.Normalize()
		if !effect.Applies(context) {
			continue
		}

		filteredEffects = append(filteredEffects, effect)
	}

	return filteredEffects
}

func FilterOwnedEffects(effects []CardEffect, context EffectContext) []OwnedEffect {
	filteredEffects := FilterEffects(effects, context)
	owned := make([]OwnedEffect, 0, len(filteredEffects))
	for _, effect := range filteredEffects {
		owned = append(owned, OwnedEffect{Owner: context.Owner, Effect: effect})
	}

	return owned
}

func OrderOwnedEffects(effects []OwnedEffect) []OwnedEffect {
	ordered := append([]OwnedEffect(nil), effects...)
	sort.SliceStable(ordered, func(i, j int) bool {
		left := ordered[i].Effect.Normalize()
		right := ordered[j].Effect.Normalize()
		if left.Priority != right.Priority {
			return left.Priority > right.Priority
		}
		if ordered[i].Owner != ordered[j].Owner {
			return ordered[i].Owner == OwnerFish
		}
		return false
	})

	return ordered
}

func FlattenOwnedEffects(effects []OwnedEffect) []CardEffect {
	flattened := make([]CardEffect, 0, len(effects))
	for _, effect := range effects {
		flattened = append(flattened, effect.Effect.Normalize())
	}

	return flattened
}

func (effect CardEffect) Normalize() CardEffect {
	resolved := effect
	if resolved.Type == EffectTypeUnknown {
		resolved.Type = resolved.inferType()
	}
	if resolved.Priority == 0 {
		resolved.Priority = resolved.defaultPriority()
	}
	return resolved
}

func (effect CardEffect) Validate() error {
	resolved := effect.Normalize()
	if resolved.Type == EffectTypeUnknown {
		return fmt.Errorf("effect type is required")
	}
	if resolved.Priority < 0 {
		return fmt.Errorf("effect priority must be greater than or equal to 0")
	}
	if resolved.Trigger < TriggerOnDraw || resolved.Trigger > TriggerOnOwnerColorLose {
		return fmt.Errorf("unknown effect trigger %d", resolved.Trigger)
	}
	return nil
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
	normalizedEffects := make([]CardEffect, 0, len(effects))
	for _, effect := range effects {
		normalizedEffects = append(normalizedEffects, effect.Normalize())
	}
	return FishCard{
		Move:              move,
		Effects:           normalizedEffects,
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
	normalizedEffects := make([]CardEffect, 0, len(effects))
	for _, effect := range effects {
		normalizedEffects = append(normalizedEffects, effect.Normalize())
	}
	return PlayerCard{
		Move:    move,
		Effects: normalizedEffects,
	}
}

func NewNamedPlayerCard(name, summary string, move domain.Move, effects ...CardEffect) PlayerCard {
	card := NewPlayerCard(move, effects...)
	card.Name = name
	card.Summary = summary
	return card
}

func CloneFishCard(card FishCard) FishCard {
	effects := make([]CardEffect, 0, len(card.Effects))
	for _, effect := range card.Effects {
		effects = append(effects, effect.Normalize())
	}
	clonedCard := NewFishCard(card.Move, effects...)
	clonedCard.Name = card.Name
	clonedCard.Summary = card.Summary
	clonedCard.DiscardVisibility = card.effectiveDiscardVisibility()
	return clonedCard
}

func ClonePlayerCard(card PlayerCard) PlayerCard {
	effects := make([]CardEffect, 0, len(card.Effects))
	for _, effect := range card.Effects {
		effects = append(effects, effect.Normalize())
	}
	clonedCard := NewPlayerCard(card.Move, effects...)
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

func (effect CardEffect) appliesOnColorDraw(context EffectContext) bool {
	return effect.appliesOnRoundDraw(context) && effect.matchesTriggerMove(context)
}

func (effect CardEffect) matchesTriggerMove(context EffectContext) bool {
	return context.HasTriggerMove && context.CardMove == context.TriggerMove
}

func (effect CardEffect) inferType() EffectType {
	switch {
	case effect.DistanceShift != 0:
		return EffectTypeAdvanceHorizontal
	case effect.DepthShift != 0:
		return EffectTypeAdvanceVertical
	case effect.CaptureDistanceBonus != 0:
		return EffectTypeLegacyCaptureWindow
	case effect.SurfaceDepthBonus != 0:
		return EffectTypeLegacySurfaceWindow
	case effect.ExhaustionCaptureDistanceBonus != 0:
		return EffectTypeLegacyExhaustionWindow
	default:
		return EffectTypeUnknown
	}
}

func (effect CardEffect) defaultPriority() int {
	effectType := effect.Type
	if effectType == EffectTypeUnknown {
		effectType = effect.inferType()
	}
	switch effectType {
	case EffectTypeReshuffleCurrentDeck, EffectTypeInstantExhaustFish, EffectTypeExhaustPlayerColor:
		return 90
	case EffectTypeHideDiscardTemporary, EffectTypeReorderPlayerDiscard, EffectTypeForceDiscardColor:
		return 80
	case EffectTypeSuccessfulSplashApproach:
		return 70
	case EffectTypeLegacyCaptureWindow, EffectTypeLegacySurfaceWindow, EffectTypeLegacyExhaustionWindow:
		return 60
	case EffectTypeAdvanceHorizontal, EffectTypeAdvanceVertical:
		return 50
	default:
		return 10
	}
}

func (card FishCard) effectiveDiscardVisibility() DiscardVisibility {
	return NormalizeDiscardVisibility(card.DiscardVisibility)
}

func (card FishCard) EffectiveDiscardVisibility() DiscardVisibility {
	return card.effectiveDiscardVisibility()
}

func NormalizeDiscardVisibility(visibility DiscardVisibility) DiscardVisibility {
	if visibility == "" {
		return DiscardVisibilityFull
	}

	return visibility
}
