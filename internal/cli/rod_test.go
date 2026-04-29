package cli

import (
	"bytes"
	"pesca/internal/content/attachmentpresets"
	"pesca/internal/content/rodpresets"
	"pesca/internal/player/rod"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChooseRodPreset(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(strings.NewReader("2\ns\n"), &out)

	preset, err := ui.ChooseRodPreset("Pesca: duelo contra el pez", sampleRodPresets())

	require.NoError(t, err)
	assert.Equal(t, "Versatil estandar", preset.Name)
	assert.Contains(t, out.String(), "Preset de cana")
	assert.Contains(t, out.String(), "Confirmar cana")
	assert.Contains(t, out.String(), "Apertura        : dist 5 | prof 3")
	assert.Contains(t, out.String(), "Track           : dist 5 | prof 3")
	assert.Contains(t, out.String(), clearSequence)
}

func TestChooseAttachmentPreset(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(strings.NewReader("2\ns\n"), &out)
	baseRod, err := rod.NewState(rod.DefaultConfig())
	require.NoError(t, err)

	preset, err := ui.ChooseAttachmentPreset("Pesca: duelo contra el pez", baseRod, sampleAttachmentPresets())

	require.NoError(t, err)
	assert.Equal(t, "Kit de fondo", preset.Name)
	assert.Contains(t, out.String(), "Preset de aditamentos")
	assert.Contains(t, out.String(), "Confirmar aditamentos")
	assert.Contains(t, out.String(), "Resultado final : apertura dist 4 | prof 4")
	assert.Contains(t, out.String(), "Track final     : dist 5 | prof 4")
	assert.Contains(t, out.String(), clearSequence)
}

func sampleRodPresets() []rodpresets.Preset {
	return rodpresets.DefaultPresets()
}

func sampleAttachmentPresets() []attachmentpresets.Preset {
	return attachmentpresets.DefaultPresets()
}
