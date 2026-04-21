package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

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
		fmt.Fprintf(os.Stderr, "error inicializando encuentro: %v\n", err)
		os.Exit(1)
	}

	engine, err := game.NewEngine(
		manager,
		rules.ClassicEvaluator{},
		progression.TrackPolicy{},
		endings.EncounterCondition{},
		game.State{Encounter: encounterState},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error inicializando partida: %v\n", err)
		os.Exit(1)
	}

	presenter := presentation.NewPresenter(presentation.DefaultCatalog())
	session, err := app.NewSession(engine, cli.NewUI(os.Stdin, os.Stdout), presenter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creando sesion: %v\n", err)
		os.Exit(1)
	}

	if err := session.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error ejecutando partida: %v\n", err)
		os.Exit(1)
	}
}
