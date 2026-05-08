package app

import (
	"errors"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/playerprofiles"
	"pesca/internal/player/loadout"
	"pesca/internal/player/rod"
	"pesca/internal/presentation"
	"pesca/internal/run"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultRunLoopConfigStateProvider(t *testing.T) {
	provider := &defaultRunLoopConfig{}
	playerRod, err := rod.NewState(rod.DefaultConfig())
	require.NoError(t, err)
	playerLoadout, err := loadout.NewState(playerRod, nil)
	require.NoError(t, err)
	route := run.DefaultRoute()

	state, err := provider.Initialize(playerLoadout, route, 3)
	require.NoError(t, err)
	assert.Equal(t, run.StatusInProgress, provider.CurrentStatus(&state))
	assert.Equal(t, route[0], provider.CurrentNode(&state))
	require.NotNil(t, provider.NextNode(&state))

	require.NoError(t, provider.Advance(&state, route))
	assert.Equal(t, route[1], provider.CurrentNode(&state))

	result := run.EncounterResult{
		Outcome:       run.EncounterOutcomeEscaped,
		Status:        "escaped",
		EndReason:     "track_escape",
		ThreadDamage:  1,
		NodeResolved:  true,
		FinishedMatch: true,
	}
	require.NoError(t, provider.ApplyResult(&state, result))
	assert.Equal(t, 2, state.Thread.Current)

	require.NoError(t, provider.Complete(&state))
	assert.Equal(t, run.StatusVictory, provider.CurrentStatus(&state))
}

func TestDefaultEncounterHandlerResolve(t *testing.T) {
	handler := &defaultEncounterHandler{}
	presenter := presentation.NewPresenter(presentation.DefaultCatalog())
	_, err := handler.Resolve(
		"Test Run",
		nil,
		&mockRunUI{},
		&mockRunUI{},
		presenter,
		playerprofiles.DefaultPresets()[0],
		loadout.State{},
		RunEncounterBootstrapConfig{},
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "randomizer is required")
}

func TestDefaultNodeRenderer(t *testing.T) {
	ui := &mockRunUI{}
	presenter := presentation.NewPresenter(presentation.DefaultCatalog())
	renderer := &defaultNodeRenderer{ui: ui, presenter: presenter}
	state := run.State{}
	route := []run.NodeState{{NodeID: "start", Kind: run.NodeKindStart}}
	result := run.EncounterResult{Outcome: run.EncounterOutcomeCaptured}
	nextNode := &run.NodeState{NodeID: "end", Kind: run.NodeKindEnd}

	ui.On("ShowRunIntro", presenter.RunIntro("Test Run", state, route)).Return(nil).Once()
	ui.On("ShowRunNode", presenter.RunNode("Test Run", state)).Return(nil).Once()
	ui.On("ShowRunNodeSummary", presenter.RunNodeSummary("Test Run", route[0], result, state, nextNode)).Return(nil).Once()
	ui.On("ShowRunSummary", presenter.RunSummary("Test Run", state)).Return(nil).Once()

	require.NoError(t, renderer.ShowIntro("Test Run", state, route))
	require.NoError(t, renderer.ShowNode("Test Run", state))
	require.NoError(t, renderer.ShowNodeSummary("Test Run", route[0], result, state, nextNode))
	require.NoError(t, renderer.ShowSummary("Test Run", state))
	ui.AssertExpectations(t)
}

func TestNoopLoopCallbacks(t *testing.T) {
	callbacks := &noopLoopCallbacks{}
	state := &run.State{}
	result := run.EncounterResult{Outcome: run.EncounterOutcomeCaptured}
	require.NoError(t, callbacks.OnNodeStart(run.NodeState{}, state))
	require.NoError(t, callbacks.OnNodeComplete(run.NodeState{}, result, state))
	require.NoError(t, callbacks.OnRunComplete(run.State{}))
}

func TestBuildEncounterConfig(t *testing.T) {
	session := &RunSession{}
	assert.Equal(t, fishprofiles.DefaultEncounterFishPoolID, session.buildEncounterConfig(run.State{}).FishPoolID)
	assert.Equal(t, fishprofiles.DefaultEncounterFishPoolID, session.buildEncounterConfig(run.State{Progress: run.ProgressState{Current: run.NodeState{Kind: run.NodeKindBoss}}}).FishPoolID)
}

func TestRunSessionSetLoopConfigAllowsInjection(t *testing.T) {
	session := &RunSession{}
	config := &RunLoopConfig{}
	session.setLoopConfig(config)
	assert.Same(t, config, session.loopConfig)
	session.setLoopConfig(nil)
	assert.Nil(t, session.loopConfig)
}

func TestMockEncounterHandlerSignature(t *testing.T) {
	_ = errors.New
	_ = fishprofiles.Spawn{}
}
