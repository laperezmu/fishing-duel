package presentation

import (
	"pesca/internal/cards"
	"pesca/internal/content/attachmentpresets"
	"pesca/internal/content/playerprofiles"
	"pesca/internal/content/rodpresets"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/match"
	"pesca/internal/player/loadout"
	"pesca/internal/player/rod"
	"pesca/internal/run"
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
		ResolvedEffects: []match.ResolvedEffectState{{
			Owner:    cards.OwnerFish,
			Type:     cards.EffectTypeLegacyCaptureWindow,
			Priority: 60,
		}},
		Status:    match.NewStatusSnapshot(match.State{Encounter: encounterState}),
		Encounter: match.EncounterEventSnapshot{LastEvent: encounterState.LastEvent},
	}))

	assert.Equal(t, "Lanzar", round.PlayerLabel)
	assert.Equal(t, "Afianzar", round.FishLabel)
	assert.Equal(t, domain.PlayerWin, round.Outcome)
	assert.Equal(t, "aventaja el jugador", round.OutcomeLabel)
	assert.Equal(t, "chapotea: sigue enganchado", round.EventLabel)
	require.Len(t, round.Resolved, 1)
	assert.Equal(t, "pez | ventana captura | p60", round.Resolved[0])
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

func TestPresenterRoundWithDifferentOutcomes(t *testing.T) {
	presenter := NewPresenter(DefaultCatalog())

	t.Run("renders fish win outcome", func(t *testing.T) {
		round := presenter.Round(match.NewRoundSnapshot(match.RoundResult{
			PlayerMove: domain.Blue,
			FishMove:   domain.Red,
			Outcome:    domain.FishWin,
			Status:     match.NewStatusSnapshot(match.State{}),
		}))

		assert.Equal(t, "gana el pez", round.OutcomeLabel)
	})

	t.Run("renders draw outcome", func(t *testing.T) {
		round := presenter.Round(match.NewRoundSnapshot(match.RoundResult{
			PlayerMove: domain.Blue,
			FishMove:   domain.Red,
			Outcome:    domain.Draw,
			Status:     match.NewStatusSnapshot(match.State{}),
		}))

		assert.Equal(t, "empate", round.OutcomeLabel)
	})
}

func TestPresenterSummaryWithEscaped(t *testing.T) {
	presenter := NewPresenter(DefaultCatalog())
	state, err := encounter.NewState(encounter.DefaultConfig())
	require.NoError(t, err)
	state.Status = encounter.StatusEscaped
	state.EndReason = encounter.EndReasonSplashEscape

	summary := presenter.Summary(match.NewSummarySnapshot(match.State{Encounter: state}))

	assert.Equal(t, encounter.StatusEscaped, summary.EncounterStatus)
	assert.Equal(t, "pez escapado", summary.OutcomeLabel)
}

func TestPresenterSummaryWithTrackCapture(t *testing.T) {
	presenter := NewPresenter(DefaultCatalog())
	state, err := encounter.NewState(encounter.DefaultConfig())
	require.NoError(t, err)
	state.Status = encounter.StatusCaptured
	state.EndReason = encounter.EndReasonTrackCapture

	summary := presenter.Summary(match.NewSummarySnapshot(match.State{Encounter: state}))

	assert.Equal(t, "captura por acercarlo a la orilla y subirlo a la superficie", summary.EndReasonLabel)
}

func TestPresenterSummaryWithTrackEscape(t *testing.T) {
	presenter := NewPresenter(DefaultCatalog())
	state, err := encounter.NewState(encounter.DefaultConfig())
	require.NoError(t, err)
	state.Status = encounter.StatusEscaped
	state.EndReason = encounter.EndReasonTrackEscape

	summary := presenter.Summary(match.NewSummarySnapshot(match.State{Encounter: state}))

	assert.Equal(t, "escape por superar la distancia maxima alcanzable", summary.EndReasonLabel)
}

func TestPresenterSummaryWithDepthEscape(t *testing.T) {
	presenter := NewPresenter(DefaultCatalog())
	state, err := encounter.NewState(encounter.DefaultConfig())
	require.NoError(t, err)
	state.Status = encounter.StatusEscaped
	state.EndReason = encounter.EndReasonDepthEscape

	summary := presenter.Summary(match.NewSummarySnapshot(match.State{Encounter: state}))

	assert.Equal(t, "escape por bajar mas alla de la profundidad alcanzable", summary.EndReasonLabel)
}

func TestPresenterSplashWithEscaped(t *testing.T) {
	presenter := NewPresenter(DefaultCatalog())
	snapshot := match.EncounterEventSnapshot{
		LastEvent: encounter.Event{
			Kind:    encounter.EventKindSplash,
			Escaped: true,
		},
	}

	view := presenter.Splash(snapshot, 2)

	assert.Contains(t, view.EventLabel, "chapotea")
	assert.Contains(t, view.EventLabel, "se suelta")
}

