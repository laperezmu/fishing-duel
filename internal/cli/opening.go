package cli

import (
	"fmt"
	"io"
	"pesca/internal/content/watercontexts"
	"pesca/internal/encounter"
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

func (ui *UI) ResolveCast(title string, context encounter.WaterContext) (encounter.CastResult, error) {
	bands := encounter.OrderedCastBands()
	totalSlots := len(bands) * castBarSectionWidth
	positions := ui.castFrames
	if len(positions) == 0 {
		positions = buildCastPositions(totalSlots)
	}

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

	currentPosition := 0
	for {
		for _, position := range positions {
			currentPosition = position
			if _, err := io.WriteString(ui.out, renderCastScreen(title, context, currentPosition, totalSlots, castBarSectionWidth, "")); err != nil {
				return encounter.CastResult{}, err
			}

			select {
			case err := <-errCh:
				return encounter.CastResult{}, err
			case <-inputCh:
				if _, err := io.WriteString(ui.out, clearSequence); err != nil {
					return encounter.CastResult{}, err
				}
				return encounter.CastResult{Band: castBandForPosition(currentPosition, castBarSectionWidth)}, nil
			default:
			}

			time.Sleep(ui.castDelay)
		}
	}
}

func (ui *UI) ShowEncounterOpening(_ string, opening encounter.Opening) error {
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

func buildCastPositions(totalSlots int) []int {
	if totalSlots <= 1 {
		return []int{0}
	}

	positions := make([]int, 0, totalSlots*2-2)
	for position := 0; position < totalSlots; position++ {
		positions = append(positions, position)
	}
	for position := totalSlots - 2; position > 0; position-- {
		positions = append(positions, position)
	}

	return positions
}

func castBandForPosition(position int, slotWidth int) encounter.CastBand {
	bands := encounter.OrderedCastBands()
	if len(bands) == 0 {
		return ""
	}
	if position < 0 {
		return bands[0]
	}

	bandIndex := position / slotWidth
	if bandIndex >= len(bands) {
		bandIndex = len(bands) - 1
	}

	return bands[bandIndex]
}
