package match

import (
	"pesca/internal/cards"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/player/loadout"
	"pesca/internal/player/rod"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStatusSnapshot(t *testing.T) {
	state := sampleState(t)

	snapshot := NewStatusSnapshot(state)

	assert.Equal(t, 3, snapshot.RoundNumber)
	assert.Equal(t, 2, snapshot.Track.Distance)
	assert.Equal(t, 1, snapshot.Track.Depth)
	assert.Equal(t, 5, snapshot.Track.MaxDistance)
	assert.Equal(t, 4, snapshot.Track.MaxDepth)
	assert.Equal(t, 2, snapshot.Stats.PlayerWins)
	require.Len(t, snapshot.Player.Moves, 1)
	require.Len(t, snapshot.FishDiscard.CurrentCycle.Entries, 1)

	state.Player.Moves.Moves[0].RemainingUses = 0
	state.Deck.CurrentCycle.Entries[0].Name = "Mutado"

	assert.Equal(t, 2, snapshot.Player.Moves[0].RemainingUses)
	assert.True(t, snapshot.Player.Moves[0].HasTopCard)
	assert.Equal(t, "Oleaje abierto", snapshot.FishDiscard.CurrentCycle.Entries[0].Name)
}

func TestNewRoundAndSummarySnapshot(t *testing.T) {
	state := sampleState(t)
	state.Encounter.Status = encounter.StatusCaptured
	state.Encounter.EndReason = encounter.EndReasonTrackCapture
	state.Encounter.LastEvent = encounter.Event{Kind: encounter.EventKindSplash, Escaped: false}

	roundSnapshot := NewRoundSnapshot(RoundResult{
		Round:      2,
		PlayerMove: domain.Blue,
		FishMove:   domain.Red,
		Outcome:    domain.PlayerWin,
		ResolvedEffects: []ResolvedEffectState{{
			Owner:    cards.OwnerFish,
			Trigger:  cards.TriggerOnDraw,
			Type:     cards.EffectTypeLegacyCaptureWindow,
			Priority: 60,
		}},
		Trace: NewResolutionTraceSnapshot(sampleState(t), state, []ResolvedEffectState{{
			Owner:    cards.OwnerFish,
			Trigger:  cards.TriggerOnDraw,
			Type:     cards.EffectTypeLegacyCaptureWindow,
			Priority: 60,
		}}),
		Status:    NewStatusSnapshot(state),
		Encounter: EncounterEventSnapshot{LastEvent: state.Encounter.LastEvent},
	})
	summarySnapshot := NewSummarySnapshot(state)

	assert.Equal(t, domain.Blue, roundSnapshot.PlayerMove)
	assert.Equal(t, domain.Red, roundSnapshot.FishMove)
	assert.Equal(t, encounter.EventKindSplash, roundSnapshot.Encounter.LastEvent.Kind)
	require.Len(t, roundSnapshot.ResolvedEffects, 1)
	assert.Equal(t, cards.OwnerFish, roundSnapshot.ResolvedEffects[0].Owner)
	assert.Equal(t, 2, roundSnapshot.Trace.After.Track.Distance)
	assert.Equal(t, encounter.StatusCaptured, summarySnapshot.Encounter.Status)
	assert.Equal(t, encounter.EndReasonTrackCapture, summarySnapshot.Encounter.EndReason)
	assert.Equal(t, 2, summarySnapshot.Stats.PlayerWins)
}

func sampleState(t *testing.T) State {
	t.Helper()

	encounterState, err := encounter.NewState(encounter.DefaultConfig())
	require.NoError(t, err)
	encounterState.Distance = 2
	encounterState.Depth = 1

	playerLoadout, err := loadout.NewState(rod.State{OpeningMaxDistance: 4, OpeningMaxDepth: 2, TrackMaxDistance: 5, TrackMaxDepth: 4}, nil)
	require.NoError(t, err)

	return State{
		Round:     RoundState{Number: 2},
		Encounter: encounterState,
		Deck: DeckState{
			CurrentCycle: FishDiscardCycleState{Entries: []FishDiscardEntryState{{Name: "Oleaje abierto"}}},
		},
		Player: PlayerState{
			Loadout: playerLoadout,
			Moves: PlayerMoveResources{Moves: []PlayerMoveState{{
				Move:          domain.Blue,
				RemainingUses: 2,
				ActiveCards:   []cards.PlayerCard{cards.NewPlayerCard(domain.Blue)},
			}}},
		},
		Lifecycle: LifecycleState{Stats: Stats{PlayerWins: 2, FishWins: 1, Draws: 3}},
	}
}

func TestNewEncounterEventSnapshot(t *testing.T) {
	t.Run("creates snapshot from encounter state", func(t *testing.T) {
		encounterState, err := encounter.NewState(encounter.DefaultConfig())
		require.NoError(t, err)
		encounterState.LastEvent = encounter.Event{Kind: encounter.EventKindSplash, Escaped: false}

		snapshot := NewEncounterEventSnapshot(encounterState)

		assert.Equal(t, encounter.EventKindSplash, snapshot.LastEvent.Kind)
		assert.False(t, snapshot.LastEvent.Escaped)
	})

	t.Run("creates snapshot with splash state", func(t *testing.T) {
		encounterState, err := encounter.NewState(encounter.DefaultConfig())
		require.NoError(t, err)
		encounterState.LastEvent = encounter.Event{Kind: encounter.EventKindSplash, Escaped: false}
		encounterState.Splash = &encounter.SplashState{TotalJumps: 3, ResolvedJumps: 1, TimeLimit: 5000}

		snapshot := NewEncounterEventSnapshot(encounterState)

		require.NotNil(t, snapshot.Splash)
		assert.Equal(t, 3, snapshot.Splash.TotalJumps)
		assert.Equal(t, 1, snapshot.Splash.ResolvedJumps)
		assert.Equal(t, 2, snapshot.Splash.CurrentJump)
	})

	t.Run("returns nil splash when state has no splash", func(t *testing.T) {
		encounterState, err := encounter.NewState(encounter.DefaultConfig())
		require.NoError(t, err)

		snapshot := NewEncounterEventSnapshot(encounterState)

		assert.Nil(t, snapshot.Splash)
	})
}
