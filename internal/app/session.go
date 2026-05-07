package app

import (
	"fmt"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/match"
	"pesca/internal/presentation"
)

type UI interface {
	ShowIntro(presentation.IntroView) error
	ChooseMove(presentation.StatusView, []presentation.MoveOption) (domain.Move, error)
	ResolveSplash(presentation.SplashView) (encounter.SplashResolution, error)
	ShowRound(presentation.RoundView) error
	ShowGameOver(presentation.SummaryView) error
}

type Presenter interface {
	Intro() presentation.IntroView
	Status(match.StatusSnapshot) presentation.StatusView
	Splash(match.EncounterEventSnapshot, int) presentation.SplashView
	Round(match.RoundSnapshot) presentation.RoundView
	Summary(match.SummarySnapshot) presentation.SummaryView
}

type Engine interface {
	State() match.State
	PlayRound(domain.Move) (match.RoundResult, error)
	ResolveSplash(encounter.SplashResolution) error
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

	for !s.engine.State().Lifecycle.Finished {
		status := s.presenter.Status(match.NewStatusSnapshot(s.engine.State()))
		move, err := s.ui.ChooseMove(status, status.MoveOptions)
		if err != nil {
			return fmt.Errorf("choose move: %w", err)
		}

		result, err := s.engine.PlayRound(move)
		if err != nil {
			return fmt.Errorf("play round: %w", err)
		}

		for result.Encounter.Splash != nil {
			state := s.engine.State()
			resolution, err := s.ui.ResolveSplash(s.presenter.Splash(result.Encounter, state.Player.Loadout.SplashSuccessDistanceBonus()))
			if err != nil {
				return fmt.Errorf("resolve splash: %w", err)
			}
			if resolution.SuccessfulJumps > 0 {
				resolution.DistanceRewardApplied = state.Player.Loadout.SplashSuccessDistanceBonus() * resolution.SuccessfulJumps
			}
			if err := s.engine.ResolveSplash(resolution); err != nil {
				return fmt.Errorf("apply splash resolution: %w", err)
			}
			result.Status = match.NewStatusSnapshot(s.engine.State())
			result.Encounter = match.NewEncounterEventSnapshot(s.engine.State().Encounter)
		}

		if err := s.ui.ShowRound(s.presenter.Round(match.NewRoundSnapshot(result))); err != nil {
			return fmt.Errorf("show round: %w", err)
		}
	}

	if err := s.ui.ShowGameOver(s.presenter.Summary(match.NewSummarySnapshot(s.engine.State()))); err != nil {
		return fmt.Errorf("show game over: %w", err)
	}

	return nil
}
