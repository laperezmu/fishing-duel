package presentation

import (
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/match"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPresenterIntro(t *testing.T) {
	tests := []struct {
		title      string
		catalog    Catalog
		wantTitle  string
		wantLabels []string
	}{
		{
			title:      "returns default title and move labels when using the default catalog",
			catalog:    DefaultCatalog(),
			wantTitle:  "Pesca: duelo contra el pez",
			wantLabels: []string{"Tirar", "Recoger", "Soltar"},
		},
		{
			title: "returns custom move labels when the catalog overrides player text",
			catalog: Catalog{
				Title: "Custom",
				PlayerMoveLabels: map[domain.Move]string{
					domain.Blue:   "Lanzar",
					domain.Red:    "Cobrar",
					domain.Yellow: "Liberar",
				},
			},
			wantTitle:  "Custom",
			wantLabels: []string{"Lanzar", "Cobrar", "Liberar"},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			intro := NewPresenter(test.catalog).Intro()

			require.Len(t, intro.Options, 3)
			assert.Equal(t, test.wantTitle, intro.Title)
			assert.Equal(t, test.wantLabels[0], intro.Options[0].Label)
			assert.Equal(t, test.wantLabels[1], intro.Options[1].Label)
			assert.Equal(t, test.wantLabels[2], intro.Options[2].Label)
		})
	}
}

func TestPresenterStatus(t *testing.T) {
	presenter := NewPresenter(DefaultCatalog())
	state := match.State{
		Round: 2,
		Deck: match.DeckState{
			ActiveCards:  4,
			DiscardCards: 5,
			RecycleCount: 1,
		},
		Encounter: encounter.State{
			Config: encounter.Config{
				CaptureDistance:           0,
				EscapeDistance:            5,
				ExhaustionCaptureDistance: 2,
			},
			Distance: 3,
		},
		Stats: match.Stats{
			PlayerWins: 2,
			FishWins:   1,
			Draws:      3,
		},
	}

	status := presenter.Status(state)

	assert.Equal(t, 3, status.RoundNumber)
	assert.Equal(t, 3, status.FishDistance)
	assert.Equal(t, 0, status.CaptureDistance)
	assert.Equal(t, 5, status.EscapeDistance)
	assert.Equal(t, 2, status.ExhaustionCaptureDistance)
	assert.Equal(t, 4, status.ActiveCards)
	assert.Equal(t, 5, status.DiscardCards)
	assert.Equal(t, 1, status.RecycleCount)
	assert.Equal(t, 2, status.PlayerWins)
	assert.Equal(t, 1, status.FishWins)
	assert.Equal(t, 3, status.Draws)
}

func TestPresenterRound(t *testing.T) {
	presenter := NewPresenter(newCustomCatalog())
	round := presenter.Round(match.RoundResult{
		PlayerMove: domain.Blue,
		FishMove:   domain.Red,
		Outcome:    domain.PlayerWin,
		State:      match.State{Encounter: newCapturedEncounterState(t)},
	})

	assert.Equal(t, "Lanzar", round.PlayerLabel)
	assert.Equal(t, "Afianzar", round.FishLabel)
	assert.Equal(t, domain.PlayerWin, round.Outcome)
	assert.Equal(t, "aventaja el jugador", round.OutcomeLabel)
}

func TestPresenterSummary(t *testing.T) {
	presenter := NewPresenter(newCustomCatalog())
	summary := presenter.Summary(match.State{Encounter: newCapturedEncounterState(t)})

	assert.Equal(t, encounter.StatusCaptured, summary.EncounterStatus)
	assert.Equal(t, "presa asegurada", summary.OutcomeLabel)
	assert.Equal(t, "sin mazo, pesca cerrada", summary.EndReasonLabel)
}

func newCustomCatalog() Catalog {
	return Catalog{
		Title: "Custom",
		PlayerMoveLabels: map[domain.Move]string{
			domain.Blue:   "Lanzar",
			domain.Red:    "Cobrar",
			domain.Yellow: "Liberar",
		},
		FishMoveLabels: map[domain.Move]string{
			domain.Blue:   "Golpear",
			domain.Red:    "Afianzar",
			domain.Yellow: "Huir",
		},
		RoundOutcomes: map[domain.RoundOutcome]string{
			domain.PlayerWin: "aventaja el jugador",
		},
		EncounterResults: map[encounter.Status]string{
			encounter.StatusCaptured: "presa asegurada",
		},
		EndReasons: map[encounter.EndReason]string{
			encounter.EndReasonDeckCapture: "sin mazo, pesca cerrada",
		},
	}
}

func newCapturedEncounterState(t *testing.T) encounter.State {
	t.Helper()

	state, err := encounter.NewState(encounter.DefaultConfig())
	require.NoError(t, err)
	state.Status = encounter.StatusCaptured
	state.EndReason = encounter.EndReasonDeckCapture

	return state
}
