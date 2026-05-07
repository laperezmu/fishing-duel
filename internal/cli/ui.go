package cli

import (
	"bufio"
	"fmt"
	"io"
	"pesca/internal/content/anglerprofiles"
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
	opening    *presentation.OpeningView
	spawn      *presentation.SpawnView
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

	return choosePreset(ui, title, presets, renderPlayerDeckSelectionScreen, "Elige un preset del jugador: ", ui.confirmPlayerDeckPreset, "seleccion cancelada, elige otro preset")
}

func (ui *UI) ChooseAnglerProfile(title string, profiles []anglerprofiles.Profile) (anglerprofiles.Profile, error) {
	if len(profiles) == 0 {
		return anglerprofiles.Profile{}, fmt.Errorf("no hay pescadores iniciales disponibles")
	}

	presenter := presentation.NewPresenter(presentation.DefaultCatalog())
	return choosePreset(ui, title, profiles, renderAnglerProfileSelectionScreen, "Elige un pescador inicial: ", func(title string, profile anglerprofiles.Profile) (bool, error) {
		return ui.confirmAnglerProfile(title, presenter.AnglerProfile(profile))
	}, "seleccion cancelada, elige otro pescador")
}

func (ui *UI) ChooseFishDeckPreset(title string, presets []fishprofiles.FishDeckPreset) (fishprofiles.FishDeckPreset, error) {
	if len(presets) == 0 {
		return fishprofiles.FishDeckPreset{}, fmt.Errorf("no hay presets de baraja disponibles")
	}

	return choosePreset(ui, title, presets, renderFishDeckSelectionScreen, "Elige un preset del pez: ", ui.confirmFishDeckPreset, "seleccion cancelada, elige otro preset")
}

func (ui *UI) ChooseRodPreset(title string, presets []rodpresets.Preset) (rodpresets.Preset, error) {
	if len(presets) == 0 {
		return rodpresets.Preset{}, fmt.Errorf("no hay presets de cana disponibles")
	}

	return choosePreset(ui, title, presets, renderRodSelectionScreen, "Elige una cana del jugador: ", ui.confirmRodPreset, "seleccion cancelada, elige otra cana")
}

func (ui *UI) ChooseAttachmentPreset(title string, baseRod rod.State, presets []attachmentpresets.Preset) (attachmentpresets.Preset, error) {
	if len(presets) == 0 {
		return attachmentpresets.Preset{}, fmt.Errorf("no hay presets de aditamentos disponibles")
	}

	return choosePreset(ui, title, presets, renderAttachmentSelectionScreen, "Elige un preset de aditamentos: ", func(title string, preset attachmentpresets.Preset) (bool, error) {
		return ui.confirmAttachmentPreset(title, baseRod, preset)
	}, "seleccion cancelada, elige otros aditamentos")
}

