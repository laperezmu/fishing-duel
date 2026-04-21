package game_test

import (
	"testing"

	"pesca/internal/deck"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/endings"
	"pesca/internal/game"
	"pesca/internal/progression"
	"pesca/internal/rules"
)

func TestEngineCapturesWhenTrackReachesPlayer(t *testing.T) {
	engine := newEngineForTest(t, []domain.Move{domain.Red, domain.Red, domain.Red}, encounter.DefaultConfig())

	for i := 0; i < 3; i++ {
		result, err := engine.PlayRound(domain.Blue)
		if err != nil {
			t.Fatalf("PlayRound() round %d error = %v", i+1, err)
		}
		if i < 2 && result.State.Finished {
			t.Fatalf("game finished too early on round %d", i+1)
		}
	}

	state := engine.State()
	if !state.Finished {
		t.Fatalf("finished = false, want true")
	}
	if state.Encounter.Status != encounter.StatusCaptured {
		t.Fatalf("status = %q, want %q", state.Encounter.Status, encounter.StatusCaptured)
	}
	if state.Encounter.EndReason != encounter.EndReasonTrackCapture {
		t.Fatalf("end reason = %q, want %q", state.Encounter.EndReason, encounter.EndReasonTrackCapture)
	}
	if state.Encounter.Distance != 0 {
		t.Fatalf("distance = %d, want 0", state.Encounter.Distance)
	}
	if state.Stats.PlayerWins != 3 {
		t.Fatalf("player wins = %d, want 3", state.Stats.PlayerWins)
	}
}

func TestEngineEscapesWhenTrackExceedsLimit(t *testing.T) {
	engine := newEngineForTest(t, []domain.Move{domain.Yellow, domain.Yellow, domain.Yellow}, encounter.DefaultConfig())

	for i := 0; i < 3; i++ {
		if _, err := engine.PlayRound(domain.Blue); err != nil {
			t.Fatalf("PlayRound() round %d error = %v", i+1, err)
		}
	}

	state := engine.State()
	if state.Encounter.Status != encounter.StatusEscaped {
		t.Fatalf("status = %q, want %q", state.Encounter.Status, encounter.StatusEscaped)
	}
	if state.Encounter.EndReason != encounter.EndReasonTrackEscape {
		t.Fatalf("end reason = %q, want %q", state.Encounter.EndReason, encounter.EndReasonTrackEscape)
	}
	if state.Encounter.Distance != 6 {
		t.Fatalf("distance = %d, want 6", state.Encounter.Distance)
	}
}

func TestEngineCapturesOnDeckExhaustionNearPlayer(t *testing.T) {
	engine := newEngineForTest(t, []domain.Move{domain.Red}, encounter.DefaultConfig())

	if _, err := engine.PlayRound(domain.Blue); err != nil {
		t.Fatalf("PlayRound() error = %v", err)
	}

	state := engine.State()
	if state.Encounter.Status != encounter.StatusCaptured {
		t.Fatalf("status = %q, want %q", state.Encounter.Status, encounter.StatusCaptured)
	}
	if state.Encounter.EndReason != encounter.EndReasonDeckCapture {
		t.Fatalf("end reason = %q, want %q", state.Encounter.EndReason, encounter.EndReasonDeckCapture)
	}
	if state.Encounter.Distance != 2 {
		t.Fatalf("distance = %d, want 2", state.Encounter.Distance)
	}
}

func TestEngineEscapesOnDeckExhaustionFarFromPlayer(t *testing.T) {
	engine := newEngineForTest(t, []domain.Move{domain.Yellow}, encounter.Config{
		InitialDistance:           3,
		CaptureDistance:           0,
		EscapeDistance:            10,
		ExhaustionCaptureDistance: 2,
		PlayerWinStep:             1,
		FishWinStep:               1,
	})

	if _, err := engine.PlayRound(domain.Blue); err != nil {
		t.Fatalf("PlayRound() error = %v", err)
	}

	state := engine.State()
	if state.Encounter.Status != encounter.StatusEscaped {
		t.Fatalf("status = %q, want %q", state.Encounter.Status, encounter.StatusEscaped)
	}
	if state.Encounter.EndReason != encounter.EndReasonDeckEscape {
		t.Fatalf("end reason = %q, want %q", state.Encounter.EndReason, encounter.EndReasonDeckEscape)
	}
	if state.Encounter.Distance != 4 {
		t.Fatalf("distance = %d, want 4", state.Encounter.Distance)
	}
	if _, err := engine.PlayRound(domain.Blue); err != game.ErrGameFinished {
		t.Fatalf("PlayRound() after finish error = %v, want %v", err, game.ErrGameFinished)
	}
}

func newEngineForTest(t *testing.T, cards []domain.Move, config encounter.Config) *game.Engine {
	t.Helper()

	encounterState, err := encounter.NewState(config)
	if err != nil {
		t.Fatalf("NewState() error = %v", err)
	}

	engine, err := game.NewEngine(
		deck.NewManager(cards, func([]domain.Move) {}, deck.RemoveCardsRecyclePolicy{CardsToRemove: 3}),
		rules.ClassicEvaluator{},
		progression.TrackPolicy{},
		endings.EncounterCondition{},
		game.State{Encounter: encounterState},
	)
	if err != nil {
		t.Fatalf("NewEngine() error = %v", err)
	}

	return engine
}
