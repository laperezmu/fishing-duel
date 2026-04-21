package app_test

import (
	"testing"

	"pesca/internal/app"
	"pesca/internal/deck"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/endings"
	"pesca/internal/game"
	"pesca/internal/presentation"
	"pesca/internal/progression"
	"pesca/internal/rules"
)

func TestSessionRunsThroughAbstractUI(t *testing.T) {
	encounterState, err := encounter.NewState(encounter.Config{
		InitialDistance:           3,
		CaptureDistance:           -1,
		EscapeDistance:            99,
		ExhaustionCaptureDistance: 2,
		PlayerWinStep:             1,
		FishWinStep:               1,
	})
	if err != nil {
		t.Fatalf("NewState() error = %v", err)
	}

	manager := deck.NewManager(
		deck.NewStandardFishDeck(),
		func([]domain.Move) {},
		deck.RemoveCardsRecyclePolicy{CardsToRemove: 3},
	)

	engine, err := game.NewEngine(
		manager,
		rules.ClassicEvaluator{},
		progression.TrackPolicy{},
		endings.EncounterCondition{},
		game.State{Encounter: encounterState},
	)
	if err != nil {
		t.Fatalf("NewEngine() error = %v", err)
	}

	ui := &stubUI{moves: repeatMove(domain.Blue, 18)}
	presenter := presentation.NewPresenter(presentation.DefaultCatalog())
	session, err := app.NewSession(engine, ui, presenter)
	if err != nil {
		t.Fatalf("NewSession() error = %v", err)
	}

	if err := session.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if ui.intro.Title == "" {
		t.Fatalf("intro was not shown")
	}
	if len(ui.rounds) != 18 {
		t.Fatalf("rounds shown = %d, want 18", len(ui.rounds))
	}
	if ui.summary.TotalRounds != 18 {
		t.Fatalf("summary rounds = %d, want 18", ui.summary.TotalRounds)
	}
	if ui.summary.Outcome != "pez capturado" {
		t.Fatalf("summary outcome = %q, want pez capturado", ui.summary.Outcome)
	}
	if ui.summary.EndReason != "captura al agotar la baraja con distancia 2 o menor" {
		t.Fatalf("summary end reason = %q", ui.summary.EndReason)
	}
	if ui.movesConsumed != 18 {
		t.Fatalf("moves consumed = %d, want 18", ui.movesConsumed)
	}
}

type stubUI struct {
	intro         presentation.IntroView
	rounds        []presentation.RoundView
	summary       presentation.SummaryView
	moves         []domain.Move
	movesConsumed int
}

func (ui *stubUI) ShowIntro(view presentation.IntroView) error {
	ui.intro = view
	return nil
}

func (ui *stubUI) ChooseMove(_ presentation.StatusView, _ []presentation.MoveOption) (domain.Move, error) {
	move := ui.moves[ui.movesConsumed]
	ui.movesConsumed++
	return move, nil
}

func (ui *stubUI) ShowRound(view presentation.RoundView) error {
	ui.rounds = append(ui.rounds, view)
	return nil
}

func (ui *stubUI) ShowGameOver(view presentation.SummaryView) error {
	ui.summary = view
	return nil
}

func repeatMove(move domain.Move, count int) []domain.Move {
	moves := make([]domain.Move, count)
	for i := range moves {
		moves[i] = move
	}
	return moves
}
