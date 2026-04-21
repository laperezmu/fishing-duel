package cli

import (
	"bytes"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/game"
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

	move, err := ui.ChooseMove(presenter.Status(dummyState()), presenter.Intro().Options)
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
}

func TestChooseMoveShowsLastRoundSummary(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(strings.NewReader("2\n"), &out)
	presenter := presentation.NewPresenter(presentation.DefaultCatalog())
	require.NoError(t, ui.ShowIntro(presenter.Intro()))

	err := ui.ShowRound(presentation.RoundView{
		Status:      presentation.StatusView{Distance: 2},
		PlayerMove:  domain.Blue,
		FishMove:    domain.Yellow,
		PlayerLabel: "Tirar",
		FishLabel:   "Zafarse",
		Outcome:     "gana el jugador",
	})
	require.NoError(t, err)

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
	require.NoError(t, err)

	printed := out.String()
	assert.Contains(t, printed, "Ultimo lance")
	assert.Contains(t, printed, colorizeMove(domain.Blue, "Tirar"))
	assert.Contains(t, printed, colorizeMove(domain.Yellow, "Zafarse"))
	assert.Contains(t, printed, outcomeColor("gana el jugador"))
	assert.Contains(t, printed, "Distancia : 2")
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
