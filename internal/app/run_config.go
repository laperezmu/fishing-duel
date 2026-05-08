package app

import (
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/playerprofiles"
	"pesca/internal/game"
	"pesca/internal/player/loadout"
	"pesca/internal/presentation"
	"pesca/internal/run"
)

type RunLoopConfig struct {
	StateProvider    RunStateProvider
	EncounterHandler EncounterHandler
	NodeRenderer     NodeRenderer
	LoopCallbacks    RunLoopCallbacks
}

type RunStateProvider interface {
	Initialize(loadout loadout.State, route []run.NodeState, thread int) (run.State, error)
	Advance(state *run.State, route []run.NodeState) error
	ApplyResult(state *run.State, result run.EncounterResult) error
	Complete(state *run.State) error
	CurrentStatus(state *run.State) run.Status
	CurrentNode(state *run.State) run.NodeState
	NextNode(state *run.State) *run.NodeState
}

type EncounterHandler interface {
	Resolve(
		title string,
		rng Randomizer,
		ui EncounterBootstrapUI,
		spawnUI SpawnUI,
		presenter presentation.Presenter,
		deckPreset playerprofiles.DeckPreset,
		playerLoadout loadout.State,
		config RunEncounterBootstrapConfig,
	) (ResolvedEncounter, error)
}

type ResolvedEncounter struct {
	Engine *game.Engine
	Spawn  fishprofiles.Spawn
}

type NodeRenderer interface {
	ShowIntro(title string, state run.State, route []run.NodeState) error
	ShowNode(title string, state run.State) error
	ShowNodeSummary(title string, currentNode run.NodeState, result run.EncounterResult, state run.State, nextNode *run.NodeState) error
	ShowSummary(title string, state run.State) error
}

type RunLoopCallbacks interface {
	OnNodeStart(node run.NodeState, state *run.State) error
	OnNodeComplete(node run.NodeState, result run.EncounterResult, state *run.State) error
	OnRunComplete(state run.State) error
}

type defaultRunLoopConfig struct{}

func (c *defaultRunLoopConfig) Initialize(loadout loadout.State, route []run.NodeState, thread int) (run.State, error) {
	return run.NewState(loadout, route, thread)
}

func (c *defaultRunLoopConfig) Advance(state *run.State, route []run.NodeState) error {
	return run.Advance(state, route)
}

func (c *defaultRunLoopConfig) ApplyResult(state *run.State, result run.EncounterResult) error {
	return run.ApplyEncounterResult(state, result)
}

func (c *defaultRunLoopConfig) Complete(state *run.State) error {
	return run.Complete(state)
}

func (c *defaultRunLoopConfig) CurrentStatus(state *run.State) run.Status {
	return state.Status
}

func (c *defaultRunLoopConfig) CurrentNode(state *run.State) run.NodeState {
	return state.Progress.Current
}

func (c *defaultRunLoopConfig) NextNode(state *run.State) *run.NodeState {
	return state.Progress.Next
}

type defaultEncounterHandler struct{}

func (h *defaultEncounterHandler) Resolve(
	title string,
	rng Randomizer,
	ui EncounterBootstrapUI,
	spawnUI SpawnUI,
	presenter presentation.Presenter,
	deckPreset playerprofiles.DeckPreset,
	playerLoadout loadout.State,
	config RunEncounterBootstrapConfig,
) (ResolvedEncounter, error) {
	bootstrap, err := BootstrapEncounterForRun(title, rng, ui, spawnUI, presenter, deckPreset, playerLoadout, config)
	if err != nil {
		return ResolvedEncounter{}, err
	}

	return ResolvedEncounter{Engine: bootstrap.Engine, Spawn: bootstrap.Spawn}, nil
}

type defaultNodeRenderer struct {
	ui        RunUI
	presenter presentation.Presenter
}

func (r *defaultNodeRenderer) ShowIntro(title string, state run.State, route []run.NodeState) error {
	return r.ui.ShowRunIntro(r.presenter.RunIntro(title, state, route))
}

func (r *defaultNodeRenderer) ShowNode(title string, state run.State) error {
	return r.ui.ShowRunNode(r.presenter.RunNode(title, state))
}

func (r *defaultNodeRenderer) ShowNodeSummary(title string, currentNode run.NodeState, result run.EncounterResult, state run.State, nextNode *run.NodeState) error {
	return r.ui.ShowRunNodeSummary(r.presenter.RunNodeSummary(title, currentNode, result, state, nextNode))
}

func (r *defaultNodeRenderer) ShowSummary(title string, state run.State) error {
	return r.ui.ShowRunSummary(r.presenter.RunSummary(title, state))
}

type noopLoopCallbacks struct{}

func (c *noopLoopCallbacks) OnNodeStart(node run.NodeState, state *run.State) error { return nil }
func (c *noopLoopCallbacks) OnNodeComplete(node run.NodeState, result run.EncounterResult, state *run.State) error {
	return nil
}
func (c *noopLoopCallbacks) OnRunComplete(state run.State) error { return nil }
