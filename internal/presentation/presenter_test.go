package presentation

import (
	"pesca/internal/cards"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/match"
	"pesca/internal/player/loadout"
	"pesca/internal/player/rod"
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
		Round: match.RoundState{Number: 2},
		Deck: match.DeckState{
			ActiveCards:       4,
			DiscardCards:      5,
			RecycleCount:      1,
			ShufflesOnRecycle: true,
			CardsToRemove:     3,
			CurrentCycle: match.FishDiscardCycleState{
				Number:     2,
				TotalCards: 3,
				Entries: []match.FishDiscardEntryState{
					{Visibility: cards.DiscardVisibilityFull, Move: domain.Blue, Name: "Oleaje abierto"},
					{Visibility: cards.DiscardVisibilityMasked},
				},
			},
			PreviousCycleStats: []match.FishDiscardCycleSummaryState{{
				Number:       1,
				TotalCards:   4,
				VisibleCards: 3,
				HiddenCards:  1,
			}},
		},
		Encounter: encounter.State{
			Config: encounter.Config{
				InitialDepth:              1,
				SurfaceDepth:              0,
				CaptureDistance:           0,
				ExhaustionCaptureDistance: 2,
				SplashProfile:             encounter.DefaultSplashProfile(),
			},
			Distance: 3,
			Depth:    2,
		},
		Player: match.PlayerState{
			Loadout: mustLoadout(t, rod.State{OpeningMaxDistance: 4, OpeningMaxDepth: 3, TrackMaxDistance: 5, TrackMaxDepth: 4}),
			Moves: match.PlayerMoveResources{Moves: []match.PlayerMoveState{
				{Move: domain.Blue, MaxUses: 3, RemainingUses: 2, ActiveCards: []cards.PlayerCard{cards.NewNamedPlayerCard("Anzuelo tenso", "Capturas desde un paso mas lejos este round.", domain.Blue, cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1}), cards.NewPlayerCard(domain.Blue)}},
				{Move: domain.Red, MaxUses: 3, RemainingUses: 0, RestoresOnRound: 5},
				{Move: domain.Yellow, MaxUses: 3, RemainingUses: 1},
			}},
		},
		Lifecycle: match.LifecycleState{Stats: match.Stats{PlayerWins: 2, FishWins: 1, Draws: 3}},
	}

	status := presenter.Status(match.NewStatusSnapshot(state))

	assert.Equal(t, 3, status.RoundNumber)
	assert.Equal(t, 3, status.FishDistance)
	assert.Equal(t, 2, status.FishDepth)
	assert.Equal(t, 0, status.SurfaceDepth)
	assert.Equal(t, 5, status.MaxDistance)
	assert.Equal(t, 4, status.MaxDepth)
	assert.Equal(t, 0, status.CaptureDistance)
	assert.Equal(t, 2, status.ExhaustionCaptureDistance)
	assert.Equal(t, 4, status.ActiveCards)
	assert.Equal(t, 5, status.DiscardCards)
	assert.Equal(t, 1, status.RecycleCount)
	assert.Equal(t, 2, status.PlayerWins)
	assert.Equal(t, 1, status.FishWins)
	assert.Equal(t, 3, status.Draws)
	assert.True(t, status.FishDiscard.ShufflesOnRecycle)
	assert.Equal(t, 3, status.FishDiscard.CardsToRemove)
	assert.Equal(t, 3, status.FishDiscard.CurrentCycleTotalCards)
	assert.Equal(t, "Oleaje abierto", status.FishDiscard.CurrentCycleEntries[0].Label)
	assert.Equal(t, "?", status.FishDiscard.CurrentCycleEntries[1].Label)
	require.Len(t, status.FishDiscard.PreviousCycles, 1)
	assert.Equal(t, 1, status.FishDiscard.PreviousCycles[0].HiddenCards)
	require.Len(t, status.MoveOptions, 3)
	assert.Equal(t, 2, status.MoveOptions[0].RemainingUses)
	assert.Equal(t, "Anzuelo tenso", status.MoveOptions[0].CardHint)
	assert.True(t, status.MoveOptions[0].Available)
	assert.Equal(t, 5, status.MoveOptions[1].RestoresOnRound)
	assert.False(t, status.MoveOptions[1].Available)
}

func TestPresenterRound(t *testing.T) {
	presenter := NewPresenter(newCustomCatalog())
	encounterState := newCapturedEncounterState(t)
	encounterState.LastEvent = encounter.Event{Kind: encounter.EventKindSplash, Escaped: false}
	round := presenter.Round(match.NewRoundSnapshot(match.RoundResult{
		PlayerMove: domain.Blue,
		FishMove:   domain.Red,
		Outcome:    domain.PlayerWin,
		Status:     match.NewStatusSnapshot(match.State{Encounter: encounterState}),
		Encounter:  match.EncounterEventSnapshot{LastEvent: encounterState.LastEvent},
	}))

	assert.Equal(t, "Lanzar", round.PlayerLabel)
	assert.Equal(t, "Afianzar", round.FishLabel)
	assert.Equal(t, domain.PlayerWin, round.Outcome)
	assert.Equal(t, "aventaja el jugador", round.OutcomeLabel)
	assert.Equal(t, "chapotea: sigue enganchado", round.EventLabel)
}

func TestPresenterSummary(t *testing.T) {
	presenter := NewPresenter(newCustomCatalog())
	summary := presenter.Summary(match.NewSummarySnapshot(match.State{Encounter: newCapturedEncounterState(t)}))

	assert.Equal(t, encounter.StatusCaptured, summary.EncounterStatus)
	assert.Equal(t, 1, summary.FishDepth)
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
		EncounterEvents: map[encounter.EventKind]string{
			encounter.EventKindSplash: "chapotea",
		},
		EventOutcomes: map[bool]string{
			false: "sigue enganchado",
			true:  "se suelta",
		},
		EncounterResults: map[encounter.Status]string{
			encounter.StatusCaptured: "presa asegurada",
		},
		EndReasons: map[encounter.EndReason]string{
			encounter.EndReasonDeckCapture:  "sin mazo, pesca cerrada",
			encounter.EndReasonSplashEscape: "escape por chapoteo",
		},
	}
}

func mustLoadout(t *testing.T, playerRod rod.State) loadout.State {
	t.Helper()

	playerLoadout, err := loadout.NewState(playerRod, nil)
	require.NoError(t, err)

	return playerLoadout
}

func newCapturedEncounterState(t *testing.T) encounter.State {
	t.Helper()

	state, err := encounter.NewState(encounter.DefaultConfig())
	require.NoError(t, err)
	state.Status = encounter.StatusCaptured
	state.EndReason = encounter.EndReasonDeckCapture

	return state
}
