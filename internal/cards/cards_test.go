package cards

import (
	"pesca/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCardEffectApplies(t *testing.T) {
	tests := []struct {
		title   string
		effect  CardEffect
		context EffectContext
		want    bool
	}{
		{
			title:  "matches draw effects during the draw phase",
			effect: CardEffect{Trigger: TriggerOnDraw},
			context: EffectContext{
				Owner: OwnerFish,
				Phase: PhaseDraw,
			},
			want: true,
		},
		{
			title:  "matches owner win for fish owned cards",
			effect: CardEffect{Trigger: TriggerOnOwnerWin},
			context: EffectContext{
				Owner:   OwnerFish,
				Phase:   PhaseOutcome,
				Outcome: domain.FishWin,
			},
			want: true,
		},
		{
			title:  "matches owner lose for fish owned cards",
			effect: CardEffect{Trigger: TriggerOnOwnerLose},
			context: EffectContext{
				Owner:   OwnerFish,
				Phase:   PhaseOutcome,
				Outcome: domain.PlayerWin,
			},
			want: true,
		},
		{
			title:  "matches owner win for player owned cards",
			effect: CardEffect{Trigger: TriggerOnOwnerWin},
			context: EffectContext{
				Owner:   OwnerPlayer,
				Phase:   PhaseOutcome,
				Outcome: domain.PlayerWin,
			},
			want: true,
		},
		{
			title:  "matches round draw regardless of owner",
			effect: CardEffect{Trigger: TriggerOnRoundDraw},
			context: EffectContext{
				Owner:   OwnerPlayer,
				Phase:   PhaseOutcome,
				Outcome: domain.Draw,
			},
			want: true,
		},
		{
			title:  "does not match when the owner relative outcome does not apply",
			effect: CardEffect{Trigger: TriggerOnOwnerLose},
			context: EffectContext{
				Owner:   OwnerPlayer,
				Phase:   PhaseOutcome,
				Outcome: domain.PlayerWin,
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			assert.Equal(t, test.want, test.effect.Applies(test.context))
		})
	}
}

func TestFilterEffects(t *testing.T) {
	effects := []CardEffect{{
		Trigger:       TriggerOnOwnerWin,
		DistanceShift: 1,
	}, {
		Trigger:    TriggerOnOwnerLose,
		DepthShift: 1,
	}, {
		Trigger:       TriggerOnRoundDraw,
		DistanceShift: -1,
	}, {
		Trigger:              TriggerOnDraw,
		CaptureDistanceBonus: 1,
	}}

	filteredEffects := FilterEffects(effects, EffectContext{Owner: OwnerFish, Phase: PhaseOutcome, Outcome: domain.PlayerWin})

	require.Len(t, filteredEffects, 1)
	assert.Equal(t, effects[1].Normalize(), filteredEffects[0])
}

func TestFilterEffectsForDrawPhase(t *testing.T) {
	effects := []CardEffect{{
		Trigger:              TriggerOnDraw,
		CaptureDistanceBonus: 1,
	}, {
		Trigger:    TriggerOnOwnerWin,
		DepthShift: 1,
	}}

	filteredEffects := FilterEffects(effects, EffectContext{Owner: OwnerFish, Phase: PhaseDraw})

	require.Len(t, filteredEffects, 1)
	assert.Equal(t, effects[0].Normalize(), filteredEffects[0])
}

func TestCardEffectNormalizeInfersTypeAndPriority(t *testing.T) {
	effect := CardEffect{Trigger: TriggerOnOwnerWin, DistanceShift: 1}

	normalized := effect.Normalize()

	assert.Equal(t, EffectTypeAdvanceHorizontal, normalized.Type)
	assert.Equal(t, 50, normalized.Priority)
}

