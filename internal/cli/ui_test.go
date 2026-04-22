package cli

import (
	"bytes"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/match"
	"pesca/internal/playerrig"
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
	assert.Contains(t, printed, colorizeMove(domain.Blue, "Tirar"))
	assert.Contains(t, printed, colorizeMove(domain.Red, "Recoger"))
	assert.Contains(t, printed, colorizeMove(domain.Yellow, "Soltar"))
	assert.Contains(t, printed, "[3/3]")
	assert.Contains(t, printed, "Superficie")
	assert.Contains(t, printed, "Nivel 1")
	assert.Contains(t, printed, "[SUP]")
	assert.Contains(t, printed, "[F]")
	assert.Contains(t, printed, "Profundidad actual: 1")
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

func samplePromptState(t *testing.T) match.State {
	t.Helper()

	encounterState, err := encounter.NewState(encounter.DefaultConfig())
	require.NoError(t, err)

	return match.State{
		Deck: match.DeckState{
			ActiveCards: 9,
		},
		Encounter: encounterState,
		PlayerRig: playerrig.State{MaxDistance: 5, MaxDepth: 4},
		PlayerMoves: match.PlayerMoveResources{Moves: []match.PlayerMoveState{
			{Move: domain.Blue, MaxUses: 3, RemainingUses: 3},
			{Move: domain.Red, MaxUses: 3, RemainingUses: 3},
			{Move: domain.Yellow, MaxUses: 3, RemainingUses: 3},
		}},
	}
}
