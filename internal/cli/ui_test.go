package cli

import (
	"bytes"
	"strings"
	"testing"

	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/game"
	"pesca/internal/presentation"
)

func TestShowIntroIncludesColoredOptions(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(strings.NewReader("1\n"), &out)
	presenter := presentation.NewPresenter(presentation.DefaultCatalog())

	err := ui.ShowIntro(presenter.Intro())
	if err != nil {
		t.Fatalf("ShowIntro() error = %v", err)
	}

	move, err := ui.ChooseMove(presenter.Status(dummyState()), presenter.Intro().Options)
	if err != nil {
		t.Fatalf("ChooseMove() error = %v", err)
	}
	if move != domain.Blue {
		t.Fatalf("move = %v, want %v", move, domain.Blue)
	}

	printed := out.String()
	assertContains(t, printed, clearSequence)
	assertContains(t, printed, "Tensa el sedal y arrastra al pez hacia la orilla.")
	assertContains(t, printed, "Orilla")
	assertContains(t, printed, "[ESC]")
	assertContains(t, printed, colorizeMove(domain.Blue, "Tirar"))
	assertContains(t, printed, colorizeMove(domain.Red, "Recoger"))
	assertContains(t, printed, colorizeMove(domain.Yellow, "Soltar"))
}

func TestChooseMoveShowsLastRoundSummary(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(strings.NewReader("2\n"), &out)
	presenter := presentation.NewPresenter(presentation.DefaultCatalog())
	if err := ui.ShowIntro(presenter.Intro()); err != nil {
		t.Fatalf("ShowIntro() error = %v", err)
	}

	err := ui.ShowRound(presentation.RoundView{
		Status:      presentation.StatusView{Distance: 2},
		PlayerMove:  domain.Blue,
		FishMove:    domain.Yellow,
		PlayerLabel: "Tirar",
		FishLabel:   "Zafarse",
		Outcome:     "gana el jugador",
	})
	if err != nil {
		t.Fatalf("ShowRound() error = %v", err)
	}

	_, err = ui.ChooseMove(presentation.StatusView{
		RoundNumber:               2,
		Distance:                  2,
		CaptureDistance:           0,
		EscapeDistance:            5,
		ExhaustionCaptureDistance: 2,
		ActiveCards:               8,
		DiscardCards:              1,
		RecycleCount:              0,
		PlayerWins:                1,
		FishWins:                  0,
		Draws:                     0,
	}, presenter.Intro().Options)
	if err != nil {
		t.Fatalf("ChooseMove() error = %v", err)
	}

	printed := out.String()
	assertContains(t, printed, "Ultimo lance")
	assertContains(t, printed, colorizeMove(domain.Blue, "Tirar"))
	assertContains(t, printed, colorizeMove(domain.Yellow, "Zafarse"))
	assertContains(t, printed, outcomeColor("gana el jugador"))
	assertContains(t, printed, "Distancia : 2")
}

func dummyState() game.State {
	encounterState, _ := encounter.NewState(encounter.DefaultConfig())
	return game.State{
		Deck: game.DeckState{
			ActiveCards: 9,
		},
		Encounter: encounterState,
	}
}

func assertContains(t *testing.T, text, want string) {
	t.Helper()
	if !strings.Contains(text, want) {
		t.Fatalf("output %q does not contain %q", text, want)
	}
}
