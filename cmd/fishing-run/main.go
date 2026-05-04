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
	presenter := presentation.NewPresenter(presentation.DefaultCatalog())
	runSession, err := app.NewRunSession("Pesca: run MVP", app.DefaultRandomizer(), ui, presenter)
	if err != nil {
		exitWithError("error inicializando run", err)
	}
	if err := runSession.Run(); err != nil {
		exitWithError("error ejecutando run", err)
	}
}

func exitWithError(message string, err error) {
	_, _ = fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
	os.Exit(1)
}