func (ui *UI) ChooseMove(status presentation.StatusView, options []presentation.MoveOption) (domain.Move, error) {
	message := ""
	for {
		if _, err := io.WriteString(ui.out, renderPromptScreen(ui.intro.Title, status, options, ui.opening, ui.spawn, ui.lastRound, message)); err != nil {
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

func (ui *UI) ResolveSplash(view presentation.SplashView) (encounter.SplashResolution, error) {
	positions := ui.castFrames
	if len(positions) == 0 {
		positions = buildSplashPositions(view.TotalSlots)
	}

	inputCh, errCh := ui.startSplashInputWatcher()

	currentPosition := 0
	for _, position := range positions {
		currentPosition = position
		if _, err := io.WriteString(ui.out, renderSplashScreen(view, currentPosition, "")); err != nil {
			return encounter.SplashResolution{}, err
		}

		select {
		case err := <-errCh:
			return encounter.SplashResolution{}, err
		case <-inputCh:
			return ui.resolveSplashStop(view, currentPosition)
		default:
		}

		time.Sleep(ui.castDelay)
	}

	if _, err := io.WriteString(ui.out, clearSequence); err != nil {
		return encounter.SplashResolution{}, err
	}

	return encounter.SplashResolution{Escaped: true}, nil
}

func (ui *UI) startSplashInputWatcher() (<-chan struct{}, <-chan error) {
	inputCh := make(chan struct{}, 1)
	errCh := make(chan error, 1)
	go func() {
		if !ui.scanner.Scan() {
			if err := ui.scanner.Err(); err != nil {
				errCh <- err
				return
			}
			errCh <- fmt.Errorf("entrada finalizada")
			return
		}
		inputCh <- struct{}{}
	}()

	return inputCh, errCh
}

func (ui *UI) resolveSplashStop(view presentation.SplashView, currentPosition int) (encounter.SplashResolution, error) {
	if _, err := io.WriteString(ui.out, clearSequence); err != nil {
		return encounter.SplashResolution{}, err
	}
	if currentPosition >= view.TargetStart && currentPosition < view.TargetStart+view.TargetWidth {
		return encounter.SplashResolution{SuccessfulJumps: 1}, nil
	}

	return encounter.SplashResolution{Escaped: true}, nil
}

func (ui *UI) ShowGameOver(view presentation.SummaryView) error {
	_, err := io.WriteString(ui.out, renderGameOverScreen(ui.intro.Title, view, ui.lastRound))
	return err
}

func (ui *UI) ShowRunIntro(view presentation.RunIntroView) error {
	ui.intro = presentation.IntroView{Title: view.Title}
	_, err := io.WriteString(ui.out, renderRunIntroScreen(view))
	return err
}

func (ui *UI) ShowRunNode(view presentation.RunNodeView) error {
	_, err := io.WriteString(ui.out, renderRunNodeScreen(view))
	return err
}

func (ui *UI) ShowRunSummary(view presentation.RunSummaryView) error {
	_, err := io.WriteString(ui.out, renderRunSummaryScreen(view))
	return err
}

func (ui *UI) ShowRunNodeSummary(view presentation.RunNodeSummaryView) error {
	if _, err := io.WriteString(ui.out, renderRunNodeSummaryScreen(view)); err != nil {
		return err
	}
	if !ui.scanner.Scan() {
		if err := ui.scanner.Err(); err != nil {
			return err
		}
		return fmt.Errorf("entrada finalizada")
	}
	_, err := io.WriteString(ui.out, clearSequence)
	return err
}

func (ui *UI) ShowFishSpawn(_ string, spawn presentation.SpawnView) error {
	resolvedSpawn := spawn
	ui.spawn = &resolvedSpawn
	return nil
}

func (ui *UI) confirmPlayerDeckPreset(title string, preset playerprofiles.DeckPreset) (bool, error) {
	return confirmSelection(ui, func(message string) string {
		return renderPlayerDeckConfirmationScreen(title, preset, message)
	}, "Confirmar preset del jugador? [s/n]: ")
}

func (ui *UI) confirmAnglerProfile(title string, profile presentation.AnglerProfileView) (bool, error) {
	return confirmSelection(ui, func(message string) string {
		return renderAnglerProfileConfirmationScreen(title, profile, message)
	}, "Confirmar pescador? [s/n]: ")
}

func (ui *UI) confirmFishDeckPreset(title string, preset fishprofiles.FishDeckPreset) (bool, error) {
	return confirmSelection(ui, func(message string) string {
		return renderFishDeckConfirmationScreen(title, preset, message)
	}, "Confirmar preset del pez? [s/n]: ")
}

func (ui *UI) confirmRodPreset(title string, preset rodpresets.Preset) (bool, error) {
	return confirmSelection(ui, func(message string) string {
		return renderRodConfirmationScreen(title, preset, message)
	}, "Confirmar cana? [s/n]: ")
}

func (ui *UI) confirmAttachmentPreset(title string, baseRod rod.State, preset attachmentpresets.Preset) (bool, error) {
	return confirmSelectionWithError(ui, func(message string) (string, error) {
		previewLoadout, err := loadout.NewState(baseRod, preset.BuildAttachments())
		if err != nil {
			return "", err
		}

		return renderAttachmentConfirmationScreen(title, preset, previewLoadout, message), nil
	}, "Confirmar aditamentos? [s/n]: ")
}

func choosePreset[T any](ui *UI, title string, options []T, render func(string, []T, string) string, prompt string, confirm func(string, T) (bool, error), cancelledMessage string) (T, error) {
	var zero T
	message := ""
	for {
		if _, err := io.WriteString(ui.out, render(title, options, message)); err != nil {
			return zero, err
		}
		if _, err := fmt.Fprint(ui.out, prompt); err != nil {
			return zero, err
		}
		if !ui.scanner.Scan() {
			if err := ui.scanner.Err(); err != nil {
				return zero, err
			}
			return zero, fmt.Errorf("entrada finalizada")
		}

		selectedIndex, err := parsePresetChoice(ui.scanner.Text(), len(options))
		if err != nil {
			message = err.Error()
			continue
		}

		selected := options[selectedIndex]
		confirmed, err := confirm(title, selected)
		if err != nil {
			return zero, err
		}
		if confirmed {
			if _, err := io.WriteString(ui.out, clearSequence); err != nil {
				return zero, err
			}
			return selected, nil
		}

		message = cancelledMessage
	}
}

func confirmSelection(ui *UI, render func(string) string, prompt string) (bool, error) {
	return confirmSelectionWithError(ui, func(message string) (string, error) {
		return render(message), nil
	}, prompt)
}

func confirmSelectionWithError(ui *UI, render func(string) (string, error), prompt string) (bool, error) {
	message := ""
	for {
		screen, err := render(message)
		if err != nil {
			return false, err
		}
		if _, err := io.WriteString(ui.out, screen); err != nil {
			return false, err
		}
		if _, err := fmt.Fprint(ui.out, prompt); err != nil {
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

func buildSplashPositions(totalSlots int) []int {
	if totalSlots <= 1 {
		return []int{0}
	}

	positions := make([]int, 0, totalSlots*2-1)
	for position := 0; position < totalSlots; position++ {
		positions = append(positions, position)
	}
	for position := totalSlots - 2; position >= 0; position-- {
		positions = append(positions, position)
	}

	return positions
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
