package game_test

import (
	"pesca/internal/cards"
	"pesca/internal/deck"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/endings"
	"pesca/internal/game"
	"pesca/internal/match"
	"pesca/internal/player/loadout"
	"pesca/internal/player/playermoves"
	"pesca/internal/player/rod"
	"pesca/internal/progression"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewEngine(t *testing.T) {
	tests := []struct {
		title             string
		fishDeck          game.FishDeck
		playerMoves       game.PlayerMoveController
		roundEvaluator    game.RoundEvaluator
		progressionPolicy game.MatchProgressionPolicy
		endCondition      game.MatchEndCondition
		wantErr           string
	}{
		{
			title:             "returns an error when fish deck is missing",
			playerMoves:       &mockPlayerMoveController{},
			roundEvaluator:    &mockRoundEvaluator{},
			progressionPolicy: &mockProgressionPolicy{},
			endCondition:      &mockEndCondition{},
			wantErr:           "fish deck is required",
		},
		{
			title:             "returns an error when player move controller is missing",
			fishDeck:          &mockFishDeck{},
			roundEvaluator:    &mockRoundEvaluator{},
			progressionPolicy: &mockProgressionPolicy{},
			endCondition:      &mockEndCondition{},
			wantErr:           "player moves are required",
		},
		{
			title:             "returns an error when round evaluator is missing",
			fishDeck:          &mockFishDeck{},
			playerMoves:       &mockPlayerMoveController{},
			progressionPolicy: &mockProgressionPolicy{},
			endCondition:      &mockEndCondition{},
			wantErr:           "round evaluator is required",
		},
		{
			title:          "returns an error when progression policy is missing",
			fishDeck:       &mockFishDeck{},
			playerMoves:    &mockPlayerMoveController{},
			roundEvaluator: &mockRoundEvaluator{},
			endCondition:   &mockEndCondition{},
			wantErr:        "progression policy is required",
		},
		{
			title:             "returns an error when end condition is missing",
			fishDeck:          &mockFishDeck{},
			playerMoves:       &mockPlayerMoveController{},
			roundEvaluator:    &mockRoundEvaluator{},
			progressionPolicy: &mockProgressionPolicy{},
			wantErr:           "end condition is required",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			engine, err := game.NewEngine(
				test.fishDeck,
				test.playerMoves,
				test.roundEvaluator,
				test.progressionPolicy,
				test.endCondition,
				match.State{},
			)

			require.Error(t, err)
			assert.EqualError(t, err, test.wantErr)
			assert.Nil(t, engine)
		})
	}

	t.Run("initializes player moves prepares the deck and refreshes the initial state", func(t *testing.T) {
		fixture := newEngineFixture(t, match.State{Player: match.PlayerState{Loadout: samplePlayerLoadoutState(t)}})

		state := fixture.engine.State()

		assert.Equal(t, 7, state.Deck.ActiveCards)
		assert.Equal(t, 2, state.Deck.DiscardCards)
		assert.Equal(t, 1, state.Deck.RecycleCount)
		assert.True(t, state.Deck.Exhausted)
		assert.True(t, state.Deck.ShufflesOnRecycle)
		assert.Equal(t, 3, state.Deck.CardsToRemove)
		assert.Equal(t, 2, state.Deck.CurrentCycle.Number)
		assert.Equal(t, 2, state.Deck.CurrentCycle.TotalCards)
		require.Len(t, state.Deck.CurrentCycle.Entries, 2)
		assert.Equal(t, cards.DiscardVisibilityFull, state.Deck.CurrentCycle.Entries[0].Visibility)
		require.Len(t, state.Deck.PreviousCycleStats, 1)
		assert.Equal(t, 1, state.Deck.PreviousCycleStats[0].Number)
		assert.Equal(t, samplePlayerLoadoutState(t), state.Player.Loadout)
		assert.Equal(t, samplePlayerMoveResources(), state.Player.Moves)
		fixture.assertExpectations(t)
	})
}

