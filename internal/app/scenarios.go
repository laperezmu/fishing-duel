package app

import (
	"pesca/internal/content/fishprofiles"
	"pesca/internal/content/watercontexts"
)

func DefaultSandboxScenarios() []SandboxScenario {
	seedSurface := int64(11)
	seedTempo := int64(23)
	seedLegacy := int64(31)

	return []SandboxScenario{
		{
			ID:                 "surface-control-shoreline",
			Name:               "Control de superficie en ensenada",
			Description:        "Replay base para validar agua cercana y setup reproducible sin prompts.",
			Seed:               &seedSurface,
			PlayerDeckPresetID: "classic",
			FishPresetID:       fishprofiles.ProfileID("classic"),
			WaterContextID:     watercontexts.ShorelineCove,
		},
		{
			ID:                 "draw-tempo-open-channel",
			Name:               "Tempo de apertura en canal",
			Description:        "Replay de QA para validar efectos de apertura con cast medio.",
			Seed:               &seedTempo,
			PlayerDeckPresetID: "hooked-opening",
			FishPresetID:       fishprofiles.ProfileID("hooked-opening"),
			WaterContextID:     watercontexts.OpenChannel,
		},
		{
			ID:                 "legacy-mixed-current",
			Name:               "Regresion legacy en corriente mixta",
			Description:        "Replay para convivencias entre efectos legacy y respuestas por outcome.",
			Seed:               &seedLegacy,
			PlayerDeckPresetID: "mixed-current",
			FishPresetID:       fishprofiles.ProfileID("mixed-current"),
			WaterContextID:     watercontexts.BrokenCurrent,
		},
	}
}
