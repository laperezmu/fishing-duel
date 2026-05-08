package app

import (
	"errors"
	"pesca/internal/content/anglerprofiles"
	"pesca/internal/content/attachmentpresets"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/playerprofiles"
	"pesca/internal/content/rodpresets"
	"pesca/internal/content/watercontexts"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/player/loadout"
	"pesca/internal/player/rod"
	"pesca/internal/presentation"
	"pesca/internal/run"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewRunSession(t *testing.T) {
	t.Run("returns a run session when all dependencies are provided", func(t *testing.T) {
		engine := &mockRNG{}
		ui := &mockRunUI{}
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())

		session, err := NewRunSession("Test Run", engine, ui, presenter)

		require.NoError(t, err)
		assert.NotNil(t, session)
	})

	t.Run("returns an error when title is missing", func(t *testing.T) {
		engine := &mockRNG{}
		ui := &mockRunUI{}
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())

		session, err := NewRunSession("", engine, ui, presenter)

		require.Error(t, err)
		assert.EqualError(t, err, "title is required")
		assert.Nil(t, session)
	})

	t.Run("returns an error when rng is missing", func(t *testing.T) {
		ui := &mockRunUI{}
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())

		session, err := NewRunSession("Test Run", nil, ui, presenter)

		require.Error(t, err)
		assert.EqualError(t, err, "randomizer is required")
		assert.Nil(t, session)
	})

	t.Run("returns an error when ui is missing", func(t *testing.T) {
		engine := &mockRNG{}
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())

		session, err := NewRunSession("Test Run", engine, nil, presenter)

		require.Error(t, err)
		assert.EqualError(t, err, "run ui is required")
		assert.Nil(t, session)
	})
}

func TestRunSessionRun(t *testing.T) {
	t.Run("returns error when setup fails", func(t *testing.T) {
		rng := &mockRNG{}
		ui := &mockRunUI{}
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())
		ui.On("ChooseAnglerProfile", mock.Anything, mock.Anything).Return(anglerprofiles.Profile{}, errors.New("setup failed")).Once()

		session, _ := NewRunSession("Test Run", rng, ui, presenter)
		err := session.Run()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "setup failed")
	})

	t.Run("returns error when show intro fails", func(t *testing.T) {
		rng := &mockRNG{}
		ui := &mockRunUI{}
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())

		ui.On("ChooseAnglerProfile", mock.Anything, mock.Anything).Return(anglerprofiles.DefaultUnlockedProfiles()[0], nil).Once()

		stateProvider := &mockRunStateProvider{}
		stateProvider.On("Initialize", mock.Anything, mock.Anything, mock.Anything).Return(run.State{
			Status: run.StatusInProgress,
			Progress: run.ProgressState{
				Current: run.NodeState{NodeID: "start", Kind: run.NodeKindStart},
			},
		}, nil)

		renderer := &mockNodeRenderer{}
		renderer.On("ShowIntro", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("intro failed")).Once()

		session, _ := NewRunSession("Test Run", rng, ui, presenter)
		session.setLoopConfig(&RunLoopConfig{
			StateProvider: stateProvider,
			NodeRenderer:  renderer,
		})
		err := session.Run()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "show run intro")
	})

	t.Run("returns error when show node fails", func(t *testing.T) {
		rng := &mockRNG{}
		ui := &mockRunUI{}
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())

		ui.On("ChooseAnglerProfile", mock.Anything, mock.Anything).Return(anglerprofiles.DefaultUnlockedProfiles()[0], nil).Once()

		stateProvider := &mockRunStateProvider{}
		stateProvider.On("Initialize", mock.Anything, mock.Anything, mock.Anything).Return(run.State{
			Status: run.StatusInProgress,
			Progress: run.ProgressState{
				Current: run.NodeState{NodeID: "start", Kind: run.NodeKindStart},
			},
		}, nil)
		stateProvider.On("CurrentStatus", mock.Anything).Return(run.StatusInProgress).Once()
		stateProvider.On("CurrentNode", mock.Anything).Return(run.NodeState{NodeID: "start", Kind: run.NodeKindStart})

		renderer := &mockNodeRenderer{}
		renderer.On("ShowIntro", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		renderer.On("ShowNode", mock.Anything, mock.Anything).Return(errors.New("show node failed")).Once()

		session, _ := NewRunSession("Test Run", rng, ui, presenter)
		session.setLoopConfig(&RunLoopConfig{
			StateProvider: stateProvider,
			NodeRenderer:  renderer,
		})
		err := session.Run()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "show run node")
	})
}