func TestEnginePlayRound(t *testing.T) {
	t.Run("returns an error when the game is already finished", func(t *testing.T) {
		fixture := newEngineFixture(t, match.State{Lifecycle: match.LifecycleState{Finished: true}, Player: match.PlayerState{Loadout: samplePlayerLoadoutState(t)}})

		result, err := fixture.engine.PlayRound(domain.Blue)

		assert.ErrorIs(t, err, game.ErrGameFinished)
		assert.Equal(t, match.RoundResult{}, result)
		fixture.assertExpectations(t)
	})

	t.Run("returns a validation error when the player move controller rejects the move", func(t *testing.T) {
		fixture := newEngineFixture(t, match.State{Player: match.PlayerState{Loadout: samplePlayerLoadoutState(t)}})
		prepareRoundCall := fixture.playerMoveController.On("PrepareRound", mock.AnythingOfType("*match.PlayerMoveRuntime")).Once()
		validateMoveCall := fixture.playerMoveController.On("ValidateMove", mock.Anything, domain.Blue).Return(playermoves.ErrMoveUnavailable).Once()
		mock.InOrder(prepareRoundCall, validateMoveCall)

		result, err := fixture.engine.PlayRound(domain.Blue)

		assert.ErrorIs(t, err, playermoves.ErrMoveUnavailable)
		assert.Equal(t, match.RoundResult{}, result)
		fixture.assertExpectations(t)
	})

	t.Run("returns the deck error when drawing a fish move fails and the match remains open", func(t *testing.T) {
		fixture := newEngineFixture(t, match.State{Player: match.PlayerState{Loadout: samplePlayerLoadoutState(t)}})
		prepareRoundCall := fixture.playerMoveController.On("PrepareRound", mock.AnythingOfType("*match.PlayerMoveRuntime")).Once()
		validateMoveCall := fixture.playerMoveController.On("ValidateMove", mock.Anything, domain.Blue).Return(nil).Once()
		peekCardCall := fixture.playerMoveController.On("PeekMoveCard", mock.Anything, domain.Blue).Return(cards.NewPlayerCard(domain.Blue), nil).Once()
		drawCall := fixture.fishDeck.On("Draw").Return(cards.FishCard{}, deck.ErrNoCardsAvailable).Once()
		activeCountCall := fixture.fishDeck.On("ActiveCount").Return(4).Once()
		discardCountCall := fixture.fishDeck.On("DiscardCount").Return(5).Once()
		recycleCountCall := fixture.fishDeck.On("RecycleCount").Return(2).Once()
		exhaustedCall := fixture.fishDeck.On("Exhausted").Return(false).Once()
		endConditionCall := fixture.endCondition.On("Apply", mock.AnythingOfType("*match.EndingState")).Once()
		mock.InOrder(prepareRoundCall, validateMoveCall, peekCardCall, drawCall, activeCountCall, discardCountCall, recycleCountCall, exhaustedCall, endConditionCall)

		result, err := fixture.engine.PlayRound(domain.Blue)

		assert.ErrorIs(t, err, deck.ErrNoCardsAvailable)
		assert.Equal(t, match.RoundResult{}, result)
		assert.False(t, fixture.engine.State().Lifecycle.Finished)
		fixture.assertExpectations(t)
	})

	t.Run("returns game finished when draw fails and the end condition closes the match", func(t *testing.T) {
		fixture := newEngineFixture(t, match.State{Player: match.PlayerState{Loadout: samplePlayerLoadoutState(t)}})
		prepareRoundCall := fixture.playerMoveController.On("PrepareRound", mock.AnythingOfType("*match.PlayerMoveRuntime")).Once()
		validateMoveCall := fixture.playerMoveController.On("ValidateMove", mock.Anything, domain.Blue).Return(nil).Once()
		peekCardCall := fixture.playerMoveController.On("PeekMoveCard", mock.Anything, domain.Blue).Return(cards.NewPlayerCard(domain.Blue), nil).Once()
		drawCall := fixture.fishDeck.On("Draw").Return(cards.FishCard{}, deck.ErrNoCardsAvailable).Once()
		activeCountCall := fixture.fishDeck.On("ActiveCount").Return(0).Once()
		discardCountCall := fixture.fishDeck.On("DiscardCount").Return(0).Once()
		recycleCountCall := fixture.fishDeck.On("RecycleCount").Return(3).Once()
		exhaustedCall := fixture.fishDeck.On("Exhausted").Return(true).Once()
		endConditionCall := fixture.endCondition.On("Apply", mock.AnythingOfType("*match.EndingState")).Run(func(arguments mock.Arguments) {
			state := arguments.Get(0).(*match.EndingState)
			state.Lifecycle.Finished = true
		}).Once()
		mock.InOrder(prepareRoundCall, validateMoveCall, peekCardCall, drawCall, activeCountCall, discardCountCall, recycleCountCall, exhaustedCall, endConditionCall)

		result, err := fixture.engine.PlayRound(domain.Blue)

		assert.ErrorIs(t, err, game.ErrGameFinished)
		assert.Equal(t, match.RoundResult{}, result)
		assert.True(t, fixture.engine.State().Lifecycle.Finished)
		fixture.assertExpectations(t)
	})

	t.Run("orchestrates the round using the injected deck and player move controller", func(t *testing.T) {
		fixture := newEngineFixture(t, match.State{Player: match.PlayerState{Loadout: samplePlayerLoadoutState(t)}})
		playerCard := cards.NewPlayerCard(domain.Blue, cards.CardEffect{Trigger: cards.TriggerOnDraw, SurfaceDepthBonus: 1})
		prepareBeforeValidationCall := fixture.playerMoveController.On("PrepareRound", mock.AnythingOfType("*match.PlayerMoveRuntime")).Run(func(arguments mock.Arguments) {
			state := arguments.Get(0).(*match.PlayerMoveRuntime)
			assert.Equal(t, 0, state.Round.Number)
		}).Once()
		validateMoveCall := fixture.playerMoveController.On("ValidateMove", mock.Anything, domain.Blue).Run(func(arguments mock.Arguments) {
			state := arguments.Get(0).(match.PlayerMoveRuntime)
			assert.Equal(t, 0, state.Round.Number)
		}).Return(nil).Once()
		peekCardCall := fixture.playerMoveController.On("PeekMoveCard", mock.Anything, domain.Blue).Return(playerCard, nil).Once()
		fishCard := cards.NewFishCard(domain.Red,
			cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1},
			cards.CardEffect{Trigger: cards.TriggerOnOwnerLose, DepthShift: -1},
			cards.CardEffect{Trigger: cards.TriggerOnOwnerWin, DistanceShift: 2},
		)
		drawCall := fixture.fishDeck.On("Draw").Return(fishCard, nil).Once()
		evaluateCall := fixture.roundEvaluator.On("Evaluate", domain.Blue, domain.Red).Return(domain.PlayerWin).Once()
		consumeMoveCall := fixture.playerMoveController.On("ConsumeMove", mock.AnythingOfType("*match.PlayerMoveRuntime"), domain.Blue).Run(func(arguments mock.Arguments) {
			state := arguments.Get(0).(*match.PlayerMoveRuntime)
			assert.Equal(t, 1, state.Round.Number)
		}).Return(playerCard).Once()
		progressionCall := fixture.progressionPolicy.On("Apply", mock.AnythingOfType("*match.ProgressionState"), match.ResolvedRound{
			PlayerMove: domain.Blue,
			PlayerCard: playerCard,
			FishCard:   fishCard,
			DrawEffects: []cards.CardEffect{{
				Trigger:           cards.TriggerOnDraw,
				SurfaceDepthBonus: 1,
			}, {
				Trigger:              cards.TriggerOnDraw,
				CaptureDistanceBonus: 1,
			}},
			OutcomeEffects: []cards.CardEffect{{
				Trigger:    cards.TriggerOnOwnerLose,
				DepthShift: -1,
			}},
			Outcome: domain.PlayerWin,
		}).Run(func(arguments mock.Arguments) {
			state := arguments.Get(0).(*match.ProgressionState)
			assert.Equal(t, 1, state.Round.Thresholds.CaptureDistanceBonus)
			assert.Equal(t, 1, state.Round.Thresholds.SurfaceDepthBonus)
			state.Lifecycle.Stats.PlayerWins++
		}).Once()
		prepareDeckCall := fixture.fishDeck.On("PrepareNextRound").Once()
		activeCountCall := fixture.fishDeck.On("ActiveCount").Return(6).Once()
		discardCountCall := fixture.fishDeck.On("DiscardCount").Return(3).Once()
		recycleCountCall := fixture.fishDeck.On("RecycleCount").Return(1).Once()
		exhaustedCall := fixture.fishDeck.On("Exhausted").Return(false).Once()
		prepareAfterRoundCall := fixture.playerMoveController.On("PrepareRound", mock.AnythingOfType("*match.PlayerMoveRuntime")).Run(func(arguments mock.Arguments) {
			state := arguments.Get(0).(*match.PlayerMoveRuntime)
			assert.Equal(t, 1, state.Round.Number)
		}).Once()
		endConditionCall := fixture.endCondition.On("Apply", mock.AnythingOfType("*match.EndingState")).Run(func(arguments mock.Arguments) {
			state := arguments.Get(0).(*match.EndingState)
			assert.Equal(t, 1, state.Round.Thresholds.CaptureDistanceBonus)
			assert.Equal(t, 1, state.Round.Thresholds.SurfaceDepthBonus)
		}).Once()

		mock.InOrder(
			prepareBeforeValidationCall,
			validateMoveCall,
			peekCardCall,
			drawCall,
			evaluateCall,
			consumeMoveCall,
			progressionCall,
			prepareDeckCall,
			activeCountCall,
			discardCountCall,
			recycleCountCall,
			exhaustedCall,
			prepareAfterRoundCall,
			endConditionCall,
		)

		result, err := fixture.engine.PlayRound(domain.Blue)

		require.NoError(t, err)
		assert.Equal(t, 1, result.Round)
		assert.Equal(t, domain.Blue, result.PlayerMove)
		assert.Equal(t, playerCard, result.PlayerCard)
		assert.Equal(t, domain.Red, result.FishMove)
		assert.Equal(t, fishCard, result.FishCard)
		assert.Equal(t, domain.PlayerWin, result.Outcome)
		assert.Equal(t, 6, result.Status.FishDiscard.ActiveCards)
		assert.Equal(t, 3, result.Status.FishDiscard.DiscardCards)
		assert.Equal(t, 1, result.Status.Stats.PlayerWins)
		assert.Equal(t, 2, result.Status.RoundNumber)
		assert.Equal(t, match.RoundState{Number: 1}, fixture.engine.State().Round)
		fixture.assertExpectations(t)
	})

	t.Run("captures using a real on draw threshold effect before end evaluation", func(t *testing.T) {
		encounterState, err := encounter.NewState(encounter.DefaultConfig())
		require.NoError(t, err)
		encounterState.Distance = 1
		encounterState.Depth = 0

		playerMoveController := &mockPlayerMoveController{}
		roundEvaluator := &mockRoundEvaluator{}
		fishDeck := &mockFishDeck{}
		fishDeck.visibilitySnapshots = []deck.VisibilitySnapshot{{
			CurrentCycle: deck.VisibleDiscardCycle{Number: 1},
		}}
		playerCard := cards.NewPlayerCard(domain.Blue)

		initializeCall := playerMoveController.On("Initialize", mock.AnythingOfType("*match.PlayerMoveRuntime")).Run(func(arguments mock.Arguments) {
			state := arguments.Get(0).(*match.PlayerMoveRuntime)
			state.Moves.Moves = samplePlayerMoveResources().Moves
		}).Once()
		prepareDeckCall := fishDeck.On("PrepareNextRound").Once()
		activeCountInitCall := fishDeck.On("ActiveCount").Return(0).Once()
		discardCountInitCall := fishDeck.On("DiscardCount").Return(0).Once()
		recycleCountInitCall := fishDeck.On("RecycleCount").Return(0).Once()
		exhaustedInitCall := fishDeck.On("Exhausted").Return(false).Once()

		prepareBeforeValidationCall := playerMoveController.On("PrepareRound", mock.AnythingOfType("*match.PlayerMoveRuntime")).Once()
		validateMoveCall := playerMoveController.On("ValidateMove", mock.Anything, domain.Blue).Return(nil).Once()
		peekCardCall := playerMoveController.On("PeekMoveCard", mock.Anything, domain.Blue).Return(playerCard, nil).Once()
		fishCard := cards.NewFishCard(domain.Red, cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1})
		drawCall := fishDeck.On("Draw").Return(fishCard, nil).Once()
		evaluateCall := roundEvaluator.On("Evaluate", domain.Blue, domain.Red).Return(domain.Draw).Once()
		consumeMoveCall := playerMoveController.On("ConsumeMove", mock.AnythingOfType("*match.PlayerMoveRuntime"), domain.Blue).Return(playerCard).Once()
		prepareDeckAfterRoundCall := fishDeck.On("PrepareNextRound").Once()
		activeCountAfterRoundCall := fishDeck.On("ActiveCount").Return(0).Once()
		discardCountAfterRoundCall := fishDeck.On("DiscardCount").Return(1).Once()
		recycleCountAfterRoundCall := fishDeck.On("RecycleCount").Return(0).Once()
		exhaustedAfterRoundCall := fishDeck.On("Exhausted").Return(false).Once()
		prepareAfterRoundCall := playerMoveController.On("PrepareRound", mock.AnythingOfType("*match.PlayerMoveRuntime")).Once()

		mock.InOrder(
			initializeCall,
			prepareDeckCall,
			activeCountInitCall,
			discardCountInitCall,
			recycleCountInitCall,
			exhaustedInitCall,
			prepareBeforeValidationCall,
			validateMoveCall,
			peekCardCall,
			drawCall,
			evaluateCall,
			consumeMoveCall,
			prepareDeckAfterRoundCall,
			activeCountAfterRoundCall,
			discardCountAfterRoundCall,
			recycleCountAfterRoundCall,
			exhaustedAfterRoundCall,
			prepareAfterRoundCall,
		)

		engine, err := game.NewEngine(
			fishDeck,
			playerMoveController,
			roundEvaluator,
			progression.TrackPolicy{},
			endings.EncounterEndCondition{},
			match.State{Encounter: encounterState, Player: match.PlayerState{Loadout: samplePlayerLoadoutState(t)}},
		)
		require.NoError(t, err)

		result, err := engine.PlayRound(domain.Blue)

		require.NoError(t, err)
		assert.Equal(t, engine.State().Encounter.LastEvent, result.Encounter.LastEvent)
		assert.Equal(t, 2, result.Status.RoundNumber)
		assert.True(t, engine.State().Lifecycle.Finished)
		assert.Equal(t, encounter.StatusCaptured, engine.State().Encounter.Status)
		assert.Equal(t, encounter.EndReasonTrackCapture, engine.State().Encounter.EndReason)
		assert.Equal(t, match.RoundState{Number: 1}, engine.State().Round)
		fishDeck.AssertExpectations(t)
		playerMoveController.AssertExpectations(t)
		roundEvaluator.AssertExpectations(t)
	})
}

