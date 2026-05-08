package presentation_test

import (
	"testing"

	"pesca/internal/content/anglerprofiles"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/waterpools"
	"pesca/internal/encounter"
	"pesca/internal/match"
	"pesca/internal/presentation"
	"pesca/internal/run"

	"github.com/stretchr/testify/assert"
)

func TestPresenterOpening(t *testing.T) {
	t.Run("renders opening with water context and cast information", func(t *testing.T) {
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())
		opening := encounter.Opening{
			WaterContext: encounter.WaterContext{
				Name:    "Costa cercana",
				PoolTag: waterpools.Shoreline,
			},
			InitialDistance: 4,
			InitialDepth:    1,
		}

		view := presenter.Opening(opening)

		assert.Equal(t, "Costa cercana", view.WaterLabel)
		assert.Equal(t, 4, view.InitialDistance)
		assert.Equal(t, 1, view.InitialDepth)
	})
}

func TestPresenterSpawn(t *testing.T) {
	t.Run("renders spawn with profile and context information", func(t *testing.T) {
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())
		spawn := fishprofiles.Spawn{
			Profile: fishprofiles.Profile{
				Name: "Bass inicial",
			},
			Context: fishprofiles.SpawnContext{
				WaterPoolTag:    waterpools.Shoreline,
				InitialDistance: 3,
				InitialDepth:    2,
			},
			CandidateCount: 5,
		}

		view := presenter.Spawn(spawn)

		assert.Equal(t, "Bass inicial", view.ProfileLabel)
		assert.Equal(t, 3, view.InitialDistance)
		assert.Equal(t, 2, view.InitialDepth)
		assert.Equal(t, 5, view.CandidateCount)
	})
}

func TestPresenterCast(t *testing.T) {
	t.Run("renders cast view", func(t *testing.T) {
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())
		context := encounter.WaterContext{
			Name:    "Costa",
			PoolTag: waterpools.Shoreline,
		}

		castView := presenter.Cast(context, 2, 5, 80)

		assert.Equal(t, "Costa", castView.WaterLabel)
	})
}

func TestPresenterSplash(t *testing.T) {
	t.Run("renders splash event", func(t *testing.T) {
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())
		snapshot := match.EncounterEventSnapshot{
			LastEvent: encounter.Event{
				Kind: encounter.EventKindSplash,
			},
		}

		view := presenter.Splash(snapshot, 2)

		assert.Contains(t, view.EventLabel, "chapotea")
	})
}

func TestPresenterAnglerProfile(t *testing.T) {
	t.Run("renders angler profile view", func(t *testing.T) {
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())
		profile := anglerprofiles.Profile{
			ID:   "coastal-specialist",
			Name: "Especialista costero",
		}

		view := presenter.AnglerProfile(profile)

		assert.Equal(t, "coastal-specialist", view.ProfileID)
		assert.Equal(t, "Especialista costero", view.Name)
	})
}

func TestPresenterRunIntro(t *testing.T) {
	t.Run("renders run intro with state and route", func(t *testing.T) {
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())
		state := run.State{}
		route := []run.NodeState{}

		view := presenter.RunIntro("Test Run", state, route)

		assert.Equal(t, "Test Run", view.Title)
	})

	t.Run("renders run intro with route nodes", func(t *testing.T) {
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())
		state := run.State{Thread: run.ThreadState{Current: 5}}
		route := []run.NodeState{
			{NodeID: "start", Kind: run.NodeKindStart, ZoneID: "coast"},
			{NodeID: "fishing-1", Kind: run.NodeKindFishing, ZoneID: "coast"},
		}

		view := presenter.RunIntro("Test Run", state, route)

		assert.Equal(t, "Test Run", view.Title)
		assert.Equal(t, 5, view.Thread)
		assert.Len(t, view.RouteLabels, 2)
	})

	t.Run("renders run intro with thread state", func(t *testing.T) {
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())
		state := run.State{Thread: run.ThreadState{Current: 7, Maximum: 10}}
		route := []run.NodeState{}

		view := presenter.RunIntro("Test Run", state, route)

		assert.Equal(t, 7, view.Thread)
	})
}