func TestPresenterPrivateLabelsAndHints(t *testing.T) {
	presenter := NewPresenter(DefaultCatalog())

	t.Run("player card hint renders normalized effects", func(t *testing.T) {
		hint := presenter.playerCardHint(match.MoveResourceSnapshot{
			HasTopCard: true,
			TopCard: cards.PlayerCard{Effects: []cards.CardEffect{{
				Trigger:                        cards.TriggerOnDraw,
				DistanceShift:                  1,
				DepthShift:                     -1,
				CaptureDistanceBonus:           2,
				SurfaceDepthBonus:              1,
				ExhaustionCaptureDistanceBonus: 3,
				Type:                           cards.EffectTypeHideDiscardTemporary,
			}}},
		})

		assert.Contains(t, hint, "draw")
		assert.Contains(t, hint, "dist +1")
		assert.Contains(t, hint, "prof -1")
		assert.Contains(t, hint, "capt +2")
		assert.Contains(t, hint, "sup +1")
		assert.Contains(t, hint, "baraja +3")
	})

	t.Run("effect impact parts include special effects", func(t *testing.T) {
		parts := effectImpactParts(cards.CardEffect{Type: cards.EffectTypeSuccessfulSplashApproach})
		assert.Contains(t, parts, "splash acerca")
	})

	t.Run("helper labels cover defaults", func(t *testing.T) {
		assert.Equal(t, "Tirar", presenter.playerMoveLabel(domain.Blue))
		assert.Equal(t, "unknown", presenter.playerMoveLabel(domain.Move(999)))
		assert.Equal(t, "Embestir", presenter.fishMoveLabel(domain.Blue))
		assert.Equal(t, "unknown", presenter.fishMoveLabel(domain.Move(999)))
		assert.Equal(t, "empate", presenter.roundOutcomeLabel(domain.Draw))
		assert.Equal(t, "unknown", presenter.roundOutcomeLabel(domain.RoundOutcome(999)))
		assert.Equal(t, "pez capturado", presenter.encounterOutcomeLabel(encounter.StatusCaptured))
		assert.Equal(t, "unknown", presenter.encounterOutcomeLabel(encounter.Status("unknown")))
		assert.Equal(t, "escape por chapoteo en superficie", presenter.endReasonLabel(encounter.EndReasonSplashEscape))
		assert.Equal(t, "unknown", presenter.endReasonLabel(encounter.EndReason("unknown")))
	})

	t.Run("run helpers cover fallback branches", func(t *testing.T) {
		assert.Equal(t, "run retirada", presenter.runStatusLabel(run.StatusRetired))
		assert.Equal(t, "mystery", presenter.runStatusLabel(run.Status("mystery")))
		assert.Equal(t, "captura confirmada: Lubina", presenter.runEncounterOutcomeLabel(run.EncounterResult{Outcome: run.EncounterOutcomeCaptured, Capture: &run.CaptureRecord{FishName: "Lubina"}}))
		assert.Equal(t, "weird", presenter.runEncounterOutcomeLabel(run.EncounterResult{Outcome: run.EncounterOutcome("weird")}))
		assert.Equal(t, "odd-node", presenter.runNodeLabel(run.NodeState{NodeID: "odd-node", Kind: run.NodeKind("weird")}))
	})

	t.Run("preset names resolve or fallback", func(t *testing.T) {
		assert.Equal(t, playerprofiles.DefaultPresets()[0].Name, playerDeckPresetName(playerprofiles.DefaultPresets()[0].ID))
		assert.Equal(t, rodpresets.DefaultPresets()[0].Name, rodPresetName(rodpresets.DefaultPresets()[0].ID))
		assert.Equal(t, attachmentpresets.DefaultPresets()[0].Name, attachmentPresetName(attachmentpresets.DefaultPresets()[0].ID))
		assert.Equal(t, "unknown-preset", playerDeckPresetName("unknown-preset"))
		assert.Equal(t, "unknown-preset", rodPresetName("unknown-preset"))
		assert.Equal(t, "unknown-preset", attachmentPresetName("unknown-preset"))
	})

	t.Run("discard labels cover unknown visibility and fallback names", func(t *testing.T) {
		assert.Equal(t, "Carta vista", presenter.fishDiscardEntryLabel(match.FishDiscardEntryState{Name: "Carta vista", Visibility: cards.DiscardVisibility("mystery")}))
		assert.Equal(t, "Embestir", presenter.fishDiscardEntryLabel(match.FishDiscardEntryState{Move: domain.Blue, Visibility: cards.DiscardVisibility("shadow")}))
	})

	t.Run("trigger and effect labels cover defaults", func(t *testing.T) {
		assert.Equal(t, "draw", triggerLabel(cards.TriggerOnDraw))
		assert.Equal(t, "efecto", triggerLabel(cards.Trigger(999)))
		assert.Equal(t, "avance horizontal", effectTypeLabel(cards.EffectTypeAdvanceHorizontal))
		assert.Equal(t, "custom", effectTypeLabel(cards.EffectType("custom")))
	})

	t.Run("splash renders non-escape event", func(t *testing.T) {
		view := presenter.Splash(match.EncounterEventSnapshot{LastEvent: encounter.Event{Kind: encounter.EventKindSplash, Escaped: false}}, 1)
		assert.Contains(t, view.EventLabel, "permanece sujeto")
	})
}