type engineFixture struct {
	engine               *game.Engine
	fishDeck             *mockFishDeck
	playerMoveController *mockPlayerMoveController
	roundEvaluator       *mockRoundEvaluator
	progressionPolicy    *mockProgressionPolicy
	endCondition         *mockEndCondition
}

func newEngineFixture(t *testing.T, initialState match.State) engineFixture {
	t.Helper()
	if initialState.Player.Loadout.TrackMaxDistance() == 0 {
		initialState.Player.Loadout = samplePlayerLoadoutState(t)
	}

	fishDeck := &mockFishDeck{}
	fishDeck.visibilitySnapshots = []deck.VisibilitySnapshot{{
		CurrentCycle: deck.VisibleDiscardCycle{
			Number:     2,
			TotalCards: 2,
			Entries: []deck.VisibleDiscardEntry{
				{Visibility: cards.DiscardVisibilityFull, Move: domain.Blue, Name: "Oleaje abierto"},
				{Visibility: cards.DiscardVisibilityMoveOnly, Move: domain.Red},
			},
		},
		PreviousCycles: []deck.VisibleDiscardCycleSummary{{
			Number:       1,
			TotalCards:   4,
			VisibleCards: 3,
			HiddenCards:  1,
		}},
		RecycleCount:      1,
		ShufflesOnRecycle: true,
		CardsToRemove:     3,
		Exhausted:         true,
	}}
	playerMoveController := &mockPlayerMoveController{}
	roundEvaluator := &mockRoundEvaluator{}
	progressionPolicy := &mockProgressionPolicy{}
	endCondition := &mockEndCondition{}

	initializeCall := playerMoveController.On("Initialize", mock.AnythingOfType("*match.PlayerMoveRuntime")).Run(func(arguments mock.Arguments) {
		state := arguments.Get(0).(*match.PlayerMoveRuntime)
		state.Moves.Moves = samplePlayerMoveResources().Moves
	}).Once()
	prepareDeckCall := fishDeck.On("PrepareNextRound").Once()
	activeCountCall := fishDeck.On("ActiveCount").Return(7).Once()
	discardCountCall := fishDeck.On("DiscardCount").Return(2).Once()
	recycleCountCall := fishDeck.On("RecycleCount").Return(1).Once()
	exhaustedCall := fishDeck.On("Exhausted").Return(true).Once()
	endConditionCall := endCondition.On("Apply", mock.AnythingOfType("*match.EndingState")).Once()
	mock.InOrder(initializeCall, prepareDeckCall, activeCountCall, discardCountCall, recycleCountCall, exhaustedCall, endConditionCall)

	engine, err := game.NewEngine(
		fishDeck,
		playerMoveController,
		roundEvaluator,
		progressionPolicy,
		endCondition,
		initialState,
	)
	require.NoError(t, err)

	return engineFixture{
		engine:               engine,
		fishDeck:             fishDeck,
		playerMoveController: playerMoveController,
		roundEvaluator:       roundEvaluator,
		progressionPolicy:    progressionPolicy,
		endCondition:         endCondition,
	}
}

