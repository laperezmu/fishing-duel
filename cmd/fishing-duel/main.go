package main

import (
	"fmt"
	"os"
	"pesca/internal/app"
	"pesca/internal/cli"
	"pesca/internal/presentation"
)

func main() {
	title := "Pesca: sandbox de encounters"
	ui := cli.NewUI(os.Stdin, os.Stdout)
	catalog := presentation.DefaultCatalog()
	catalog.Title = title
	presenter := presentation.NewPresenter(catalog)
	sandbox, err := app.NewSandboxSession(title, app.DefaultRandomizer(), ui, presenter)
	if err != nil {
		exitWithError("error inicializando sandbox", err)
	}
	if len(os.Args) > 2 && os.Args[1] == "--scenario" {
		if err := sandbox.RunScenarioByID(os.Args[2]); err != nil {
			exitWithError("error ejecutando escenario", err)
		}
		return
	}
	if err := sandbox.Run(); err != nil {
		exitWithError("error ejecutando sandbox", err)
	}
}

func exitWithError(message string, err error) {
	_, _ = fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
	os.Exit(1)
}
