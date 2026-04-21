package app_test

import (
	"errors"
	"pesca/internal/app"
	"pesca/internal/domain"
	"pesca/internal/match"
	"pesca/internal/presentation"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewSession(t *testing.T) {
	t.Run("returns a session when engine ui and presenter are provided", func(t *testing.T) {
		engine := &mockEngine{}
		ui := &mockUI{}
		presenter := &mockPresenter{}
		intro := presentation.IntroView{Title: "Pesca"}

		presenter.On("Intro").Return(intro).Once()

		session, err := app.NewSession(engine, ui, presenter)

		require.NoError(t, err)
		assert.NotNil(t, session)
		presenter.AssertExpectations(t)
	})

	missingDependencyCases := []struct {
		title     string
		engine    app.Engine
		ui        app.UI
		presenter app.Presenter
		wantErr   string
	}{
		{
			title:     "returns an error when engine is missing",
			ui:        &mockUI{},
			presenter: &mockPresenter{},
			wantErr:   "engine is required",
		},
		{
			title:     "returns an error when ui is missing",
			engine:    &mockEngine{},
			presenter: &mockPresenter{},
			wantErr:   "ui is required",
		},
		{
			title:   "returns an error when presenter is missing",
			engine:  &mockEngine{},
			ui:      &mockUI{},
			wantErr: "presenter is required",
		},
	}

	for _, test := range missingDependencyCases {
		t.Run(test.title, func(t *testing.T) {
			session, err := app.NewSession(test.engine, test.ui, test.presenter)

			require.Error(t, err)
			assert.EqualError(t, err, test.wantErr)
			assert.Nil(t, session)
		})
	}
}

func TestSessionRun(t *testing.T) {
	t.Run("shows intro runs one round and shows summary when the engine finishes", func(t *testing.T) {
		fixture := newSessionFixture(t)

		fixture.engine.On("State").Return(fixture.ongoingState).Twice()
		fixture.engine.On("PlayRound", domain.Blue).Return(fixture.roundResult, nil).Once()
		fixture.engine.On("State").Return(fixture.finishedState).Twice()
		fixture.ui.On("ShowIntro", fixture.intro).Return(nil).Once()
		fixture.presenter.On("Status", fixture.ongoingState).Return(fixture.status).Once()
		fixture.ui.On("ChooseMove", fixture.status, fixture.intro.Options).Return(domain.Blue, nil).Once()
		fixture.presenter.On("Round", fixture.roundResult).Return(fixture.round).Once()
		fixture.ui.On("ShowRound", fixture.round).Return(nil).Once()
		fixture.presenter.On("Summary", fixture.finishedState).Return(fixture.summary).Once()
		fixture.ui.On("ShowGameOver", fixture.summary).Return(nil).Once()

		require.NoError(t, fixture.session.Run())

		fixture.engine.AssertExpectations(t)
		fixture.ui.AssertExpectations(t)
		fixture.presenter.AssertExpectations(t)
	})

	errorCases := []struct {
		title   string
		setup   func(fixture sessionFixture)
		wantErr string
	}{
		{
			title: "returns a wrapped error when showing the intro fails",
			setup: func(fixture sessionFixture) {
				fixture.ui.On("ShowIntro", fixture.intro).Return(errors.New("intro failed")).Once()
			},
			wantErr: "show intro: intro failed",
		},
		{
			title: "returns a wrapped error when choosing a move fails",
			setup: func(fixture sessionFixture) {
				fixture.engine.On("State").Return(fixture.ongoingState).Twice()
				fixture.ui.On("ShowIntro", fixture.intro).Return(nil).Once()
				fixture.presenter.On("Status", fixture.ongoingState).Return(fixture.status).Once()
				fixture.ui.On("ChooseMove", fixture.status, fixture.intro.Options).Return(domain.Blue, errors.New("choose failed")).Once()
			},
			wantErr: "choose move: choose failed",
		},
		{
			title: "returns a wrapped error when playing a round fails",
			setup: func(fixture sessionFixture) {
				fixture.engine.On("State").Return(fixture.ongoingState).Twice()
				fixture.ui.On("ShowIntro", fixture.intro).Return(nil).Once()
				fixture.presenter.On("Status", fixture.ongoingState).Return(fixture.status).Once()
				fixture.ui.On("ChooseMove", fixture.status, fixture.intro.Options).Return(domain.Blue, nil).Once()
				fixture.engine.On("PlayRound", domain.Blue).Return(match.RoundResult{}, errors.New("round failed")).Once()
			},
			wantErr: "play round: round failed",
		},
		{
			title: "returns a wrapped error when showing the round fails",
			setup: func(fixture sessionFixture) {
				fixture.engine.On("State").Return(fixture.ongoingState).Twice()
				fixture.ui.On("ShowIntro", fixture.intro).Return(nil).Once()
				fixture.presenter.On("Status", fixture.ongoingState).Return(fixture.status).Once()
				fixture.ui.On("ChooseMove", fixture.status, fixture.intro.Options).Return(domain.Blue, nil).Once()
				fixture.engine.On("PlayRound", domain.Blue).Return(fixture.roundResult, nil).Once()
				fixture.presenter.On("Round", fixture.roundResult).Return(fixture.round).Once()
				fixture.ui.On("ShowRound", fixture.round).Return(errors.New("round view failed")).Once()
			},
			wantErr: "show round: round view failed",
		},
		{
			title: "returns a wrapped error when showing the summary fails",
			setup: func(fixture sessionFixture) {
				fixture.engine.On("State").Return(fixture.ongoingState).Twice()
				fixture.engine.On("PlayRound", domain.Blue).Return(fixture.roundResult, nil).Once()
				fixture.engine.On("State").Return(fixture.finishedState).Twice()
				fixture.ui.On("ShowIntro", fixture.intro).Return(nil).Once()
				fixture.presenter.On("Status", fixture.ongoingState).Return(fixture.status).Once()
				fixture.ui.On("ChooseMove", fixture.status, fixture.intro.Options).Return(domain.Blue, nil).Once()
				fixture.presenter.On("Round", fixture.roundResult).Return(fixture.round).Once()
				fixture.ui.On("ShowRound", fixture.round).Return(nil).Once()
				fixture.presenter.On("Summary", fixture.finishedState).Return(fixture.summary).Once()
				fixture.ui.On("ShowGameOver", fixture.summary).Return(errors.New("summary failed")).Once()
			},
			wantErr: "show game over: summary failed",
		},
	}

	for _, test := range errorCases {
		t.Run(test.title, func(t *testing.T) {
			fixture := newSessionFixture(t)
			test.setup(fixture)

			err := fixture.session.Run()

			require.Error(t, err)
			assert.EqualError(t, err, test.wantErr)
			fixture.engine.AssertExpectations(t)
			fixture.ui.AssertExpectations(t)
			fixture.presenter.AssertExpectations(t)
		})
	}
}

type sessionFixture struct {
	engine        *mockEngine
	ui            *mockUI
	presenter     *mockPresenter
	session       *app.Session
	intro         presentation.IntroView
	status        presentation.StatusView
	round         presentation.RoundView
	summary       presentation.SummaryView
	ongoingState  match.State
	finishedState match.State
	roundResult   match.RoundResult
}

func newSessionFixture(t *testing.T) sessionFixture {
	t.Helper()

	engine := &mockEngine{}
	ui := &mockUI{}
	presenter := &mockPresenter{}
	intro := presentation.IntroView{
		Title: "Pesca",
		Options: []presentation.MoveOption{
			{Index: 1, Move: domain.Blue, Label: "Tirar"},
			{Index: 2, Move: domain.Red, Label: "Recoger"},
			{Index: 3, Move: domain.Yellow, Label: "Soltar"},
		},
	}
	status := presentation.StatusView{RoundNumber: 1}
	ongoingState := match.State{Round: 0, Finished: false}
	finishedState := match.State{Round: 1, Finished: true}
	roundResult := match.RoundResult{Round: 1, PlayerMove: domain.Blue, Outcome: domain.PlayerWin, State: finishedState}
	round := presentation.RoundView{Outcome: domain.PlayerWin, OutcomeLabel: "gana el jugador"}
	summary := presentation.SummaryView{TotalRounds: 1}

	presenter.On("Intro").Return(intro).Once()
	session, err := app.NewSession(engine, ui, presenter)
	require.NoError(t, err)

	return sessionFixture{
		engine:        engine,
		ui:            ui,
		presenter:     presenter,
		session:       session,
		intro:         intro,
		status:        status,
		round:         round,
		summary:       summary,
		ongoingState:  ongoingState,
		finishedState: finishedState,
		roundResult:   roundResult,
	}
}

type mockEngine struct {
	mock.Mock
}

func (engine *mockEngine) State() match.State {
	return engine.Called().Get(0).(match.State)
}

func (engine *mockEngine) PlayRound(move domain.Move) (match.RoundResult, error) {
	args := engine.Called(move)
	return args.Get(0).(match.RoundResult), args.Error(1)
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

func (presenter *mockPresenter) Status(state match.State) presentation.StatusView {
	return presenter.Called(state).Get(0).(presentation.StatusView)
}

func (presenter *mockPresenter) Round(result match.RoundResult) presentation.RoundView {
	return presenter.Called(result).Get(0).(presentation.RoundView)
}

func (presenter *mockPresenter) Summary(state match.State) presentation.SummaryView {
	return presenter.Called(state).Get(0).(presentation.SummaryView)
}