func (fixture engineFixture) assertExpectations(t *testing.T) {
	t.Helper()

	fixture.fishDeck.AssertExpectations(t)
	fixture.playerMoveController.AssertExpectations(t)
	fixture.roundEvaluator.AssertExpectations(t)
	fixture.progressionPolicy.AssertExpectations(t)
	fixture.endCondition.AssertExpectations(t)
}

func samplePlayerMoveResources() match.PlayerMoveResources {
	return match.PlayerMoveResources{Moves: []match.PlayerMoveState{
		{Move: domain.Blue, MaxUses: 3, RemainingUses: 3, ActiveCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Blue), cards.NewPlayerCard(domain.Blue), cards.NewPlayerCard(domain.Blue)}},
		{Move: domain.Red, MaxUses: 3, RemainingUses: 3, ActiveCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Red), cards.NewPlayerCard(domain.Red), cards.NewPlayerCard(domain.Red)}},
		{Move: domain.Yellow, MaxUses: 3, RemainingUses: 3, ActiveCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Yellow), cards.NewPlayerCard(domain.Yellow), cards.NewPlayerCard(domain.Yellow)}},
	}}
}

func samplePlayerRodState() rod.State {
	return rod.State{
		OpeningMaxDistance: 5,
		OpeningMaxDepth:    3,
		TrackMaxDistance:   5,
		TrackMaxDepth:      3,
	}
}