func TestRunSessionRunReturnsErrorWhenShowIntroFails(t *testing.T) {
	rng := &mockRNG{}
	ui := &mockRunUI{}
	presenter := presentation.NewPresenter(presentation.DefaultCatalog())

	ui.On("ChooseAnglerProfile", mock.Anything, mock.Anything).Return(anglerprofiles.DefaultUnlockedProfiles()[0], nil).Once()

	stateProvider := &mockRunStateProvider{}
	stateProvider.On("Initialize", mock.Anything, mock.Anything, mock.Anything).Return(run.State{
		Status: run.StatusInProgress,
		Progress: run.ProgressState{
			Current: run.NodeState{NodeID: "start", Kind: run.NodeKindStart},
		},
	}, nil)

	renderer := &mockNodeRenderer{}
	renderer.On("ShowIntro", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("intro failed")).Once()

	session, _ := NewRunSession("Test Run", rng, ui, presenter)
	session.setLoopConfig(&RunLoopConfig{
		StateProvider: stateProvider,
		NodeRenderer:  renderer,
	})
	err := session.Run()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "show run intro")
}

func TestRunSessionRunReturnsErrorWhenShowNodeFails(t *testing.T) {
	rng := &mockRNG{}
	ui := &mockRunUI{}
	presenter := presentation.NewPresenter(presentation.DefaultCatalog())

	ui.On("ChooseAnglerProfile", mock.Anything, mock.Anything).Return(anglerprofiles.DefaultUnlockedProfiles()[0], nil).Once()

	stateProvider := &mockRunStateProvider{}
	stateProvider.On("Initialize", mock.Anything, mock.Anything, mock.Anything).Return(run.State{
		Status: run.StatusInProgress,
		Progress: run.ProgressState{
			Current: run.NodeState{NodeID: "start", Kind: run.NodeKindStart},
		},
	}, nil)
	stateProvider.On("CurrentStatus", mock.Anything).Return(run.StatusInProgress).Once()
	stateProvider.On("CurrentNode", mock.Anything).Return(run.NodeState{NodeID: "start", Kind: run.NodeKindStart})

	renderer := &mockNodeRenderer{}
	renderer.On("ShowIntro", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	renderer.On("ShowNode", mock.Anything, mock.Anything).Return(errors.New("show node failed")).Once()

	session, _ := NewRunSession("Test Run", rng, ui, presenter)
	session.setLoopConfig(&RunLoopConfig{
		StateProvider: stateProvider,
		NodeRenderer:  renderer,
	})
	err := session.Run()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "show run node")
}

type mockRNG struct {
	mock.Mock
}

func (rng *mockRNG) Float64() float64 {
	return rng.Called().Get(0).(float64)
}

func (rng *mockRNG) Intn(n int) int {
	return rng.Called(n).Int(0)
}

func (rng *mockRNG) Shuffle(n int, swap func(i, j int)) {
	rng.Called(n, swap)
}

type mockRunUI struct {
	mock.Mock
}

func (ui *mockRunUI) ChooseAnglerProfile(name string, profiles []anglerprofiles.Profile) (anglerprofiles.Profile, error) {
	args := ui.Called(name, profiles)
	return args.Get(0).(anglerprofiles.Profile), args.Error(1)
}

func (ui *mockRunUI) ChoosePlayerDeckPreset(name string, presets []playerprofiles.DeckPreset) (playerprofiles.DeckPreset, error) {
	args := ui.Called(name, presets)
	return args.Get(0).(playerprofiles.DeckPreset), args.Error(1)
}

func (ui *mockRunUI) ChooseRodPreset(name string, presets []rodpresets.Preset) (rodpresets.Preset, error) {
	args := ui.Called(name, presets)
	return args.Get(0).(rodpresets.Preset), args.Error(1)
}

func (ui *mockRunUI) ChooseAttachmentPreset(name string, baseRod rod.State, presets []attachmentpresets.Preset) (attachmentpresets.Preset, error) {
	args := ui.Called(name, baseRod, presets)
	return args.Get(0).(attachmentpresets.Preset), args.Error(1)
}

func (ui *mockRunUI) ShowOpening(view presentation.OpeningView) error {
	return ui.Called(view).Error(0)
}

func (ui *mockRunUI) ResolveCast(name string, context encounter.WaterContext, presenter CastPresenter) (encounter.CastResult, error) {
	args := ui.Called(name, context, presenter)
	return args.Get(0).(encounter.CastResult), args.Error(1)
}

func (ui *mockRunUI) ShowEncounterOpening(title string, opening presentation.OpeningView) error {
	return ui.Called(title, opening).Error(0)
}

