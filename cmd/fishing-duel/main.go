package main

import (
	"fmt"
	"math/rand"
	"os"
	"pesca/internal/app"
	"pesca/internal/cards"
	"pesca/internal/cli"
	"pesca/internal/deck"
	"pesca/internal/encounter"
	"pesca/internal/endings"
	"pesca/internal/game"
	"pesca/internal/match"
	"pesca/internal/playermoves"
	"pesca/internal/playerrig"
	"pesca/internal/presentation"
	"pesca/internal/progression"
	"pesca/internal/rules"
	"time"
)

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffler := func(fishCards []cards.FishCard) {
		rng.Shuffle(len(fishCards), func(i, j int) {
			fishCards[i], fishCards[j] = fishCards[j], fishCards[i]
		})
	}
	playerCardShuffler := func(playerCards []cards.PlayerCard) {
		rng.Shuffle(len(playerCards), func(i, j int) {
			playerCards[i], playerCards[j] = playerCards[j], playerCards[i]
		})
	}
	ui := cli.NewUI(os.Stdin, os.Stdout)
	customFishDeck, err := ui.ChooseCustomFishDeck("Pesca: duelo contra el pez", deck.DefaultCustomFishDecks())
	if err != nil {
		exitWithError("error eligiendo preset de baraja", err)
	}

	fishDeck := customFishDeck.Build(shuffler)

	encounterState, err := encounter.NewState(encounter.DefaultConfig())
	if err != nil {
		exitWithError("error inicializando encuentro", err)
	}
	playerRigState, err := playerrig.NewState(playerrig.DefaultConfig())
	if err != nil {
		exitWithError("error configurando herramientas del jugador", err)
	}

	playerMoveConfig := playermoves.DefaultConfig()
	playerMoveConfig.DeckShuffler = playerCardShuffler
	playerMoveController, err := playermoves.NewUsageController(playerMoveConfig)
	if err != nil {
		exitWithError("error configurando movimientos del jugador", err)
	}

	engine, err := game.NewEngine(
		fishDeck,
		playerMoveController,
		rules.NewClassicEvaluator(rules.NewFishCombatProfile()),
		progression.TrackPolicy{SplashEscapeDecider: progression.SplashEscapeDeciderFunc(func(chance float64) bool {
			return rng.Float64() < chance
		})},
		endings.EncounterEndCondition{},
		match.State{Encounter: encounterState, PlayerRig: playerRigState},
	)
	if err != nil {
		exitWithError("error inicializando partida", err)
	}

	presenter := presentation.NewPresenter(presentation.DefaultCatalog())
	session, err := app.NewSession(engine, ui, presenter)
	if err != nil {
		exitWithError("error creando sesion", err)
	}

	if err := session.Run(); err != nil {
		exitWithError("error ejecutando partida", err)
	}
}

func exitWithError(message string, err error) {
	_, _ = fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
	os.Exit(1)
}