func samplePlayerLoadoutState(t *testing.T) loadout.State {
	t.Helper()

	playerLoadout, err := loadout.NewState(samplePlayerRodState(), nil)
	require.NoError(t, err)

	return playerLoadout
}

type mockFishDeck struct {
	mock.Mock
	visibilitySnapshots     []deck.VisibilitySnapshot
	visibilitySnapshotCalls int
}

func (mockDeck *mockFishDeck) Draw() (cards.FishCard, error) {
	arguments := mockDeck.Called()
	return arguments.Get(0).(cards.FishCard), arguments.Error(1)
}

func (mockDeck *mockFishDeck) PrepareNextRound() {
	mockDeck.Called()
}

func (mockDeck *mockFishDeck) ActiveCount() int {
	return mockDeck.Called().Int(0)
}

func (mockDeck *mockFishDeck) DiscardCount() int {
	return mockDeck.Called().Int(0)
}

func (mockDeck *mockFishDeck) RecycleCount() int {
	return mockDeck.Called().Int(0)
}

func (mockDeck *mockFishDeck) Exhausted() bool {
	return mockDeck.Called().Bool(0)
}

func (mockDeck *mockFishDeck) VisibilitySnapshot() deck.VisibilitySnapshot {
	if len(mockDeck.visibilitySnapshots) == 0 {
		return deck.VisibilitySnapshot{}
	}

	index := mockDeck.visibilitySnapshotCalls
	if index >= len(mockDeck.visibilitySnapshots) {
		index = len(mockDeck.visibilitySnapshots) - 1
	}
	mockDeck.visibilitySnapshotCalls++

	return mockDeck.visibilitySnapshots[index]
}

