package progression

import (
	"pesca/internal/cards"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/match"
	"pesca/internal/player/loadout"
	"pesca/internal/player/rod"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrackPolicyApply(t *testing.T) {
	tests := []struct {
		title          string
		initialState   func(*testing.T) match.State
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
			title:        "moves the fish closer when the player wins",
			initialState: newMatchState,
			policy:       TrackPolicy{},
			round: match.ResolvedRound{
				PlayerMove:     domain.Blue,
				FishCard:       cards.NewFishCard(domain.Red),
				OutcomeEffects: nil,
				Outcome:        domain.PlayerWin,
			},
			wantDistance:   2,
			wantDepth:      1,
			wantPlayerWins: 1,
		},
		{
			title:        "moves the fish away when the fish wins",
			initialState: newMatchState,
			policy:       TrackPolicy{},
			round: match.ResolvedRound{
				PlayerMove:     domain.Blue,
				FishCard:       cards.NewFishCard(domain.Yellow),
				OutcomeEffects: nil,
				Outcome:        domain.FishWin,
			},
			wantDistance: 4,
			wantDepth:    1,
			wantFishWins: 1,
		},
		{
			title:        "applies card effects after the base round progression",
			initialState: newMatchState,
			policy:       TrackPolicy{},
			round: match.ResolvedRound{
				PlayerMove: domain.Blue,
				FishCard:   cards.NewFishCard(domain.Yellow),
				OutcomeEffects: []cards.CardEffect{{
					Trigger:    cards.TriggerOnOwnerWin,
					DepthShift: 1,
				}},
				Outcome: domain.FishWin,
			},
			wantDistance: 4,
			wantDepth:    2,
			wantFishWins: 1,
		},
		{
			title: "raises the fish toward the surface when the player wins at capture distance",
			initialState: func(t *testing.T) match.State {
				state := newMatchState(t)
				state.Encounter.Distance = state.Encounter.Config.CaptureDistance
				state.Encounter.Depth = 2
				return state
			},
			policy: TrackPolicy{},
			round: match.ResolvedRound{
				PlayerMove:     domain.Blue,
				FishCard:       cards.NewFishCard(domain.Red),
				OutcomeEffects: nil,
				Outcome:        domain.PlayerWin,
			},
			wantDistance:   0,
			wantDepth:      1,
			wantPlayerWins: 1,
		},
		{
			title:        "triggers a splash event when a card raises the fish above the surface",
			initialState: newMatchState,
			policy: TrackPolicy{SplashEscapeDecider: SplashEscapeDeciderFunc(func(float64) bool {
				return false
			})},
			round: match.ResolvedRound{
				PlayerMove: domain.Blue,
				FishCard:   cards.NewFishCard(domain.Red),
				OutcomeEffects: []cards.CardEffect{{
					Trigger:    cards.TriggerOnOwnerLose,
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
			title:        "marks a splash escape when the decider resolves a slip",
			initialState: newMatchState,
			policy: TrackPolicy{SplashEscapeDecider: SplashEscapeDeciderFunc(func(float64) bool {
				return true
			})},
			round: match.ResolvedRound{
				PlayerMove: domain.Blue,
				FishCard:   cards.NewFishCard(domain.Red),
				OutcomeEffects: []cards.CardEffect{{
					Trigger:    cards.TriggerOnOwnerLose,
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
			state := test.initialState(t)
			progressionState := state.ProgressionState()

			test.policy.Apply(&progressionState, test.round)

			assert.Equal(t, test.wantDistance, state.Encounter.Distance)
			assert.Equal(t, test.wantDepth, state.Encounter.Depth)
			assert.Equal(t, test.wantPlayerWins, state.Lifecycle.Stats.PlayerWins)
			assert.Equal(t, test.wantFishWins, state.Lifecycle.Stats.FishWins)
			assert.Equal(t, test.wantDraws, state.Lifecycle.Stats.Draws)
			assert.Equal(t, test.wantEvent, state.Encounter.LastEvent)
			assert.Equal(t, test.wantEndReason, state.Encounter.EndReason)
		})
	}
}

func TestTrackPolicyApplyUsesRoundThresholdBonuses(t *testing.T) {
	state := newMatchState(t)
	state.Encounter.Distance = 1
	state.Encounter.Depth = 2
	state.Round.Thresholds.CaptureDistanceBonus = 1
	progressionState := state.ProgressionState()

	TrackPolicy{}.Apply(&progressionState, match.ResolvedRound{
		PlayerMove:     domain.Blue,
		FishCard:       cards.NewFishCard(domain.Red),
		OutcomeEffects: nil,
		Outcome:        domain.PlayerWin,
	})

	assert.Equal(t, 1, state.Encounter.Distance)
	assert.Equal(t, 1, state.Encounter.Depth)
	assert.Equal(t, 1, state.Lifecycle.Stats.PlayerWins)
}

func TestAccumulateCardEffects(t *testing.T) {
	tests := []struct {
		title        string
		initialDelta encounter.Delta
		effects      []cards.CardEffect
		wantDelta    encounter.Delta
	}{
		{
			title:        "returns the original delta when there are no card effects",
			initialDelta: encounter.Delta{DistanceShift: -1},
			wantDelta:    encounter.Delta{DistanceShift: -1},
		},
		{
			title:        "adds all card effect shifts into the encounter delta",
			initialDelta: encounter.Delta{DistanceShift: -1},
			effects: []cards.CardEffect{{
				DistanceShift: 2,
			}, {
				DepthShift: -1,
			}},
			wantDelta: encounter.Delta{DistanceShift: 1, DepthShift: -1},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			assert.Equal(t, test.wantDelta, accumulateCardEffects(test.initialDelta, test.effects))
		})
	}
}

func newMatchState(t *testing.T) match.State {
	t.Helper()

	encounterState, err := encounter.NewState(encounter.DefaultConfig())
	require.NoError(t, err)
	playerRodState, err := rod.NewState(rod.DefaultConfig())
	require.NoError(t, err)
	playerLoadout, err := loadout.NewState(playerRodState, nil)
	require.NoError(t, err)

	return match.State{
		Encounter: encounterState,
		Player:    match.PlayerState{Loadout: playerLoadout},
	}
}
