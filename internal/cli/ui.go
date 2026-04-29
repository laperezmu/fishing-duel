package cli

import (
	"bufio"
	"fmt"
	"io"
	"pesca/internal/content/attachmentpresets"
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/playerprofiles"
	"pesca/internal/content/rodpresets"
	"pesca/internal/domain"
	"pesca/internal/encounter"
	"pesca/internal/player/loadout"
	"pesca/internal/player/rod"
	"pesca/internal/presentation"
	"strconv"
	"strings"
	"time"
)

type UI struct {
	scanner    *bufio.Scanner
	out        io.Writer
	intro      presentation.IntroView
	lastRound  *presentation.RoundView
	opening    *encounter.Opening
	castDelay  time.Duration
	castFrames []int
}

func NewUI(in io.Reader, out io.Writer) *UI {
	return &UI{
		scanner:   bufio.NewScanner(in),
		out:       out,
		castDelay: 50 * time.Millisecond,
	}
}

func (ui *UI) ShowIntro(view presentation.IntroView) error {
	ui.intro = view
	return nil
}

func (ui *UI) ChoosePlayerDeckPreset(title string, presets []playerprofiles.DeckPreset) (playerprofiles.DeckPreset, error) {
	if len(presets) == 0 {
		return playerprofiles.DeckPreset{}, fmt.Errorf("no hay presets de baraja del jugador disponibles")
	}

	message := ""
	for {
		if _, err := io.WriteString(ui.out, renderPlayerDeckSelectionScreen(title, presets, message)); err != nil {
			return playerprofiles.DeckPreset{}, err
		}
		if _, err := fmt.Fprint(ui.out, "Elige un preset del jugador: "); err != nil {
			return playerprofiles.DeckPreset{}, err
		}
		if !ui.scanner.Scan() {
			if err := ui.scanner.Err(); err != nil {
				return playerprofiles.DeckPreset{}, err
			}
			return playerprofiles.DeckPreset{}, fmt.Errorf("entrada finalizada")
		}

		selectedIndex, err := parsePresetChoice(ui.scanner.Text(), len(presets))
		if err != nil {
			message = err.Error()
			continue
		}

		selectedPreset := presets[selectedIndex]
		confirmed, err := ui.confirmPlayerDeckPreset(title, selectedPreset)
		if err != nil {
			return playerprofiles.DeckPreset{}, err
		}
		if confirmed {
			if _, err := io.WriteString(ui.out, clearSequence); err != nil {
				return playerprofiles.DeckPreset{}, err
			}
			return selectedPreset, nil
		}

		message = "seleccion cancelada, elige otro preset"
	}
}

func (ui *UI) ChooseFishDeckPreset(title string, presets []fishprofiles.FishDeckPreset) (fishprofiles.FishDeckPreset, error) {
	if len(presets) == 0 {
		return fishprofiles.FishDeckPreset{}, fmt.Errorf("no hay presets de baraja disponibles")
	}

	message := ""
	for {
		if _, err := io.WriteString(ui.out, renderFishDeckSelectionScreen(title, presets, message)); err != nil {
			return fishprofiles.FishDeckPreset{}, err
		}
		if _, err := fmt.Fprint(ui.out, "Elige un preset del pez: "); err != nil {
			return fishprofiles.FishDeckPreset{}, err
		}
		if !ui.scanner.Scan() {
			if err := ui.scanner.Err(); err != nil {
				return fishprofiles.FishDeckPreset{}, err
			}
			return fishprofiles.FishDeckPreset{}, fmt.Errorf("entrada finalizada")
		}

		selectedIndex, err := parsePresetChoice(ui.scanner.Text(), len(presets))
		if err != nil {
			message = err.Error()
			continue
		}

		selectedPreset := presets[selectedIndex]
		confirmed, err := ui.confirmFishDeckPreset(title, selectedPreset)
		if err != nil {
			return fishprofiles.FishDeckPreset{}, err
		}
		if confirmed {
			if _, err := io.WriteString(ui.out, clearSequence); err != nil {
				return fishprofiles.FishDeckPreset{}, err
			}
			return selectedPreset, nil
		}

		message = "seleccion cancelada, elige otro preset"
	}
}

