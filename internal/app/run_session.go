package app

import (
	"fmt"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/playerprofiles"
	"pesca/internal/player/loadout"
	"pesca/internal/presentation"
	"pesca/internal/run"
)

type RunUI interface {
	EncounterBootstrapUI
	UI
	ShowRunIntro(presentation.RunIntroView) error
	ShowRunNode(presentation.RunNodeView) error
	ShowRunSummary(presentation.RunSummaryView) error
}

type RunSession struct {
	title     string
	rng       Randomizer
	ui        RunUI
	presenter presentation.Presenter
	route     []run.NodeState
}

func NewRunSession(title string, rng Randomizer, ui RunUI, presenter presentation.Presenter) (*RunSession, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if rng == nil {
		return nil, fmt.Errorf("randomizer is required")
	}
	if ui == nil {
		return nil, fmt.Errorf("run ui is required")
	}

	route := run.DefaultRoute()
	if len(route) == 0 {
		return nil, fmt.Errorf("run route is required")
	}

	return &RunSession{title: title, rng: rng, ui: ui, presenter: presenter, route: route}, nil
}

func (session *RunSession) Run() error {
	playerDeckPreset, playerLoadout, err := resolvePlayerSetup(session.title, session.ui)
	if err != nil {
		return err
	}

	state, err := session.initializeRunState(playerLoadout)
	if err != nil {
		return err
	}

	for state.Status == run.StatusInProgress {
		if err := session.playCurrentNode(&state, playerDeckPreset, playerLoadout); err != nil {
			return err
		}
	}

	if err := session.ui.ShowRunSummary(session.presenter.RunSummary(session.title, state)); err != nil {
		return fmt.Errorf("show run summary: %w", err)
	}

	return nil
}

func (session *RunSession) initializeRunState(playerLoadout loadout.State) (run.State, error) {
	state, err := run.NewState(playerLoadout, session.route, run.DefaultThreadMaximum)
	if err != nil {
		return run.State{}, fmt.Errorf("initialize run state: %w", err)
	}
	if err := session.ui.ShowRunIntro(session.presenter.RunIntro(session.title, state, session.route)); err != nil {
		return run.State{}, fmt.Errorf("show run intro: %w", err)
	}

	return state, nil
}

func (session *RunSession) playCurrentNode(state *run.State, playerDeckPreset playerprofiles.DeckPreset, playerLoadout loadout.State) error {
	if err := session.ui.ShowRunNode(session.presenter.RunNode(session.title, *state)); err != nil {
		return fmt.Errorf("show run node: %w", err)
	}

	switch state.Progress.Current.Kind {
	case run.NodeKindStart, run.NodeKindService, run.NodeKindCheckpoint:
		return session.advanceRunNode(state)
	case run.NodeKindFishing, run.NodeKindBoss:
		return session.playEncounterNode(state, playerDeckPreset, playerLoadout)
	case run.NodeKindEnd:
		if err := run.Complete(state); err != nil {
			return fmt.Errorf("complete run: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unsupported node kind %q", state.Progress.Current.Kind)
	}
}

func (session *RunSession) advanceRunNode(state *run.State) error {
	if err := run.Advance(state, session.route); err != nil {
		return fmt.Errorf("advance run node: %w", err)
	}

	return nil
}

func (session *RunSession) playEncounterNode(state *run.State, playerDeckPreset playerprofiles.DeckPreset, playerLoadout loadout.State) error {
	title := fmt.Sprintf("%s - %s", session.title, state.Progress.Current.NodeID)
	bootstrap, err := BootstrapEncounterForRun(title, session.rng, session.ui, session.presenter, playerDeckPreset, playerLoadout, session.buildEncounterConfig(*state))
	if err != nil {
		return fmt.Errorf("bootstrap encounter: %w", err)
	}

	encounterSession, err := NewSession(bootstrap.Engine, session.ui, session.presenter)
	if err != nil {
		return fmt.Errorf("initialize encounter session: %w", err)
	}
	if err := encounterSession.Run(); err != nil {
		return fmt.Errorf("run encounter session: %w", err)
	}

	encounterResult, err := ResolveEncounterResult(bootstrap.Engine.State(), bootstrap.Spawn)
	if err != nil {
		return fmt.Errorf("resolve encounter result: %w", err)
	}
	if err := run.ApplyEncounterResult(state, encounterResult); err != nil {
		return fmt.Errorf("apply encounter result: %w", err)
	}
	if state.Status != run.StatusInProgress {
		return nil
	}
	return session.advanceRunNode(state)
}

func (session *RunSession) buildEncounterConfig(state run.State) EncounterBootstrapConfig {
	if state.Progress.Current.Kind == run.NodeKindBoss {
		return EncounterBootstrapConfig{FishPoolID: fishprofiles.DefaultEncounterFishPoolID}
	}

	return EncounterBootstrapConfig{FishPoolID: fishprofiles.DefaultEncounterFishPoolID}
}
