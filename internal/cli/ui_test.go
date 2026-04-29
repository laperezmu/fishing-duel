package cli

import (
	"bytes"
	"pesca/internal/cards"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/playerprofiles"
	"pesca/internal/content/watercontexts"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/match"
	"pesca/internal/player/playermoves"
	"pesca/internal/player/playerrig"
	"pesca/internal/presentation"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShowIntroIncludesColoredOptions(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(strings.NewReader("1\n"), &out)
	presenter := presentation.NewPresenter(presentation.DefaultCatalog())

	err := ui.ShowIntro(presenter.Intro())
	require.NoError(t, err)

	status := presenter.Status(samplePromptState(t))
	move, err := ui.ChooseMove(status, status.MoveOptions)
	require.NoError(t, err)
	assert.Equal(t, domain.Blue, move)

	printed := out.String()
	assert.Contains(t, printed, clearSequence)
	assert.Contains(t, printed, "Tensa el sedal y arrastra al pez hacia la orilla.")
	assert.Contains(t, printed, "Orilla")
	assert.Contains(t, printed, "[ESC]")
	assert.Contains(t, printed, "ESC")
	assert.Contains(t, printed, "ESC | ")
	assert.Contains(t, printed, "  0 | ")
	assert.Contains(t, printed, "  1 | ")
	assert.Contains(t, printed, colorizeMove(domain.Blue, "Tirar"))
	assert.Contains(t, printed, colorizeMove(domain.Red, "Recoger"))
	assert.Contains(t, printed, colorizeMove(domain.Yellow, "Soltar"))
	assert.Contains(t, printed, "[3/3]")
	assert.Contains(t, printed, "{Anzuelo tenso}")
	assert.Contains(t, printed, "[F]")
	assert.Contains(t, printed, "~~~~")
	assert.Contains(t, printed, "Profundidad actual: 1")
	assert.Contains(t, printed, "Baraja agotada: captura con distancia <= 2 y profundidad <= 1")
	assert.Contains(t, printed, "Historial del pez")
	assert.Contains(t, printed, "Ciclo activo : C2 Oleaje abierto | ? | 1 oculta")
	assert.Contains(t, printed, "Reciclado   : rebaraja | retira 3 cartas | 1 ciclo cerrado")
	assert.Contains(t, printed, "Ciclos cerrados: C1 4 usadas, 1 oculta")
}

func TestChooseMoveShowsLastRoundSummary(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(strings.NewReader("2\n"), &out)
	presenter := presentation.NewPresenter(presentation.DefaultCatalog())
	require.NoError(t, ui.ShowIntro(presenter.Intro()))

	err := ui.ShowRound(presentation.RoundView{
		Status:       presentation.StatusView{FishDistance: 2, FishDepth: 1},
		PlayerMove:   domain.Blue,
		FishMove:     domain.Yellow,
		PlayerLabel:  "Tirar",
		FishLabel:    "Zafarse",
		Outcome:      domain.PlayerWin,
		OutcomeLabel: "gana el jugador",
		EventLabel:   "chapotea: permanece sujeto",
	})
	require.NoError(t, err)

	nextRoundStatus := presenter.Status(samplePromptState(t))
	nextRoundStatus.RoundNumber = 2
	nextRoundStatus.FishDistance = 2
	nextRoundStatus.FishDepth = 1
	nextRoundStatus.ActiveCards = 8
	nextRoundStatus.DiscardCards = 1
	nextRoundStatus.PlayerWins = 1

	_, err = ui.ChooseMove(nextRoundStatus, nextRoundStatus.MoveOptions)
	require.NoError(t, err)

	printed := out.String()
	assert.Contains(t, printed, "Ultimo lance")
	assert.Contains(t, printed, colorizeMove(domain.Blue, "Tirar"))
	assert.Contains(t, printed, colorizeMove(domain.Yellow, "Zafarse"))
	assert.Contains(t, printed, colorizeRoundOutcome(domain.PlayerWin, "gana el jugador"))
	assert.Contains(t, printed, "Distancia : 2")
	assert.Contains(t, printed, "Profundidad : 1")
	assert.Contains(t, printed, "Evento    : chapotea: permanece sujeto")
	assert.Contains(t, printed, "Historial del pez")
}

func TestChooseMoveRejectsUnavailableMoveUntilPlayerSelectsAvailableOption(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(strings.NewReader("1\n2\n"), &out)
	presenter := presentation.NewPresenter(presentation.DefaultCatalog())
	require.NoError(t, ui.ShowIntro(presenter.Intro()))

	status := presenter.Status(samplePromptState(t))
	status.MoveOptions[0].RemainingUses = 0
	status.MoveOptions[0].Available = false
	status.MoveOptions[0].RestoresOnRound = 3

	move, err := ui.ChooseMove(status, status.MoveOptions)
	require.NoError(t, err)
	assert.Equal(t, domain.Red, move)
	assert.Contains(t, out.String(), "la accion tirar recarga en la ronda 3")
}

func TestChoosePlayerDeckPreset(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(strings.NewReader("2\ns\n"), &out)

	preset, err := ui.ChoosePlayerDeckPreset("Pesca: duelo contra el pez", samplePlayerDeckPresets())

	require.NoError(t, err)
	assert.Equal(t, "Apertura", preset.Name)
	assert.Contains(t, out.String(), "Preset del jugador")
	assert.Contains(t, out.String(), "Azul - Anzuelo tenso")
	assert.Contains(t, out.String(), clearSequence)
}

func TestChooseFishDeckPreset(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(strings.NewReader("2\ns\n"), &out)

	preset, err := ui.ChooseFishDeckPreset("Pesca: duelo contra el pez", sampleFishDeckPresets())

	require.NoError(t, err)
	assert.Equal(t, "Apertura", preset.Name)
	assert.Contains(t, out.String(), "Preset del pez")
	assert.Contains(t, out.String(), "Confirmar preset")
	assert.Contains(t, out.String(), "Apertura")
	assert.Contains(t, out.String(), "Rojo - Tiron de apertura")
	assert.Contains(t, out.String(), clearSequence)
}

func TestChooseFishDeckPresetReturnsToSelectionAfterCancellingConfirmation(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(strings.NewReader("1\nn\n2\ns\n"), &out)

	preset, err := ui.ChooseFishDeckPreset("Pesca: duelo contra el pez", sampleFishDeckPresets())

	require.NoError(t, err)
	assert.Equal(t, "Apertura", preset.Name)
	assert.Contains(t, out.String(), "seleccion cancelada, elige otro preset")
}

func TestChooseFishDeckPresetRejectsInvalidInput(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(strings.NewReader("9\n2\nquizas\ns\n"), &out)

	preset, err := ui.ChooseFishDeckPreset("Pesca: duelo contra el pez", sampleFishDeckPresets())

	require.NoError(t, err)
	assert.Equal(t, "Apertura", preset.Name)
	assert.Contains(t, out.String(), "opcion no valida, usa un numero entre 1 y 2")
	assert.Contains(t, out.String(), "respuesta no valida, usa s o n")
}

func TestChooseWaterContext(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(strings.NewReader("2\ns\n"), &out)

	preset, err := ui.ChooseWaterContext("Pesca: duelo contra el pez", sampleWaterContextPresets())

	require.NoError(t, err)
	assert.Equal(t, "Canal abierto", preset.Name)
	assert.Contains(t, out.String(), "Situacion de agua")
	assert.Contains(t, out.String(), "Confirmar situacion de agua")
	assert.NotContains(t, out.String(), "offshore")
	assert.Contains(t, out.String(), clearSequence)
}

func TestResolveCastUsesOscillatingBarAndStoresOpeningSummary(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(strings.NewReader("\n1\n"), &out)
	ui.castDelay = 0
	ui.castFrames = []int{0}
	waterContext := sampleWaterContextPresets()[0].BuildContext()

	castResult, err := ui.ResolveCast("Pesca: duelo contra el pez", waterContext)
	require.NoError(t, err)
	assert.Equal(t, encounter.CastBandVeryShort, castResult.Band)

	opening, err := encounter.ResolveOpening(encounter.DefaultConfig(), waterContext, castResult)
	require.NoError(t, err)
	require.NoError(t, ui.ShowEncounterOpening("Pesca: duelo contra el pez", opening))

	presenter := presentation.NewPresenter(presentation.DefaultCatalog())
	require.NoError(t, ui.ShowIntro(presenter.Intro()))
	status := presenter.Status(samplePromptState(t))
	_, err = ui.ChooseMove(status, status.MoveOptions)
	require.NoError(t, err)

	printed := out.String()
	assert.Contains(t, printed, "Lectura del agua")
	assert.Contains(t, printed, "Pulsa Enter para detener la barra")
	assert.Contains(t, printed, "Barra      : [")
	assert.Contains(t, printed, "Apertura del lance")
	assert.Contains(t, printed, "Agua       : Ensenada cercana")
	assert.Contains(t, printed, "Lance      : muy corto")
	assert.Contains(t, printed, "Inicio     : distancia 0 | profundidad 1")
}

func TestPresetSelectionScreensHideCardDetailsFromTheList(t *testing.T) {
	playerSelection := renderPlayerDeckSelectionSection(samplePlayerDeckPresets())
	fishSelection := renderFishDeckSelectionSection(sampleFishDeckPresets())

	assert.NotContains(t, playerSelection, "Azul - Anzuelo tenso")
	assert.NotContains(t, fishSelection, "Rojo - Tiron de apertura")

	playerConfirmation := renderPlayerDeckConfirmationSection(samplePlayerDeckPresets()[1])
	fishConfirmation := renderFishDeckConfirmationSection(sampleFishDeckPresets()[1])

	assert.Contains(t, playerConfirmation, "Azul - Anzuelo tenso")
	assert.Contains(t, fishConfirmation, "Rojo - Tiron de apertura")
}

func samplePromptState(t *testing.T) match.State {
	t.Helper()

	encounterState, err := encounter.NewState(encounter.DefaultConfig())
	require.NoError(t, err)

	return match.State{
		Deck: match.DeckState{
			ActiveCards:       9,
			DiscardCards:      3,
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
		Encounter: encounterState,
		PlayerRig: playerrig.State{MaxDistance: 5, MaxDepth: 4},
		PlayerMoves: match.PlayerMoveResources{Moves: []match.PlayerMoveState{
			{Move: domain.Blue, MaxUses: 3, RemainingUses: 3, ActiveCards: []cards.PlayerCard{cards.NewNamedPlayerCard("Anzuelo tenso", "Capturas desde un paso mas lejos este round.", domain.Blue, cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1}), cards.NewPlayerCard(domain.Blue), cards.NewPlayerCard(domain.Blue)}},
			{Move: domain.Red, MaxUses: 3, RemainingUses: 3, ActiveCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Red), cards.NewPlayerCard(domain.Red), cards.NewPlayerCard(domain.Red)}},
			{Move: domain.Yellow, MaxUses: 3, RemainingUses: 3, ActiveCards: []cards.PlayerCard{cards.NewPlayerCard(domain.Yellow), cards.NewPlayerCard(domain.Yellow), cards.NewPlayerCard(domain.Yellow)}},
		}},
	}
}

