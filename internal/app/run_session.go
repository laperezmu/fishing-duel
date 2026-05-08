package app

import (
	"fmt"
	"pesca/internal/content/anglerprofiles"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/presentation"
	"pesca/internal/run"
)

type RunUI interface {
	EncounterBootstrapUI
	RunSetupUI
	UI
	ShowRunIntro(presentation.RunIntroView) error
	ShowRunNode(presentation.RunNodeView) error
	ShowRunNodeSummary(presentation.RunNodeSummaryView) error
	ShowRunSummary(presentation.RunSummaryView) error
}

type RunSession struct {
	title      string
	rng        Randomizer
	ui         RunUI
	presenter  presentation.Presenter
	route      []run.NodeState
	loopConfig *RunLoopConfig
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

	session := &RunSession{
		title:      title,
		rng:        rng,
		ui:         ui,
		presenter:  presenter,
		route:      route,
		loopConfig: nil,
	}

	session.loopConfig = session.createDefaultLoopConfig()

	return session, nil
}

func (session *RunSession) createDefaultLoopConfig() *RunLoopConfig {
	return &RunLoopConfig{
		StateProvider:    &defaultRunLoopConfig{},
		EncounterHandler: &defaultEncounterHandler{},
		NodeRenderer:     &defaultNodeRenderer{ui: session.ui, presenter: session.presenter},
		LoopCallbacks:    &noopLoopCallbacks{},
	}
}

func (session *RunSession) Run() error {
	resolvedStart, err := resolveRunSetup(session.title, session.ui)
	if err != nil {
		return err
	}

	state, err := session.initializeRunState(resolvedStart)
	if err != nil {
		return err
	}

	for session.loopConfig.StateProvider.CurrentStatus(&state) == run.StatusInProgress {
		if err := session.playCurrentNode(&state, resolvedStart); err != nil {
			return err
		}
	}

	if err := session.loopConfig.NodeRenderer.ShowSummary(session.title, state); err != nil {
		return fmt.Errorf("show run summary: %w", err)
	}

	return nil
}

func (session *RunSession) initializeRunState(resolvedStart anglerprofiles.ResolvedStart) (run.State, error) {
	state, err := session.loopConfig.StateProvider.Initialize(resolvedStart.Loadout, session.route, resolvedStart.StartingThread)
	if err != nil {
		return run.State{}, fmt.Errorf("initialize run state: %w", err)
	}
	if err := session.loopConfig.NodeRenderer.ShowIntro(session.title, state, session.route); err != nil {
		return run.State{}, fmt.Errorf("show run intro: %w", err)
	}

	return state, nil
}

func (session *RunSession) playCurrentNode(state *run.State, resolvedStart anglerprofiles.ResolvedStart) error {
	currentNode := session.loopConfig.StateProvider.CurrentNode(state)

	if err := session.loopConfig.NodeRenderer.ShowNode(session.title, *state); err != nil {
		return fmt.Errorf("show run node: %w", err)
	}

	if err := session.loopConfig.LoopCallbacks.OnNodeStart(currentNode, state); err != nil {
		return err
	}

	switch currentNode.Kind {
	case run.NodeKindStart, run.NodeKindService, run.NodeKindCheckpoint:
		if err := session.advanceRunNode(state); err != nil {
			return err
		}
	case run.NodeKindFishing, run.NodeKindBoss:
		if err := session.playEncounterNode(state, resolvedStart); err != nil {
			return err
		}
	case run.NodeKindEnd:
		if err := session.loopConfig.StateProvider.Complete(state); err != nil {
			return fmt.Errorf("complete run: %w", err)
		}
	default:
		return fmt.Errorf("unsupported node kind %q", currentNode.Kind)
	}

	return nil
}

func (session *RunSession) advanceRunNode(state *run.State) error {
	if err := session.loopConfig.StateProvider.Advance(state, session.route); err != nil {
		return fmt.Errorf("advance run node: %w", err)
	}

	return nil
}

func (session *RunSession) playEncounterNode(state *run.State, resolvedStart anglerprofiles.ResolvedStart) error {
	currentNode := session.loopConfig.StateProvider.CurrentNode(state)
	nextNode := session.loopConfig.StateProvider.NextNode(state)
	title := fmt.Sprintf("%s - %s", session.title, currentNode.NodeID)
	waterPreset, err := run.ResolveWaterPreset(currentNode.WaterPresetID)
	if err != nil {
		return fmt.Errorf("resolve node water preset: %w", err)
	}

	resolvedEncounter, err := session.loopConfig.EncounterHandler.Resolve(
		title,
		session.rng,
		session.ui,
		session.ui,
		session.presenter,
		resolvedStart.DeckPreset,
		state.Loadout,
		RunEncounterBootstrapConfig{
			Encounter: session.buildEncounterConfig(*state),
			Water:     waterPreset,
		},
	)
	if err != nil {
		return fmt.Errorf("bootstrap encounter: %w", err)
	}

	encounterSession, err := NewSession(resolvedEncounter.Engine, session.ui, session.presenter)
	if err != nil {
		return fmt.Errorf("initialize encounter session: %w", err)
	}
	if err := encounterSession.Run(); err != nil {
		return fmt.Errorf("run encounter session: %w", err)
	}

	encounterResult, err := ResolveEncounterResult(resolvedEncounter.Engine.State(), resolvedEncounter.Spawn)
	if err != nil {
		return fmt.Errorf("resolve encounter result: %w", err)
	}
	if err := session.loopConfig.StateProvider.ApplyResult(state, encounterResult); err != nil {
		return fmt.Errorf("apply encounter result: %w", err)
	}

	if session.loopConfig.StateProvider.CurrentStatus(state) != run.StatusInProgress {
		return nil
	}

	if err := session.loopConfig.NodeRenderer.ShowNodeSummary(session.title, currentNode, encounterResult, *state, nextNode); err != nil {
		return fmt.Errorf("show run node summary: %w", err)
	}

	if err := session.loopConfig.LoopCallbacks.OnNodeComplete(currentNode, encounterResult, state); err != nil {
		return err
	}

	return session.advanceRunNode(state)
}

func (session *RunSession) buildEncounterConfig(state run.State) EncounterBootstrapConfig {
	if state.Progress.Current.Kind == run.NodeKindBoss {
		return EncounterBootstrapConfig{FishPoolID: fishprofiles.DefaultEncounterFishPoolID}
	}

	return EncounterBootstrapConfig{FishPoolID: fishprofiles.DefaultEncounterFishPoolID}
}

func (session *RunSession) setLoopConfig(config *RunLoopConfig) {
	session.loopConfig = config
}
