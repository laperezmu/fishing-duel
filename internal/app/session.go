package app

import (
	"fmt"

	"pesca/internal/domain"
	"pesca/internal/game"
	"pesca/internal/presentation"
)

type UI interface {
	ShowIntro(presentation.IntroView) error
	ChooseMove(presentation.StatusView, []presentation.MoveOption) (domain.Move, error)
	ShowRound(presentation.RoundView) error
	ShowGameOver(presentation.SummaryView) error
}

type Presenter interface {
	Intro() presentation.IntroView
	Status(game.State) presentation.StatusView
	Round(game.RoundResult) presentation.RoundView
	Summary(game.State) presentation.SummaryView
}

type Engine interface {
	State() game.State
	PlayRound(domain.Move) (game.RoundResult, error)
}

type Session struct {
	engine    Engine
	ui        UI
	presenter Presenter
	intro     presentation.IntroView
}

func NewSession(engine Engine, ui UI, presenter Presenter) (*Session, error) {
	if engine == nil {
		return nil, fmt.Errorf("engine is required")
	}
	if ui == nil {
		return nil, fmt.Errorf("ui is required")
	}
	if presenter == nil {
		return nil, fmt.Errorf("presenter is required")
	}

	intro := presenter.Intro()
	return &Session{
		engine:    engine,
		ui:        ui,
		presenter: presenter,
		intro:     intro,
	}, nil
}

func (s *Session) Run() error {
	if err := s.ui.ShowIntro(s.intro); err != nil {
		return fmt.Errorf("show intro: %w", err)
	}

	for !s.engine.State().Finished {
		status := s.presenter.Status(s.engine.State())
		move, err := s.ui.ChooseMove(status, s.intro.Options)
		if err != nil {
			return fmt.Errorf("choose move: %w", err)
		}

		result, err := s.engine.PlayRound(move)
		if err != nil {
			return fmt.Errorf("play round: %w", err)
		}

		if err := s.ui.ShowRound(s.presenter.Round(result)); err != nil {
			return fmt.Errorf("show round: %w", err)
		}
	}

	if err := s.ui.ShowGameOver(s.presenter.Summary(s.engine.State())); err != nil {
		return fmt.Errorf("show game over: %w", err)
	}

	return nil
}
