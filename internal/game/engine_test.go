package game_test

import (
	"pesca/internal/cards"
	"pesca/internal/deck"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/endings"
	"pesca/internal/game"
	"pesca/internal/match"
	"pesca/internal/player/playermoves"
	"pesca/internal/player/playerrig"
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
		fixture := newEngineFixture(t, match.State{PlayerRig: samplePlayerRigState()})

		state := fixture.engine.State()

		assert.Equal(t, 7, state.Deck.ActiveCards)
		assert.Equal(t, 2, state.Deck.DiscardCards)
		assert.Equal(t, 1, state.Deck.RecycleCount)
		assert.True(t, state.Deck.Exhausted)
		assert.Equal(t, samplePlayerRigState(), state.PlayerRig)
		assert.Equal(t, samplePlayerMoveResources(), state.PlayerMoves)
		fixture.assertExpectations(t)
	})
}

func TestEnginePlayRound(t *testing.T) {
	t.Run("returns an error when the game is already finished", func(t *testing.T) {
		fixture := newEngineFixture(t, match.State{Finished: true, PlayerRig: samplePlayerRigState()})

		result, err := fixture.engine.PlayRound(domain.Blue)

		assert.ErrorIs(t, err, game.ErrGameFinished)
		assert.Equal(t, match.RoundResult{}, result)
		fixture.assertExpectations(t)
	})

	t.Run("returns a validation error when the player move controller rejects the move", func(t *testing.T) {
		fixture := newEngineFixture(t, match.State{PlayerRig: samplePlayerRigState()})
		prepareRoundCall := fixture.playerMoveController.On("PrepareRound", mock.AnythingOfType("*match.State")).Once()
		validateMoveCall := fixture.playerMoveController.On("ValidateMove", mock.Anything, domain.Blue).Return(playermoves.ErrMoveUnavailable).Once()
		mock.InOrder(prepareRoundCall, validateMoveCall)

		result, err := fixture.engine.PlayRound(domain.Blue)

		assert.ErrorIs(t, err, playermoves.ErrMoveUnavailable)
		assert.Equal(t, match.RoundResult{}, result)
		fixture.assertExpectations(t)
	})

	t.Run("returns the deck error when drawing a fish move fails and the match remains open", func(t *testing.T) {
		fixture := newEngineFixture(t, match.State{PlayerRig: samplePlayerRigState()})
		prepareRoundCall := fixture.playerMoveController.On("PrepareRound", mock.AnythingOfType("*match.State")).Once()
		validateMoveCall := fixture.playerMoveController.On("ValidateMove", mock.Anything, domain.Blue).Return(nil).Once()
		peekCardCall := fixture.playerMoveController.On("PeekMoveCard", mock.Anything, domain.Blue).Return(cards.NewPlayerCard(domain.Blue), nil).Once()
		drawCall := fixture.fishDeck.On("Draw").Return(cards.FishCard{}, deck.ErrNoCardsAvailable).Once()
		activeCountCall := fixture.fishDeck.On("ActiveCount").Return(4).Once()
		discardCountCall := fixture.fishDeck.On("DiscardCount").Return(5).Once()
		recycleCountCall := fixture.fishDeck.On("RecycleCount").Return(2).Once()
		exhaustedCall := fixture.fishDeck.On("Exhausted").Return(false).Once()
		endConditionCall := fixture.endCondition.On("Apply", mock.AnythingOfType("*match.State")).Once()
		mock.InOrder(prepareRoundCall, validateMoveCall, peekCardCall, drawCall, activeCountCall, discardCountCall, recycleCountCall, exhaustedCall, endConditionCall)

		result, err := fixture.engine.PlayRound(domain.Blue)

		assert.ErrorIs(t, err, deck.ErrNoCardsAvailable)
		assert.Equal(t, match.RoundResult{}, result)
		assert.False(t, fixture.engine.State().Finished)
		fixture.assertExpectations(t)
	})

	t.Run("returns game finished when draw fails and the end condition closes the match", func(t *testing.T) {
		fixture := newEngineFixture(t, match.State{PlayerRig: samplePlayerRigState()})
		prepareRoundCall := fixture.playerMoveController.On("PrepareRound", mock.AnythingOfType("*match.State")).Once()
		validateMoveCall := fixture.playerMoveController.On("ValidateMove", mock.Anything, domain.Blue).Return(nil).Once()
		peekCardCall := fixture.playerMoveController.On("PeekMoveCard", mock.Anything, domain.Blue).Return(cards.NewPlayerCard(domain.Blue), nil).Once()
		drawCall := fixture.fishDeck.On("Draw").Return(cards.FishCard{}, deck.ErrNoCardsAvailable).Once()
		activeCountCall := fixture.fishDeck.On("ActiveCount").Return(0).Once()
		discardCountCall := fixture.fishDeck.On("DiscardCount").Return(0).Once()
		recycleCountCall := fixture.fishDeck.On("RecycleCount").Return(3).Once()
		exhaustedCall := fixture.fishDeck.On("Exhausted").Return(true).Once()
		endConditionCall := fixture.endCondition.On("Apply", mock.AnythingOfType("*match.State")).Run(func(arguments mock.Arguments) {
			state := arguments.Get(0).(*match.State)
			state.Finished = true
		}).Once()
		mock.InOrder(prepareRoundCall, validateMoveCall, peekCardCall, drawCall, activeCountCall, discardCountCall, recycleCountCall, exhaustedCall, endConditionCall)

		result, err := fixture.engine.PlayRound(domain.Blue)

		assert.ErrorIs(t, err, game.ErrGameFinished)
		assert.Equal(t, match.RoundResult{}, result)
		assert.True(t, fixture.engine.State().Finished)
		fixture.assertExpectations(t)
	})

	t.Run("orchestrates the round using the injected deck and player move controller", func(t *testing.T) {
		fixture := newEngineFixture(t, match.State{PlayerRig: samplePlayerRigState()})
		playerCard := cards.NewPlayerCard(domain.Blue, cards.CardEffect{Trigger: cards.TriggerOnDraw, SurfaceDepthBonus: 1})
		prepareBeforeValidationCall := fixture.playerMoveController.On("PrepareRound", mock.AnythingOfType("*match.State")).Run(func(arguments mock.Arguments) {
			state := arguments.Get(0).(*match.State)
			assert.Equal(t, 0, state.Round)
		}).Once()
		validateMoveCall := fixture.playerMoveController.On("ValidateMove", mock.Anything, domain.Blue).Run(func(arguments mock.Arguments) {
			state := arguments.Get(0).(match.State)
			assert.Equal(t, 0, state.Round)
		}).Return(nil).Once()
		peekCardCall := fixture.playerMoveController.On("PeekMoveCard", mock.Anything, domain.Blue).Return(playerCard, nil).Once()
		fishCard := cards.NewFishCard(domain.Red,
			cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1},
			cards.CardEffect{Trigger: cards.TriggerOnOwnerLose, DepthShift: -1},
			cards.CardEffect{Trigger: cards.TriggerOnOwnerWin, DistanceShift: 2},
		)
		drawCall := fixture.fishDeck.On("Draw").Return(fishCard, nil).Once()
		evaluateCall := fixture.roundEvaluator.On("Evaluate", domain.Blue, domain.Red).Return(domain.PlayerWin).Once()
		consumeMoveCall := fixture.playerMoveController.On("ConsumeMove", mock.AnythingOfType("*match.State"), domain.Blue).Run(func(arguments mock.Arguments) {
			state := arguments.Get(0).(*match.State)
			assert.Equal(t, 1, state.Round)
		}).Return(playerCard).Once()
		progressionCall := fixture.progressionPolicy.On("Apply", mock.AnythingOfType("*match.State"), match.ResolvedRound{
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
			state := arguments.Get(0).(*match.State)
			assert.Equal(t, 1, state.RoundState.Thresholds.CaptureDistanceBonus)
			assert.Equal(t, 1, state.RoundState.Thresholds.SurfaceDepthBonus)
			state.Stats.PlayerWins++
		}).Once()
		prepareDeckCall := fixture.fishDeck.On("PrepareNextRound").Once()
		activeCountCall := fixture.fishDeck.On("ActiveCount").Return(6).Once()
		discardCountCall := fixture.fishDeck.On("DiscardCount").Return(3).Once()
		recycleCountCall := fixture.fishDeck.On("RecycleCount").Return(1).Once()
		exhaustedCall := fixture.fishDeck.On("Exhausted").Return(false).Once()
		prepareAfterRoundCall := fixture.playerMoveController.On("PrepareRound", mock.AnythingOfType("*match.State")).Run(func(arguments mock.Arguments) {
			state := arguments.Get(0).(*match.State)
			assert.Equal(t, 1, state.Round)
		}).Once()
		endConditionCall := fixture.endCondition.On("Apply", mock.AnythingOfType("*match.State")).Run(func(arguments mock.Arguments) {
			state := arguments.Get(0).(*match.State)
			assert.Equal(t, 1, state.RoundState.Thresholds.CaptureDistanceBonus)
			assert.Equal(t, 1, state.RoundState.Thresholds.SurfaceDepthBonus)
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
		assert.Equal(t, 6, result.State.Deck.ActiveCards)
		assert.Equal(t, 3, result.State.Deck.DiscardCards)
		assert.Equal(t, 1, result.State.Stats.PlayerWins)
		assert.Equal(t, match.RoundState{}, result.State.RoundState)
		assert.Equal(t, match.RoundState{}, fixture.engine.State().RoundState)
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
		playerCard := cards.NewPlayerCard(domain.Blue)

		initializeCall := playerMoveController.On("Initialize", mock.AnythingOfType("*match.State")).Run(func(arguments mock.Arguments) {
			state := arguments.Get(0).(*match.State)
			state.PlayerMoves = samplePlayerMoveResources()
		}).Once()
		prepareDeckCall := fishDeck.On("PrepareNextRound").Once()
		activeCountInitCall := fishDeck.On("ActiveCount").Return(0).Once()
		discardCountInitCall := fishDeck.On("DiscardCount").Return(0).Once()
		recycleCountInitCall := fishDeck.On("RecycleCount").Return(0).Once()
		exhaustedInitCall := fishDeck.On("Exhausted").Return(false).Once()

		prepareBeforeValidationCall := playerMoveController.On("PrepareRound", mock.AnythingOfType("*match.State")).Once()
		validateMoveCall := playerMoveController.On("ValidateMove", mock.Anything, domain.Blue).Return(nil).Once()
		peekCardCall := playerMoveController.On("PeekMoveCard", mock.Anything, domain.Blue).Return(playerCard, nil).Once()
		fishCard := cards.NewFishCard(domain.Red, cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1})
		drawCall := fishDeck.On("Draw").Return(fishCard, nil).Once()
		evaluateCall := roundEvaluator.On("Evaluate", domain.Blue, domain.Red).Return(domain.Draw).Once()
		consumeMoveCall := playerMoveController.On("ConsumeMove", mock.AnythingOfType("*match.State"), domain.Blue).Return(playerCard).Once()
		prepareDeckAfterRoundCall := fishDeck.On("PrepareNextRound").Once()
		activeCountAfterRoundCall := fishDeck.On("ActiveCount").Return(0).Once()
		discardCountAfterRoundCall := fishDeck.On("DiscardCount").Return(1).Once()
		recycleCountAfterRoundCall := fishDeck.On("RecycleCount").Return(0).Once()
		exhaustedAfterRoundCall := fishDeck.On("Exhausted").Return(false).Once()
		prepareAfterRoundCall := playerMoveController.On("PrepareRound", mock.AnythingOfType("*match.State")).Once()

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
			match.State{Encounter: encounterState, PlayerRig: samplePlayerRigState()},
		)
		require.NoError(t, err)

		result, err := engine.PlayRound(domain.Blue)

		require.NoError(t, err)
		assert.True(t, result.State.Finished)
		assert.Equal(t, encounter.StatusCaptured, result.State.Encounter.Status)
		assert.Equal(t, encounter.EndReasonTrackCapture, result.State.Encounter.EndReason)
		assert.Equal(t, match.RoundState{}, result.State.RoundState)
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
	if initialState.PlayerRig.MaxDistance == 0 {
		initialState.PlayerRig = samplePlayerRigState()
	}

	fishDeck := &mockFishDeck{}
	playerMoveController := &mockPlayerMoveController{}
	roundEvaluator := &mockRoundEvaluator{}
	progressionPolicy := &mockProgressionPolicy{}
	endCondition := &mockEndCondition{}

	initializeCall := playerMoveController.On("Initialize", mock.AnythingOfType("*match.State")).Run(func(arguments mock.Arguments) {
		state := arguments.Get(0).(*match.State)
		state.PlayerMoves = samplePlayerMoveResources()
	}).Once()
	prepareDeckCall := fishDeck.On("PrepareNextRound").Once()
	activeCountCall := fishDeck.On("ActiveCount").Return(7).Once()
	discardCountCall := fishDeck.On("DiscardCount").Return(2).Once()
	recycleCountCall := fishDeck.On("RecycleCount").Return(1).Once()
	exhaustedCall := fishDeck.On("Exhausted").Return(true).Once()
	endConditionCall := endCondition.On("Apply", mock.AnythingOfType("*match.State")).Once()
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

func samplePlayerRigState() playerrig.State {
	return playerrig.State{
		MaxDistance: 5,
		MaxDepth:    3,
	}
}

type mockFishDeck struct {
	mock.Mock
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

type mockPlayerMoveController struct {
	mock.Mock
}

func (mockController *mockPlayerMoveController) Initialize(state *match.State) {
	mockController.Called(state)
}

func (mockController *mockPlayerMoveController) PrepareRound(state *match.State) {
	mockController.Called(state)
}

func (mockController *mockPlayerMoveController) ValidateMove(state match.State, playerMove domain.Move) error {
	arguments := mockController.Called(state, playerMove)
	return arguments.Error(0)
}

func (mockController *mockPlayerMoveController) PeekMoveCard(state match.State, playerMove domain.Move) (cards.PlayerCard, error) {
	arguments := mockController.Called(state, playerMove)
	return arguments.Get(0).(cards.PlayerCard), arguments.Error(1)
}

func (mockController *mockPlayerMoveController) ConsumeMove(state *match.State, playerMove domain.Move) cards.PlayerCard {
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

func (mockPolicy *mockProgressionPolicy) Apply(state *match.State, round match.ResolvedRound) {
	mockPolicy.Called(state, round)
}

type mockEndCondition struct {
	mock.Mock
}

func (mockCondition *mockEndCondition) Apply(state *match.State) {
	mockCondition.Called(state)
}
