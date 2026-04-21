package presentation

import (
	"testing"

	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/game"
)

func TestDefaultCatalogBuildsExpectedMoveOptions(t *testing.T) {
	presenter := NewPresenter(DefaultCatalog())
	intro := presenter.Intro()

	if intro.Title != "Pesca: duelo contra el pez" {
		t.Fatalf("title = %q, want %q", intro.Title, "Pesca: duelo contra el pez")
	}
	if len(intro.Options) != 3 {
		t.Fatalf("options = %d, want 3", len(intro.Options))
	}
	if intro.Options[0].Label != "Tirar" || intro.Options[1].Label != "Recoger" || intro.Options[2].Label != "Soltar" {
		t.Fatalf("unexpected option labels: %#v", intro.Options)
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
	if err != nil {
		t.Fatalf("NewState() error = %v", err)
	}
	encounterState.Status = encounter.StatusCaptured
	encounterState.EndReason = encounter.EndReasonDeckCapture

	round := presenter.Round(game.RoundResult{
		PlayerMove: domain.Blue,
		FishMove:   domain.Red,
		Outcome:    domain.PlayerWin,
		State:      game.State{Encounter: encounterState},
	})
	if round.PlayerLabel != "Lanzar" || round.FishLabel != "Afianzar" || round.Outcome != "aventaja el jugador" {
		t.Fatalf("unexpected round view: %#v", round)
	}

	summary := presenter.Summary(game.State{Encounter: encounterState})
	if summary.Outcome != "presa asegurada" || summary.EndReason != "sin mazo, pesca cerrada" {
		t.Fatalf("unexpected summary view: %#v", summary)
	}
}