func TestOrderOwnedEffectsPrefersHigherPriorityAndFishTies(t *testing.T) {
	effects := []OwnedEffect{{
		Owner:  OwnerPlayer,
		Effect: CardEffect{Trigger: TriggerOnOwnerWin, DistanceShift: -1, Priority: 40},
	}, {
		Owner:  OwnerFish,
		Effect: CardEffect{Trigger: TriggerOnOwnerWin, DepthShift: 1, Priority: 40},
	}, {
		Owner:  OwnerPlayer,
		Effect: CardEffect{Trigger: TriggerOnDraw, CaptureDistanceBonus: 1, Priority: 60},
	}}

	ordered := OrderOwnedEffects(effects)

	require.Len(t, ordered, 3)
	assert.Equal(t, 60, ordered[0].Effect.Priority)
	assert.Equal(t, OwnerFish, ordered[1].Owner)
	assert.Equal(t, OwnerPlayer, ordered[2].Owner)
}

func TestCardEffectAppliesColorSpecificTriggers(t *testing.T) {
	effect := CardEffect{Trigger: TriggerOnOwnerColorWin, DistanceShift: 1}
	context := EffectContext{
		Owner:          OwnerPlayer,
		Phase:          PhaseOutcome,
		Outcome:        domain.PlayerWin,
		CardMove:       domain.Red,
		TriggerMove:    domain.Red,
		HasTriggerMove: true,
	}

	assert.True(t, effect.Applies(context))
	assert.False(t, effect.Applies(EffectContext{
		Owner:          OwnerPlayer,
		Phase:          PhaseOutcome,
		Outcome:        domain.PlayerWin,
		CardMove:       domain.Blue,
		TriggerMove:    domain.Red,
		HasTriggerMove: true,
	}))
}

func TestNewFishCard(t *testing.T) {
	effect := CardEffect{Trigger: TriggerOnOwnerWin, DepthShift: 1}
	card := NewFishCard(domain.Red, effect)

	require.Len(t, card.Effects, 1)
	assert.Equal(t, domain.Red, card.Move)
	assert.Equal(t, effect.Normalize(), card.Effects[0])
}

func TestNewPlayerCard(t *testing.T) {
	effect := CardEffect{Trigger: TriggerOnOwnerLose, DistanceShift: 1}
	card := NewPlayerCard(domain.Blue, effect)

	require.Len(t, card.Effects, 1)
	assert.Equal(t, domain.Blue, card.Move)
	assert.Equal(t, effect.Normalize(), card.Effects[0])
}

func TestNewNamedPlayerCard(t *testing.T) {
	card := NewNamedPlayerCard("Anzuelo tenso", "Capturas desde un paso mas lejos este round.", domain.Blue, CardEffect{Trigger: TriggerOnDraw, CaptureDistanceBonus: 1})

	require.Len(t, card.Effects, 1)
	assert.Equal(t, "Anzuelo tenso", card.Name)
	assert.Equal(t, "Capturas desde un paso mas lejos este round.", card.Summary)
}

func TestCloneCardsPreserveMetadata(t *testing.T) {
	fishCard := NewNamedFishCard("Tiron de apertura", "Permite capturar desde un paso mas lejos este round.", domain.Red, CardEffect{Trigger: TriggerOnDraw, CaptureDistanceBonus: 1})
	playerCard := NewNamedPlayerCard("Anzuelo tenso", "Capturas desde un paso mas lejos este round.", domain.Blue, CardEffect{Trigger: TriggerOnDraw, CaptureDistanceBonus: 1})

	clonedFishCard := CloneFishCard(fishCard)
	clonedPlayerCard := ClonePlayerCard(playerCard)

	assert.Equal(t, fishCard, clonedFishCard)
	assert.Equal(t, playerCard, clonedPlayerCard)

	fishCard.Effects[0].CaptureDistanceBonus = 9
	playerCard.Effects[0].CaptureDistanceBonus = 9

	assert.Equal(t, 1, clonedFishCard.Effects[0].CaptureDistanceBonus)
	assert.Equal(t, 1, clonedPlayerCard.Effects[0].CaptureDistanceBonus)
}

