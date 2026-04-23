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
	assert.Equal(t, effects[1], filteredEffects[0])
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
	assert.Equal(t, effects[0], filteredEffects[0])
}

func TestNewFishCard(t *testing.T) {
	effect := CardEffect{Trigger: TriggerOnOwnerWin, DepthShift: 1}
	card := NewFishCard(domain.Red, effect)

	require.Len(t, card.Effects, 1)
	assert.Equal(t, domain.Red, card.Move)
	assert.Equal(t, effect, card.Effects[0])
}

func TestNewPlayerCard(t *testing.T) {
	effect := CardEffect{Trigger: TriggerOnOwnerLose, DistanceShift: 1}
	card := NewPlayerCard(domain.Blue, effect)

	require.Len(t, card.Effects, 1)
	assert.Equal(t, domain.Blue, card.Move)
	assert.Equal(t, effect, card.Effects[0])
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
