package cli

import (
	"bufio"
	"fmt"
	"io"
	"pesca/internal/domain"
	"pesca/internal/presentation"
	"strconv"
	"strings"
)

type UI struct {
	scanner   *bufio.Scanner
	out       io.Writer
	intro     presentation.IntroView
	lastRound *presentation.RoundView
}

func NewUI(in io.Reader, out io.Writer) *UI {
	return &UI{
		scanner: bufio.NewScanner(in),
		out:     out,
	}
}

func (ui *UI) ShowIntro(view presentation.IntroView) error {
	ui.intro = view
	return nil
}

func (ui *UI) ChooseMove(status presentation.StatusView, options []presentation.MoveOption) (domain.Move, error) {
	message := ""
	for {
		if _, err := io.WriteString(ui.out, renderPromptScreen(ui.intro.Title, status, options, ui.lastRound, message)); err != nil {
			return 0, err
		}
		if _, err := fmt.Fprint(ui.out, "Elige 1, 2 o 3: "); err != nil {
			return 0, err
		}
		if !ui.scanner.Scan() {
			if err := ui.scanner.Err(); err != nil {
				return 0, err
			}
			return 0, fmt.Errorf("entrada finalizada")
		}

		move, err := parseMove(ui.scanner.Text(), options)
		if err == nil {
			ui.lastRound = nil
			if _, err := io.WriteString(ui.out, clearSequence); err != nil {
				return 0, err
			}
			return move, nil
		}

		message = err.Error()
	}
}

func (ui *UI) ShowRound(view presentation.RoundView) error {
	round := view
	ui.lastRound = &round
	return nil
}

func (ui *UI) ShowGameOver(view presentation.SummaryView) error {
	_, err := io.WriteString(ui.out, renderGameOverScreen(ui.intro.Title, view, ui.lastRound))
	return err
}

func parseMove(input string, options []presentation.MoveOption) (domain.Move, error) {
	trimmed := strings.TrimSpace(strings.ToLower(input))

	for _, option := range options {
		if trimmed == strings.ToLower(option.Label) {
			return selectMoveOption(option)
		}
	}

	choice, err := strconv.Atoi(trimmed)
	if err != nil {
		return 0, fmt.Errorf("opcion no valida, usa 1, 2 o 3")
	}

	for _, option := range options {
		if option.Index == choice {
			return selectMoveOption(option)
		}
	}

	return 0, fmt.Errorf("opcion no valida, usa 1, 2 o 3")
}

func selectMoveOption(option presentation.MoveOption) (domain.Move, error) {
	if option.Available {
		return option.Move, nil
	}

	if option.RestoresOnRound > 0 {
		return 0, fmt.Errorf("la accion %s recarga en la ronda %d", strings.ToLower(option.Label), option.RestoresOnRound)
	}

	return 0, fmt.Errorf("la accion %s no tiene usos disponibles", strings.ToLower(option.Label))
}
