package app_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"pesca/internal/app"
	"pesca/internal/deck"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/endings"
	"pesca/internal/game"
	"pesca/internal/presentation"
	"pesca/internal/progression"
	"pesca/internal/rules"
)

func TestNewSessionValidatesDependencies(t *testing.T) {
	engine := newEngineForSessionTest(t)
	ui := &mockUI{}
	presenter := &mockPresenter{}

	tests := []struct {
		name      string
		engine    *game.Engine
		ui        app.UI
		presenter app.Presenter
		wantErr   string
	}{
		{name: "requires engine", ui: ui, presenter: presenter, wantErr: "engine is required"},
		{name: "requires ui", engine: engine, presenter: presenter, wantErr: "ui is required"},
		{name: "requires presenter", engine: engine, ui: ui, wantErr: "presenter is required"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			session, err := app.NewSession(test.engine, test.ui, test.presenter)

			require.Error(t, err)
			assert.EqualError(t, err, test.wantErr)
			assert.Nil(t, session)
		})
	}
}

func TestSessionRunUsesInjectedDependencies(t *testing.T) {
	engine := newEngineForSessionTest(t)
	intro := presentation.IntroView{
		Title: "Pesca",
		Options: []presentation.MoveOption{
			{Index: 1, Move: domain.Blue, Label: "Tirar"},
			{Index: 2, Move: domain.Red, Label: "Recoger"},
			{Index: 3, Move: domain.Yellow, Label: "Soltar"},
		},
	}
	status := presentation.StatusView{}
	round := presentation.RoundView{}
	summary := presentation.SummaryView{TotalRounds: 18}

	ui := &mockUI{}
	presenter := &mockPresenter{}
	presenter.On("Intro").Return(intro).Once()

	session, err := app.NewSession(engine, ui, presenter)
	require.NoError(t, err)

	ui.On("ShowIntro", intro).Return(nil).Once()
	presenter.On("Status", mock.Anything).Return(status).Times(18)
	ui.On("ChooseMove", status, intro.Options).Return(domain.Blue, nil).Times(18)
	presenter.On("Round", mock.Anything).Return(round).Times(18)
	ui.On("ShowRound", round).Return(nil).Times(18)
	presenter.On("Summary", mock.Anything).Return(summary).Once()
	ui.On("ShowGameOver", summary).Return(nil).Once()

	require.NoError(t, session.Run())

	ui.AssertExpectations(t)
	presenter.AssertExpectations(t)
	ui.AssertNumberOfCalls(t, "ChooseMove", 18)
	ui.AssertNumberOfCalls(t, "ShowRound", 18)
	presenter.AssertNumberOfCalls(t, "Status", 18)
	presenter.AssertNumberOfCalls(t, "Round", 18)
}

func newEngineForSessionTest(t *testing.T) *game.Engine {
	t.Helper()

	encounterState, err := encounter.NewState(encounter.Config{
		InitialDistance:           3,
		CaptureDistance:           -1,
		EscapeDistance:            99,
		ExhaustionCaptureDistance: 2,
		PlayerWinStep:             1,
		FishWinStep:               1,
	})
	require.NoError(t, err)

	engine, err := game.NewEngine(
		deck.NewManager(
			deck.NewStandardFishDeck(),
			func([]domain.Move) {},
			deck.RemoveCardsRecyclePolicy{CardsToRemove: 3},
		),
		rules.NewClassicEvaluator(rules.NewFishCombatProfile()),
		progression.TrackPolicy{},
		endings.EncounterCondition{},
		game.State{Encounter: encounterState},
	)
	require.NoError(t, err)

	return engine
}

type mockUI struct {
	mock.Mock
}

func (ui *mockUI) ShowIntro(view presentation.IntroView) error {
	return ui.Called(view).Error(0)
}

func (ui *mockUI) ChooseMove(view presentation.StatusView, options []presentation.MoveOption) (domain.Move, error) {
	args := ui.Called(view, options)
	return args.Get(0).(domain.Move), args.Error(1)
}

func (ui *mockUI) ShowRound(view presentation.RoundView) error {
	return ui.Called(view).Error(0)
}

func (ui *mockUI) ShowGameOver(view presentation.SummaryView) error {
	return ui.Called(view).Error(0)
}

type mockPresenter struct {
	mock.Mock
}

func (presenter *mockPresenter) Intro() presentation.IntroView {
	return presenter.Called().Get(0).(presentation.IntroView)
}

func (presenter *mockPresenter) Status(state game.State) presentation.StatusView {
	return presenter.Called(state).Get(0).(presentation.StatusView)
}

func (presenter *mockPresenter) Round(result game.RoundResult) presentation.RoundView {
	return presenter.Called(result).Get(0).(presentation.RoundView)
}

func (presenter *mockPresenter) Summary(state game.State) presentation.SummaryView {
	return presenter.Called(state).Get(0).(presentation.SummaryView)
}