func TestPresenterRunNode(t *testing.T) {
	t.Run("renders run node view", func(t *testing.T) {
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())
		state := run.State{}

		view := presenter.RunNode("Test Run", state)

		assert.Equal(t, "Test Run", view.Title)
	})

	t.Run("renders run node with progress state", func(t *testing.T) {
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())
		state := run.State{
			Progress: run.ProgressState{
				Current: run.NodeState{NodeID: "fishing-1", Kind: run.NodeKindFishing, ZoneID: "coast"},
			},
			Thread: run.ThreadState{Current: 5, Maximum: 10},
		}

		view := presenter.RunNode("Test Run", state)

		assert.Equal(t, "Test Run", view.Title)
		assert.Equal(t, run.NodeKindFishing, view.NodeKind)
		assert.Equal(t, 5, view.Thread)
		assert.Equal(t, 10, view.ThreadMax)
		assert.Equal(t, 0, view.CaptureCount)
	})

	t.Run("renders run node with captures", func(t *testing.T) {
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())
		state := run.State{
			Progress: run.ProgressState{
				Current: run.NodeState{NodeID: "fishing-1", Kind: run.NodeKindFishing, ZoneID: "coast"},
			},
			Captures: []run.CaptureRecord{
				{FishID: "bass-1", FishName: "Lubina"},
				{FishID: "mackerel-1", FishName: "Caballa"},
			},
		}

		view := presenter.RunNode("Test Run", state)

		assert.Equal(t, 2, view.CaptureCount)
	})
}

func TestPresenterRunSummary(t *testing.T) {
	t.Run("renders run summary with victory status", func(t *testing.T) {
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())
		state := run.State{}

		view := presenter.RunSummary("Test Run", state)

		assert.Equal(t, "Test Run", view.Title)
	})
}

func TestPresenterRunNodeSummary(t *testing.T) {
	t.Run("renders run node summary", func(t *testing.T) {
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())
		currentNode := run.NodeState{NodeID: "fishing-1"}
		result := run.EncounterResult{Outcome: run.EncounterOutcomeCaptured}
		state := run.State{}
		nextNode := &run.NodeState{NodeID: "fishing-2"}

		view := presenter.RunNodeSummary("Test Run", currentNode, result, state, nextNode)

		assert.Equal(t, "Test Run", view.Title)
	})

	t.Run("renders escaped outcome", func(t *testing.T) {
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())
		currentNode := run.NodeState{NodeID: "fishing-1"}
		result := run.EncounterResult{Outcome: run.EncounterOutcomeEscaped}
		state := run.State{}
		nextNode := &run.NodeState{NodeID: "fishing-2"}

		view := presenter.RunNodeSummary("Test Run", currentNode, result, state, nextNode)

		assert.Equal(t, "Test Run", view.Title)
	})

	t.Run("renders defeated outcome", func(t *testing.T) {
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())
		currentNode := run.NodeState{NodeID: "boss-1"}
		result := run.EncounterResult{Outcome: run.EncounterOutcomeDefeated}
		state := run.State{}
		nextNode := &run.NodeState{NodeID: "fishing-2"}

		view := presenter.RunNodeSummary("Test Run", currentNode, result, state, nextNode)

		assert.Equal(t, "Test Run", view.Title)
	})

	t.Run("renders escaped with thread damage", func(t *testing.T) {
		presenter := presentation.NewPresenter(presentation.DefaultCatalog())
		currentNode := run.NodeState{NodeID: "fishing-1"}
		result := run.EncounterResult{Outcome: run.EncounterOutcomeEscaped, ThreadDamage: 3, Capture: nil}
		state := run.State{}
		nextNode := &run.NodeState{NodeID: "fishing-2"}

		view := presenter.RunNodeSummary("Test Run", currentNode, result, state, nextNode)

		assert.Equal(t, "Test Run", view.Title)
	})
}
