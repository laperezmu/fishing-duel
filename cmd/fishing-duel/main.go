package main

import (
	"fmt"
	"math/rand"
	"os"
	"pesca/internal/app"
	"pesca/internal/cards"
	"pesca/internal/cli"
	"pesca/internal/content/attachmentpresets"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/playerprofiles"
	"pesca/internal/content/rodpresets"
	"pesca/internal/content/watercontexts"
	"pesca/internal/encounter"
	"pesca/internal/endings"
	"pesca/internal/game"
	"pesca/internal/match"
	"pesca/internal/player/playermoves"
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
	playerDeckPreset, err := ui.ChoosePlayerDeckPreset("Pesca: duelo contra el pez", playerprofiles.DefaultPresets())
	if err != nil {
		exitWithError("error eligiendo preset del jugador", err)
	}
	rodPreset, err := ui.ChooseRodPreset("Pesca: duelo contra el pez", rodpresets.DefaultPresets())
	if err != nil {
		exitWithError("error eligiendo cana del jugador", err)
	}
	playerRod, err := rodPreset.BuildRod()
	if err != nil {
		exitWithError("error configurando cana del jugador", err)
	}
	attachmentPreset, err := ui.ChooseAttachmentPreset("Pesca: duelo contra el pez", playerRod, attachmentpresets.DefaultPresets())
	if err != nil {
		exitWithError("error eligiendo aditamentos del jugador", err)
	}
	playerLoadout, err := rodPreset.BuildLoadoutWithAttachments(attachmentPreset.BuildAttachments())
	if err != nil {
		exitWithError("error configurando loadout del jugador", err)
	}
	opening, err := app.ResolveEncounterOpening("Pesca: duelo contra el pez", encounter.DefaultConfig(), playerLoadout, watercontexts.DefaultPresets(), ui)
	if err != nil {
		exitWithError("error resolviendo la apertura de pesca", err)
	}
	spawn, err := app.ResolveFishSpawn("Pesca: duelo contra el pez", opening, playerLoadout, fishprofiles.DefaultProfiles(), ui)
	if err != nil {
		exitWithError("error resolviendo la aparicion del pez", err)
	}

	fishDeck := spawn.Profile.BuildPreset().BuildDeck(shuffler)

	encounterState, err := encounter.NewState(opening.Config)
	if err != nil {
		exitWithError("error inicializando encuentro", err)
	}
	playerMoveConfig := playerDeckPreset.BuildConfig(playerCardShuffler)
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
		match.State{Encounter: encounterState, Player: match.PlayerState{Loadout: playerLoadout}},
	)
	if err != nil {
		exitWithError("error inicializando partida", err)
	}

	if err := runSession(engine, ui); err != nil {
		exitWithError("error ejecutando partida", err)
	}
}

func runSession(engine *game.Engine, ui *cli.UI) error {
	presenter := presentation.NewPresenter(presentation.DefaultCatalog())
	session, err := app.NewSession(engine, ui, presenter)
	if err != nil {
		return err
	}

	return session.Run()
}

func exitWithError(message string, err error) {
	_, _ = fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
	os.Exit(1)
}