func (ui *mockRunUI) ShowFishSpawn(title string, spawn presentation.SpawnView) error {
	return ui.Called(title, spawn).Error(0)
}

func (ui *mockRunUI) ChooseWaterContext(name string, contexts []watercontexts.Preset) (watercontexts.Preset, error) {
	args := ui.Called(name, contexts)
	return args.Get(0).(watercontexts.Preset), args.Error(1)
}

func (ui *mockRunUI) ChooseFishDeckPreset(name string, presets []fishprofiles.Profile) (fishprofiles.Profile, error) {
	args := ui.Called(name, presets)
	return args.Get(0).(fishprofiles.Profile), args.Error(1)
}

func (ui *mockRunUI) ShowRunIntro(view presentation.RunIntroView) error {
	return ui.Called(view).Error(0)
}

func (ui *mockRunUI) ShowRunNode(view presentation.RunNodeView) error {
	return ui.Called(view).Error(0)
}

func (ui *mockRunUI) ShowRunNodeSummary(view presentation.RunNodeSummaryView) error {
	return ui.Called(view).Error(0)
}

func (ui *mockRunUI) ShowRunSummary(view presentation.RunSummaryView) error {
	return ui.Called(view).Error(0)
}

func (ui *mockRunUI) ShowIntro(view presentation.IntroView) error {
	return ui.Called(view).Error(0)
}

func (ui *mockRunUI) ChooseMove(view presentation.StatusView, options []presentation.MoveOption) (domain.Move, error) {
	args := ui.Called(view, options)
	return args.Get(0).(domain.Move), args.Error(1)
}

func (ui *mockRunUI) ShowRound(view presentation.RoundView) error {
	return ui.Called(view).Error(0)
}

func (ui *mockRunUI) ResolveSplash(view presentation.SplashView) (encounter.SplashResolution, error) {
	args := ui.Called(view)
	return args.Get(0).(encounter.SplashResolution), args.Error(1)
}

func (ui *mockRunUI) ShowGameOver(view presentation.SummaryView) error {
	return ui.Called(view).Error(0)
}

type mockRunStateProvider struct {
	mock.Mock
}

func (m *mockRunStateProvider) Initialize(loadout loadout.State, route []run.NodeState, thread int) (run.State, error) {
	args := m.Called(loadout, route, thread)
	return args.Get(0).(run.State), args.Error(1)
}

func (m *mockRunStateProvider) Advance(state *run.State, route []run.NodeState) error {
	return m.Called(state, route).Error(0)
}

func (m *mockRunStateProvider) ApplyResult(state *run.State, result run.EncounterResult) error {
	return m.Called(state, result).Error(0)
}

func (m *mockRunStateProvider) Complete(state *run.State) error {
	return m.Called(state).Error(0)
}

func (m *mockRunStateProvider) CurrentStatus(state *run.State) run.Status {
	return m.Called(state).Get(0).(run.Status)
}

func (m *mockRunStateProvider) CurrentNode(state *run.State) run.NodeState {
	return m.Called(state).Get(0).(run.NodeState)
}

func (m *mockRunStateProvider) NextNode(state *run.State) *run.NodeState {
	args := m.Called(state)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*run.NodeState)
}

type mockEncounterHandler struct {
	mock.Mock
}

func (m *mockEncounterHandler) Resolve(
	title string,
	rng Randomizer,
	ui EncounterBootstrapUI,
	spawnUI SpawnUI,
	presenter presentation.Presenter,
	deckPreset playerprofiles.DeckPreset,
	playerLoadout loadout.State,
	config RunEncounterBootstrapConfig,
) (ResolvedEncounter, error) {
	args := m.Called(title, rng, ui, spawnUI, presenter, deckPreset, playerLoadout, config)
	return args.Get(0).(ResolvedEncounter), args.Error(1)
}

type mockNodeRenderer struct {
	mock.Mock
}

func (m *mockNodeRenderer) ShowIntro(title string, state run.State, route []run.NodeState) error {
	return m.Called(title, state, route).Error(0)
}

func (m *mockNodeRenderer) ShowNode(title string, state run.State) error {
	return m.Called(title, state).Error(0)
}

func (m *mockNodeRenderer) ShowNodeSummary(title string, currentNode run.NodeState, result run.EncounterResult, state run.State, nextNode *run.NodeState) error {
	return m.Called(title, currentNode, result, state, nextNode).Error(0)
}

func (m *mockNodeRenderer) ShowSummary(title string, state run.State) error {
	return m.Called(title, state).Error(0)
}