type mockPlayerMoveController struct {
	mock.Mock
}

func (mockController *mockPlayerMoveController) Initialize(state *match.PlayerMoveRuntime) {
	mockController.Called(state)
}

func (mockController *mockPlayerMoveController) PrepareRound(state *match.PlayerMoveRuntime) {
	mockController.Called(state)
}

func (mockController *mockPlayerMoveController) ValidateMove(state match.PlayerMoveRuntime, playerMove domain.Move) error {
	arguments := mockController.Called(state, playerMove)
	return arguments.Error(0)
}

func (mockController *mockPlayerMoveController) PeekMoveCard(state match.PlayerMoveRuntime, playerMove domain.Move) (cards.PlayerCard, error) {
	arguments := mockController.Called(state, playerMove)
	return arguments.Get(0).(cards.PlayerCard), arguments.Error(1)
}

func (mockController *mockPlayerMoveController) ConsumeMove(state *match.PlayerMoveRuntime, playerMove domain.Move) cards.PlayerCard {
	arguments := mockController.Called(state, playerMove)
	return arguments.Get(0).(cards.PlayerCard)
}

type mockRoundEvaluator struct {
	mock.Mock
}

func (mockEvaluator *mockRoundEvaluator) Evaluate(playerMove, fishMove domain.Move) domain.RoundOutcome {
	arguments := mockEvaluator.Called(playerMove, fishMove)
	return arguments.Get(0).(domain.RoundOutcome)
}

