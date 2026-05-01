package cli

import (
	"fmt"
	"io"
	"pesca/internal/app"
	"pesca/internal/content/watercontexts"
	"pesca/internal/encounter"
	"pesca/internal/presentation"
	"time"
)

const castBarSectionWidth = 4

func (ui *UI) ChooseWaterContext(title string, presets []watercontexts.Preset) (watercontexts.Preset, error) {
	if len(presets) == 0 {
		return watercontexts.Preset{}, fmt.Errorf("no hay situaciones de agua disponibles")
	}

	message := ""
	for {
		if _, err := io.WriteString(ui.out, renderWaterContextSelectionScreen(title, presets, message)); err != nil {
			return watercontexts.Preset{}, err
		}
		if _, err := fmt.Fprint(ui.out, "Elige una situacion de agua: "); err != nil {
			return watercontexts.Preset{}, err
		}
		if !ui.scanner.Scan() {
			if err := ui.scanner.Err(); err != nil {
				return watercontexts.Preset{}, err
			}
			return watercontexts.Preset{}, fmt.Errorf("entrada finalizada")
		}

		selectedIndex, err := parsePresetChoice(ui.scanner.Text(), len(presets))
		if err != nil {
			message = err.Error()
			continue
		}

		selectedPreset := presets[selectedIndex]
		confirmed, err := ui.confirmWaterContext(title, selectedPreset)
		if err != nil {
			return watercontexts.Preset{}, err
		}
		if confirmed {
			if _, err := io.WriteString(ui.out, clearSequence); err != nil {
				return watercontexts.Preset{}, err
			}
			return selectedPreset, nil
		}

		message = "seleccion cancelada, elige otra situacion"
	}
}

func (ui *UI) ResolveCast(_ string, context encounter.WaterContext, presenter app.CastPresenter) (encounter.CastResult, error) {
	controller := encounter.NewCastController(castBarSectionWidth, ui.castFrames)

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

	currentPosition := controller.CurrentPosition()
	for {
		for step := 0; step < len(ui.castFrames)+controller.TotalSlots(); step++ {
			castView := presenter.Cast(context, currentPosition, controller.TotalSlots(), castBarSectionWidth)
			if _, err := io.WriteString(ui.out, renderCastScreen(castView, "")); err != nil {
				return encounter.CastResult{}, err
			}

			select {
			case err := <-errCh:
				return encounter.CastResult{}, err
			case <-inputCh:
				if _, err := io.WriteString(ui.out, clearSequence); err != nil {
					return encounter.CastResult{}, err
				}
				return encounter.CastResult{Band: controller.ResolveBand(currentPosition)}, nil
			default:
			}

			time.Sleep(ui.castDelay)
			currentPosition = controller.Advance()
		}
	}
}

func (ui *UI) ShowEncounterOpening(_ string, opening presentation.OpeningView) error {
	ui.opening = &opening
	return nil
}

func (ui *UI) confirmWaterContext(title string, preset watercontexts.Preset) (bool, error) {
	message := ""
	for {
		if _, err := io.WriteString(ui.out, renderWaterContextConfirmationScreen(title, preset, message)); err != nil {
			return false, err
		}
		if _, err := fmt.Fprint(ui.out, "Confirmar situacion de agua? [s/n]: "); err != nil {
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