func (ui *UI) ChooseRodPreset(title string, presets []rodpresets.Preset) (rodpresets.Preset, error) {
	if len(presets) == 0 {
		return rodpresets.Preset{}, fmt.Errorf("no hay presets de cana disponibles")
	}

	message := ""
	for {
		if _, err := io.WriteString(ui.out, renderRodSelectionScreen(title, presets, message)); err != nil {
			return rodpresets.Preset{}, err
		}
		if _, err := fmt.Fprint(ui.out, "Elige una cana del jugador: "); err != nil {
			return rodpresets.Preset{}, err
		}
		if !ui.scanner.Scan() {
			if err := ui.scanner.Err(); err != nil {
				return rodpresets.Preset{}, err
			}
			return rodpresets.Preset{}, fmt.Errorf("entrada finalizada")
		}

		selectedIndex, err := parsePresetChoice(ui.scanner.Text(), len(presets))
		if err != nil {
			message = err.Error()
			continue
		}

		selectedPreset := presets[selectedIndex]
		confirmed, err := ui.confirmRodPreset(title, selectedPreset)
		if err != nil {
			return rodpresets.Preset{}, err
		}
		if confirmed {
			if _, err := io.WriteString(ui.out, clearSequence); err != nil {
				return rodpresets.Preset{}, err
			}
			return selectedPreset, nil
		}

		message = "seleccion cancelada, elige otra cana"
	}
}

func (ui *UI) ChooseAttachmentPreset(title string, baseRod rod.State, presets []attachmentpresets.Preset) (attachmentpresets.Preset, error) {
	if len(presets) == 0 {
		return attachmentpresets.Preset{}, fmt.Errorf("no hay presets de aditamentos disponibles")
	}

	message := ""
	for {
		if _, err := io.WriteString(ui.out, renderAttachmentSelectionScreen(title, presets, message)); err != nil {
			return attachmentpresets.Preset{}, err
		}
		if _, err := fmt.Fprint(ui.out, "Elige un preset de aditamentos: "); err != nil {
			return attachmentpresets.Preset{}, err
		}
		if !ui.scanner.Scan() {
			if err := ui.scanner.Err(); err != nil {
				return attachmentpresets.Preset{}, err
			}
			return attachmentpresets.Preset{}, fmt.Errorf("entrada finalizada")
		}

		selectedIndex, err := parsePresetChoice(ui.scanner.Text(), len(presets))
		if err != nil {
			message = err.Error()
			continue
		}

		selectedPreset := presets[selectedIndex]
		confirmed, err := ui.confirmAttachmentPreset(title, baseRod, selectedPreset)
		if err != nil {
			return attachmentpresets.Preset{}, err
		}
		if confirmed {
			if _, err := io.WriteString(ui.out, clearSequence); err != nil {
				return attachmentpresets.Preset{}, err
			}
			return selectedPreset, nil
		}

		message = "seleccion cancelada, elige otros aditamentos"
	}
}

func (ui *UI) ChooseMove(status presentation.StatusView, options []presentation.MoveOption) (domain.Move, error) {
	message := ""
	for {
		if _, err := io.WriteString(ui.out, renderPromptScreen(ui.intro.Title, status, options, ui.opening, ui.lastRound, message)); err != nil {
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

func (ui *UI) confirmPlayerDeckPreset(title string, preset playerprofiles.DeckPreset) (bool, error) {
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

func (ui *UI) confirmFishDeckPreset(title string, preset fishprofiles.FishDeckPreset) (bool, error) {
	message := ""
	for {
		if _, err := io.WriteString(ui.out, renderFishDeckConfirmationScreen(title, preset, message)); err != nil {
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

func (ui *UI) confirmRodPreset(title string, preset rodpresets.Preset) (bool, error) {
	message := ""
	for {
		if _, err := io.WriteString(ui.out, renderRodConfirmationScreen(title, preset, message)); err != nil {
			return false, err
		}
		if _, err := fmt.Fprint(ui.out, "Confirmar cana? [s/n]: "); err != nil {
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

func (ui *UI) confirmAttachmentPreset(title string, baseRod rod.State, preset attachmentpresets.Preset) (bool, error) {
	message := ""
	for {
		previewLoadout, err := loadout.NewState(baseRod, preset.BuildAttachments())
		if err != nil {
			return false, err
		}
		if _, err := io.WriteString(ui.out, renderAttachmentConfirmationScreen(title, preset, previewLoadout, message)); err != nil {
			return false, err
		}
		if _, err := fmt.Fprint(ui.out, "Confirmar aditamentos? [s/n]: "); err != nil {
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
		return false, fmt.Errorf("respuesta no valida, usa s o n")
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