type mockProgressionPolicy struct {
	mock.Mock
}

func (mockPolicy *mockProgressionPolicy) Apply(state *match.ProgressionState, round match.ResolvedRound) {
	mockPolicy.Called(state, round)
}

type mockEndCondition struct {
	mock.Mock
}

func (mockCondition *mockEndCondition) Apply(state *match.EndingState) {
	mockCondition.Called(state)
}

func TestEngineResolveSplash(t *testing.T) {
	t.Run("returns an error when the game is already finished", func(t *testing.T) {
		fixture := newEngineFixture(t, match.State{Lifecycle: match.LifecycleState{Finished: true}, Player: match.PlayerState{Loadout: samplePlayerLoadoutState(t)}})

		err := fixture.engine.ResolveSplash(encounter.SplashResolution{})

		assert.ErrorIs(t, err, game.ErrGameFinished)
	})

	t.Run("applies the resolution and finishes the game when splash escape occurs", func(t *testing.T) {
		fixture := newEngineFixture(t, match.State{
			Player: match.PlayerState{Loadout: samplePlayerLoadoutState(t)},
			Encounter: encounter.State{
				Config:    encounter.DefaultConfig(),
				Splash:    &encounter.SplashState{TotalJumps: 1},
				LastEvent: encounter.Event{Kind: encounter.EventKindSplash},
			},
		})
		fixture.fishDeck.visibilitySnapshots = append(fixture.fishDeck.visibilitySnapshots, deck.VisibilitySnapshot{Exhausted: false, ShufflesOnRecycle: true}, deck.VisibilitySnapshot{Exhausted: false, ShufflesOnRecycle: true})
		fixture.fishDeck.On("ActiveCount").Return(7).Once()
		fixture.fishDeck.On("DiscardCount").Return(2).Once()
		fixture.fishDeck.On("RecycleCount").Return(1).Once()
		fixture.fishDeck.On("Exhausted").Return(false).Once()
		fixture.endCondition.On("Apply", mock.AnythingOfType("*match.EndingState")).Run(func(arguments mock.Arguments) {
			state := arguments.Get(0).(*match.EndingState)
			state.Lifecycle.Finished = true
			state.Encounter.Status = encounter.StatusEscaped
			state.Encounter.EndReason = encounter.EndReasonSplashEscape
		}).Once()

		err := fixture.engine.ResolveSplash(encounter.SplashResolution{Escaped: true})

		assert.NoError(t, err)
		assert.True(t, fixture.engine.State().Lifecycle.Finished)
		assert.Equal(t, encounter.StatusEscaped, fixture.engine.State().Encounter.Status)
		assert.Equal(t, encounter.EndReasonSplashEscape, fixture.engine.State().Encounter.EndReason)
		fixture.assertExpectations(t)
	})

	t.Run("applies the resolution and continues when splash succeeds", func(t *testing.T) {
		fixture := newEngineFixture(t, match.State{
			Player: match.PlayerState{Loadout: samplePlayerLoadoutState(t)},
			Encounter: encounter.State{
				Config:    encounter.DefaultConfig(),
				Splash:    &encounter.SplashState{TotalJumps: 1},
				LastEvent: encounter.Event{Kind: encounter.EventKindSplash},
			},
		})
		fixture.fishDeck.visibilitySnapshots = append(fixture.fishDeck.visibilitySnapshots, deck.VisibilitySnapshot{Exhausted: false, ShufflesOnRecycle: true}, deck.VisibilitySnapshot{Exhausted: false, ShufflesOnRecycle: true})
		fixture.fishDeck.On("ActiveCount").Return(7).Once()
		fixture.fishDeck.On("DiscardCount").Return(2).Once()
		fixture.fishDeck.On("RecycleCount").Return(1).Once()
		fixture.fishDeck.On("Exhausted").Return(false).Once()
		fixture.endCondition.On("Apply", mock.AnythingOfType("*match.EndingState")).Once()

		err := fixture.engine.ResolveSplash(encounter.SplashResolution{SuccessfulJumps: 1})

		assert.NoError(t, err)
		assert.False(t, fixture.engine.State().Lifecycle.Finished)
		assert.Nil(t, fixture.engine.State().Encounter.Splash)
		fixture.assertExpectations(t)
	})
}
