package main

import (
	"fmt"
	"math/rand"
	"os"
	"pesca/internal/app"
	"pesca/internal/cli"
	"pesca/internal/deck"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/endings"
	"pesca/internal/game"
	"pesca/internal/presentation"
	"pesca/internal/progression"
	"pesca/internal/rules"
	"time"
)

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffler := func(cards []domain.Move) {
		rng.Shuffle(len(cards), func(i, j int) {
			cards[i], cards[j] = cards[j], cards[i]
		})
	}

	manager := deck.NewManager(
		deck.NewStandardFishDeck(),
		shuffler,
		deck.RemoveCardsRecyclePolicy{CardsToRemove: 3},
	)

	encounterState, err := encounter.NewState(encounter.DefaultConfig())
	if err != nil {
		exitWithError("error inicializando encuentro", err)
	}

	engine, err := game.NewEngine(
		manager,
		rules.NewClassicEvaluator(rules.NewFishCombatProfile()),
		progression.TrackPolicy{},
		endings.EncounterCondition{},
		game.State{Encounter: encounterState},
	)
	if err != nil {
		exitWithError("error inicializando partida", err)
	}

	presenter := presentation.NewPresenter(presentation.DefaultCatalog())
	session, err := app.NewSession(engine, cli.NewUI(os.Stdin, os.Stdout), presenter)
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
