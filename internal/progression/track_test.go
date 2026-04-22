package progression

import (
	"pesca/internal/cards"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/match"
	"pesca/internal/playerrig"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrackPolicyApply(t *testing.T) {
	tests := []struct {
		title          string
		policy         TrackPolicy
		round          match.ResolvedRound
		wantDistance   int
		wantDepth      int
		wantPlayerWins int
		wantFishWins   int
		wantDraws      int
		wantEvent      encounter.Event
		wantEndReason  encounter.EndReason
	}{
		{
			title:  "moves the fish closer when the player wins",
			policy: TrackPolicy{},
			round: match.ResolvedRound{
				PlayerMove:         domain.Blue,
				FishCard:           cards.NewFishCard(domain.Red),
				EncounterModifiers: nil,
				Outcome:            domain.PlayerWin,
			},
			wantDistance:   2,
			wantDepth:      1,
			wantPlayerWins: 1,
		},
		{
			title:  "moves the fish away when the fish wins",
			policy: TrackPolicy{},
			round: match.ResolvedRound{
				PlayerMove:         domain.Blue,
				FishCard:           cards.NewFishCard(domain.Yellow),
				EncounterModifiers: nil,
				Outcome:            domain.FishWin,
			},
			wantDistance: 4,
			wantDepth:    1,
			wantFishWins: 1,
		},
		{
			title:  "applies a depth modifier from the fish card when the trigger matches",
			policy: TrackPolicy{},
			round: match.ResolvedRound{
				PlayerMove: domain.Blue,
				FishCard: cards.NewFishCard(domain.Yellow, cards.EncounterModifier{
					Trigger:    cards.TriggerOnFishWin,
					DepthShift: 1,
				}),
				EncounterModifiers: []cards.EncounterModifier{{
					Trigger:    cards.TriggerOnFishWin,
					DepthShift: 1,
				}},
				Outcome: domain.FishWin,
			},
			wantDistance: 4,
			wantDepth:    2,
			wantFishWins: 1,
		},
		{
			title: "triggers a splash event when a card raises the fish above the surface",
			policy: TrackPolicy{SplashEscapeDecider: SplashEscapeDeciderFunc(func(float64) bool {
				return false
			})},
			round: match.ResolvedRound{
				PlayerMove: domain.Blue,
				FishCard: cards.NewFishCard(domain.Red, cards.EncounterModifier{
					Trigger:    cards.TriggerOnPlayerWin,
					DepthShift: -2,
				}),
				EncounterModifiers: []cards.EncounterModifier{{
					Trigger:    cards.TriggerOnPlayerWin,
					DepthShift: -2,
				}},
				Outcome: domain.PlayerWin,
			},
			wantDistance:   2,
			wantDepth:      0,
			wantPlayerWins: 1,
			wantEvent: encounter.Event{
				Kind:    encounter.EventKindSplash,
				Escaped: false,
			},
		},
		{
			title: "marks a splash escape when the decider resolves a slip",
			policy: TrackPolicy{SplashEscapeDecider: SplashEscapeDeciderFunc(func(float64) bool {
				return true
			})},
			round: match.ResolvedRound{
				PlayerMove: domain.Blue,
				FishCard: cards.NewFishCard(domain.Red, cards.EncounterModifier{
					Trigger:    cards.TriggerOnPlayerWin,
					DepthShift: -2,
				}),
				EncounterModifiers: []cards.EncounterModifier{{
					Trigger:    cards.TriggerOnPlayerWin,
					DepthShift: -2,
				}},
				Outcome: domain.PlayerWin,
			},
			wantDistance:   2,
			wantDepth:      0,
			wantPlayerWins: 1,
			wantEvent: encounter.Event{
				Kind:    encounter.EventKindSplash,
				Escaped: true,
			},
			wantEndReason: encounter.EndReasonSplashEscape,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			state := newMatchState(t)

			test.policy.Apply(&state, test.round)

			assert.Equal(t, test.wantDistance, state.Encounter.Distance)
			assert.Equal(t, test.wantDepth, state.Encounter.Depth)
			assert.Equal(t, test.wantPlayerWins, state.Stats.PlayerWins)
			assert.Equal(t, test.wantFishWins, state.Stats.FishWins)
			assert.Equal(t, test.wantDraws, state.Stats.Draws)
			assert.Equal(t, test.wantEvent, state.Encounter.LastEvent)
			assert.Equal(t, test.wantEndReason, state.Encounter.EndReason)
		})
	}
}

func newMatchState(t *testing.T) match.State {
	t.Helper()

	encounterState, err := encounter.NewState(encounter.DefaultConfig())
	require.NoError(t, err)
	playerRigState, err := playerrig.NewState(playerrig.DefaultConfig())
	require.NoError(t, err)

	return match.State{
		Encounter: encounterState,
		PlayerRig: playerRigState,
	}
}
