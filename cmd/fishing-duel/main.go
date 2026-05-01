package main

import (
	"fmt"
	"os"
	"pesca/internal/app"
	"pesca/internal/cli"
	"pesca/internal/presentation"
)

func main() {
	ui := cli.NewUI(os.Stdin, os.Stdout)
	engine, err := app.BootstrapEncounter("Pesca: duelo contra el pez", app.DefaultRandomizer(), ui)
	if err != nil {
		exitWithError("error preparando encuentro", err)
	}

	presenter := presentation.NewPresenter(presentation.DefaultCatalog())
	session, err := app.NewSession(engine, ui, presenter)
	if err != nil {
		exitWithError("error inicializando partida", err)
	}
	if err := session.Run(); err != nil {
		exitWithError("error ejecutando partida", err)
	}
}

func exitWithError(message string, err error) {
	_, _ = fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
	os.Exit(1)
}
