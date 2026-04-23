package cli

import (
	"bufio"
	"fmt"
	"io"
	"pesca/internal/deck"
	"pesca/internal/domain"
	"pesca/internal/playermoves"
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

func (ui *UI) ChoosePlayerDeckPreset(title string, presets []playermoves.PlayerDeckPreset) (playermoves.PlayerDeckPreset, error) {
	if len(presets) == 0 {
		return playermoves.PlayerDeckPreset{}, fmt.Errorf("no hay presets de baraja del jugador disponibles")
	}

	message := ""
	for {
		if _, err := io.WriteString(ui.out, renderPlayerDeckSelectionScreen(title, presets, message)); err != nil {
			return playermoves.PlayerDeckPreset{}, err
		}
		if _, err := fmt.Fprint(ui.out, "Elige un preset del jugador: "); err != nil {
			return playermoves.PlayerDeckPreset{}, err
		}
		if !ui.scanner.Scan() {
			if err := ui.scanner.Err(); err != nil {
				return playermoves.PlayerDeckPreset{}, err
			}
			return playermoves.PlayerDeckPreset{}, fmt.Errorf("entrada finalizada")
		}

		selectedIndex, err := parsePresetChoice(ui.scanner.Text(), len(presets))
		if err != nil {
			message = err.Error()
			continue
		}

		selectedPreset := presets[selectedIndex]
		confirmed, err := ui.confirmPlayerDeckPreset(title, selectedPreset)
		if err != nil {
			return playermoves.PlayerDeckPreset{}, err
		}
		if confirmed {
			if _, err := io.WriteString(ui.out, clearSequence); err != nil {
				return playermoves.PlayerDeckPreset{}, err
			}
			return selectedPreset, nil
		}

		message = "seleccion cancelada, elige otro preset"
	}
}

func (ui *UI) ChooseCustomFishDeck(title string, customFishDecks []deck.CustomFishDeck) (deck.CustomFishDeck, error) {
	if len(customFishDecks) == 0 {
		return deck.CustomFishDeck{}, fmt.Errorf("no hay presets de baraja disponibles")
	}

	message := ""
	for {
		if _, err := io.WriteString(ui.out, renderCustomFishDeckSelectionScreen(title, customFishDecks, message)); err != nil {
			return deck.CustomFishDeck{}, err
		}
		if _, err := fmt.Fprint(ui.out, "Elige un preset del pez: "); err != nil {
			return deck.CustomFishDeck{}, err
		}
		if !ui.scanner.Scan() {
			if err := ui.scanner.Err(); err != nil {
				return deck.CustomFishDeck{}, err
			}
			return deck.CustomFishDeck{}, fmt.Errorf("entrada finalizada")
		}

		selectedIndex, err := parsePresetChoice(ui.scanner.Text(), len(customFishDecks))
		if err != nil {
			message = err.Error()
			continue
		}

		selectedCustomFishDeck := customFishDecks[selectedIndex]
		confirmed, err := ui.confirmCustomFishDeck(title, selectedCustomFishDeck)
		if err != nil {
			return deck.CustomFishDeck{}, err
		}
		if confirmed {
			if _, err := io.WriteString(ui.out, clearSequence); err != nil {
				return deck.CustomFishDeck{}, err
			}
			return selectedCustomFishDeck, nil
		}

		message = "seleccion cancelada, elige otro preset"
	}
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

func (ui *UI) confirmPlayerDeckPreset(title string, preset playermoves.PlayerDeckPreset) (bool, error) {
	message := ""
	for {
		if _, err := io.WriteString(ui.out, renderPlayerDeckConfirmationScreen(title, preset, message)); err != nil {
			return false, err
		}
		if _, err := fmt.Fprint(ui.out, "Confirmar preset del jugador? [s/n]: "); err != nil {
			return false, err
		}
		if !ui.scanner.Scan() {
			if err := ui.scanner.Err(); err != nil {
				return false, err
			}
			return false, fmt.Errorf("entrada finalizada")
		}

		confirmed, err := parseConfirmation(ui.scanner.Text())
		if err == nil {
			return confirmed, nil
		}

		message = err.Error()
	}
}

func (ui *UI) confirmCustomFishDeck(title string, customFishDeck deck.CustomFishDeck) (bool, error) {
	message := ""
	for {
		if _, err := io.WriteString(ui.out, renderCustomFishDeckConfirmationScreen(title, customFishDeck, message)); err != nil {
			return false, err
		}
		if _, err := fmt.Fprint(ui.out, "Confirmar preset del pez? [s/n]: "); err != nil {
			return false, err
		}
		if !ui.scanner.Scan() {
			if err := ui.scanner.Err(); err != nil {
				return false, err
			}
			return false, fmt.Errorf("entrada finalizada")
		}

		confirmed, err := parseConfirmation(ui.scanner.Text())
		if err == nil {
			return confirmed, nil
		}

		message = err.Error()
	}
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

func parsePresetChoice(input string, availablePresets int) (int, error) {
	choice, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return 0, fmt.Errorf("opcion no valida, usa un numero de preset")
	}
	if choice < 1 || choice > availablePresets {
		return 0, fmt.Errorf("opcion no valida, usa un numero entre 1 y %d", availablePresets)
	}

	return choice - 1, nil
}

func parseConfirmation(input string) (bool, error) {
	trimmed := strings.TrimSpace(strings.ToLower(input))
	trimmed = strings.ReplaceAll(trimmed, "í", "i")

	switch trimmed {
	case "s", "si":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("confirmacion no valida, usa s o n")
	}
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