func sampleFishDeckPresets() []fishprofiles.FishDeckPreset {
	return []fishprofiles.FishDeckPreset{
		{
			ID:            "classic",
			ArchetypeID:   fishprofiles.ArchetypeBaselineCycle,
			Name:          "Clasico",
			Description:   "Sin efectos.",
			Details:       []string{"3 cartas lisas."},
			FishCards:     []cards.FishCard{cards.NewFishCard(domain.Blue)},
			CardsToRemove: 3,
			Shuffle:       true,
		},
		{
			ID:          "hooked-opening",
			ArchetypeID: fishprofiles.ArchetypeDrawTempo,
			Name:        "Apertura",
			Description: "Con on_draw.",
			Details:     []string{"Rojo - Tiron de apertura: al revelarse permite capturar desde un paso mas lejos ese round."},
			FishCards: []cards.FishCard{
				cards.NewFishCard(domain.Red, cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1}),
			},
			CardsToRemove: 0,
			Shuffle:       false,
		},
	}
}

func sampleWaterContextPresets() []watercontexts.Preset {
	return watercontexts.DefaultPresets()
}

func samplePlayerDeckPresets() []playerprofiles.DeckPreset {
	return []playerprofiles.DeckPreset{
		{
			ID:          "classic",
			Name:        "Clasico",
			Description: "Sin efectos.",
			Details:     []string{"Azul: 3 cartas lisas."},
			Config:      playermoves.DefaultConfig(),
		},
		{
			ID:          "hooked-opening",
			Name:        "Apertura",
			Description: "Con on_draw.",
			Details:     []string{"Azul - Anzuelo tenso: al revelar la carta permite capturar desde un paso mas lejos ese round."},
			Config: playermoves.Config{
				InitialDecks: map[domain.Move][]cards.PlayerCard{
					domain.Blue:   {cards.NewNamedPlayerCard("Anzuelo tenso", "Capturas desde un paso mas lejos este round.", domain.Blue, cards.CardEffect{Trigger: cards.TriggerOnDraw, CaptureDistanceBonus: 1})},
					domain.Red:    {cards.NewPlayerCard(domain.Red)},
					domain.Yellow: {cards.NewPlayerCard(domain.Yellow)},
				},
				RecoveryDelayRounds: 1,
			},
		},
	}
}
