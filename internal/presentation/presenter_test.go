package presentation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/game"
)

func TestPresenterIntroBuildsExpectedMoveOptions(t *testing.T) {
	tests := []struct {
		name       string
		catalog    Catalog
		wantTitle  string
		wantLabels []string
	}{
		{
			name:      "default catalog",
			catalog:   DefaultCatalog(),
			wantTitle: "Pesca: duelo contra el pez",
			wantLabels: []string{
				"Tirar",
				"Recoger",
				"Soltar",
			},
		},
		{
			name: "custom player labels",
			catalog: Catalog{
				Title: "Custom",
				PlayerMoveLabels: map[domain.Move]string{
					domain.Blue:   "Lanzar",
					domain.Red:    "Cobrar",
					domain.Yellow: "Liberar",
				},
			},
			wantTitle: "Custom",
			wantLabels: []string{
				"Lanzar",
				"Cobrar",
				"Liberar",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			presenter := NewPresenter(test.catalog)
			intro := presenter.Intro()

			require.Len(t, intro.Options, 3)
			assert.Equal(t, test.wantTitle, intro.Title)
			assert.Equal(t, test.wantLabels[0], intro.Options[0].Label)
			assert.Equal(t, test.wantLabels[1], intro.Options[1].Label)
			assert.Equal(t, test.wantLabels[2], intro.Options[2].Label)
		})
	}
}

func TestPresenterUsesCustomCatalogTexts(t *testing.T) {
	presenter := NewPresenter(Catalog{
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
	})

	encounterState, err := encounter.NewState(encounter.DefaultConfig())
	require.NoError(t, err)
	encounterState.Status = encounter.StatusCaptured
	encounterState.EndReason = encounter.EndReasonDeckCapture

	round := presenter.Round(game.RoundResult{
		PlayerMove: domain.Blue,
		FishMove:   domain.Red,
		Outcome:    domain.PlayerWin,
		State:      game.State{Encounter: encounterState},
	})
	assert.Equal(t, "Lanzar", round.PlayerLabel)
	assert.Equal(t, "Afianzar", round.FishLabel)
	assert.Equal(t, "aventaja el jugador", round.Outcome)

	summary := presenter.Summary(game.State{Encounter: encounterState})
	assert.Equal(t, "presa asegurada", summary.Outcome)
	assert.Equal(t, "sin mazo, pesca cerrada", summary.EndReason)
}
