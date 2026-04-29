package attachmentpresets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildAttachmentsClonesSlices(t *testing.T) {
	preset := DefaultPresets()[1]
	attachments := preset.BuildAttachments()
	require.Len(t, attachments, 1)

	attachments[0].HabitatTags[0] = "changed"

	assert.Equal(t, "bottom", preset.Attachments[0].HabitatTags[0])
}

func TestDefaultPresets(t *testing.T) {
	presets := DefaultPresets()

	require.Len(t, presets, 4)
	assert.Equal(t, "Sin aditamentos", presets[0].Name)
	assert.Empty(t, presets[0].Attachments)
	assert.Equal(t, "Kit de fondo", presets[1].Name)
	assert.Equal(t, -1, presets[1].Attachments[0].OpeningDistanceModifier)
	assert.Equal(t, "Kit de estabilidad", presets[3].Name)
	assert.Equal(t, 1, presets[3].Attachments[0].TrackDepthModifier)
}