func TestFilterOwnedEffects(t *testing.T) {
	t.Parallel()

	effects := []CardEffect{
		{Trigger: TriggerOnDraw, Type: EffectTypeAdvanceHorizontal},
		{Trigger: TriggerOnOwnerWin, Type: EffectTypeAdvanceVertical},
	}
	context := EffectContext{
		Owner: OwnerPlayer,
		Phase: PhaseDraw,
	}

	owned := FilterOwnedEffects(effects, context)

	require.Len(t, owned, 1)
	assert.Equal(t, OwnerPlayer, owned[0].Owner)
	assert.Equal(t, EffectTypeAdvanceHorizontal, owned[0].Effect.Type)
}

func TestFlattenOwnedEffects(t *testing.T) {
	t.Parallel()

	owned := []OwnedEffect{
		{Owner: OwnerFish, Effect: CardEffect{Type: EffectTypeAdvanceVertical, Priority: 10}},
		{Owner: OwnerPlayer, Effect: CardEffect{Type: EffectTypeAdvanceHorizontal, Priority: 20}},
	}

	flattened := FlattenOwnedEffects(owned)

	require.Len(t, flattened, 2)
	assert.Equal(t, EffectTypeAdvanceVertical, flattened[0].Type)
	assert.Equal(t, EffectTypeAdvanceHorizontal, flattened[1].Type)
}

func TestCardEffectValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		effect  CardEffect
		wantErr bool
	}{
		{
			name: "valid effect with type and trigger",
			effect: CardEffect{
				Type:     EffectTypeAdvanceHorizontal,
				Trigger:  TriggerOnDraw,
				Priority: 10,
			},
			wantErr: false,
		},
		{
			name: "missing type",
			effect: CardEffect{
				Trigger:  TriggerOnDraw,
				Priority: 10,
			},
			wantErr: true,
		},
		{
			name: "negative priority",
			effect: CardEffect{
				Type:     EffectTypeAdvanceHorizontal,
				Trigger:  TriggerOnDraw,
				Priority: -1,
			},
			wantErr: true,
		},
		{
			name: "unknown trigger",
			effect: CardEffect{
				Type:     EffectTypeAdvanceHorizontal,
				Trigger:  Trigger(999),
				Priority: 10,
			},
			wantErr: true,
		},
		{
			name: "priority defaults for advance vertical",
			effect: CardEffect{
				Type: EffectTypeAdvanceVertical,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.effect.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNormalizeDiscardVisibility(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input DiscardVisibility
		want  DiscardVisibility
	}{
		{"masked stays masked", DiscardVisibilityMasked, DiscardVisibilityMasked},
		{"full stays full", DiscardVisibilityFull, DiscardVisibilityFull},
		{"hidden stays hidden", DiscardVisibilityHidden, DiscardVisibilityHidden},
		{"move only stays move only", DiscardVisibilityMoveOnly, DiscardVisibilityMoveOnly},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := NormalizeDiscardVisibility(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFishCardEffectiveDiscardVisibility(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		card FishCard
		want DiscardVisibility
	}{
		{
			name: "full visibility returns full",
			card: FishCard{DiscardVisibility: DiscardVisibilityFull},
			want: DiscardVisibilityFull,
		},
		{
			name: "masked visibility returns masked",
			card: FishCard{DiscardVisibility: DiscardVisibilityMasked},
			want: DiscardVisibilityMasked,
		},
		{
			name: "hidden visibility returns hidden",
			card: FishCard{DiscardVisibility: DiscardVisibilityHidden},
			want: DiscardVisibilityHidden,
		},
		{
			name: "move only visibility returns move only",
			card: FishCard{DiscardVisibility: DiscardVisibilityMoveOnly},
			want: DiscardVisibilityMoveOnly,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.card.EffectiveDiscardVisibility()
			assert.Equal(t, tt.want, got)
		})
	}
}
