package anglerprofiles

import (
	"fmt"
	"pesca/internal/content/attachmentpresets"
	"pesca/internal/content/playerprofiles"
	"pesca/internal/content/rodpresets"
	"pesca/internal/player/loadout"
)

type Profile struct {
	ID                 string
	Name               string
	Description        string
	Details            []string
	DeckPresetID       string
	RodPresetID        string
	AttachmentPresetID string
	StartingThread     int
	UnlockedByDefault  bool
	UnlockID           string
}

type ResolvedStart struct {
	Profile          Profile
	DeckPreset       playerprofiles.DeckPreset
	RodPreset        rodpresets.Preset
	AttachmentPreset attachmentpresets.Preset
	Loadout          loadout.State
	StartingThread   int
}

func (profile Profile) Validate() error {
	if profile.ID == "" {
		return fmt.Errorf("angler profile id is required")
	}
	if profile.Name == "" {
		return fmt.Errorf("angler profile name is required")
	}
	if profile.DeckPresetID == "" {
		return fmt.Errorf("angler profile deck preset id is required")
	}
	if profile.RodPresetID == "" {
		return fmt.Errorf("angler profile rod preset id is required")
	}
	if profile.AttachmentPresetID == "" {
		return fmt.Errorf("angler profile attachment preset id is required")
	}
	if profile.StartingThread <= 0 {
		return fmt.Errorf("angler profile starting thread must be greater than 0")
	}

	return nil
}

func DefaultProfiles() []Profile {
	return []Profile{
		{
			ID:                 "coastal-specialist",
			Name:               "Especialista costero",
			Description:        "Apuesta por cierres cercanos y lectura rapida del agua, con menos hilo para errores largos.",
			Details:            []string{"Baraja tactica de apertura preparada.", "Rod de control costero.", "Sin aditamentos para cierres rapidos.", "Hilo inicial: 2."},
			DeckPresetID:       "hooked-opening",
			RodPresetID:        "coastal-control",
			AttachmentPresetID: "no-attachments",
			StartingThread:     2,
			UnlockedByDefault:  true,
		},
		{
			ID:                 "deep-angler",
			Name:               "Pescador de fondo",
			Description:        "Abre mejor pescas profundas y empuja habitats de fondo, a costa de una run mas tensa.",
			Details:            []string{"Baraja de respuesta vertical.", "Rod de presion de fondo.", "Kit de fondo para profundizar la apertura.", "Hilo inicial: 2."},
			DeckPresetID:       "vertical-pressure",
			RodPresetID:        "bottom-pressure",
			AttachmentPresetID: "bottom-kit",
			StartingThread:     2,
			UnlockedByDefault:  true,
		},
		{
			ID:                 "steady-handler",
			Name:               "Patron sereno",
			Description:        "Prioriza estabilidad general y mas hilo inicial para sostener una run larga.",
			Details:            []string{"Baraja clasica como base estable.", "Rod versatil estandar.", "Kit de estabilidad para aguantar el track.", "Hilo inicial: 4."},
			DeckPresetID:       "classic",
			RodPresetID:        "versatile-standard",
			AttachmentPresetID: "stability-kit",
			StartingThread:     4,
			UnlockedByDefault:  true,
		},
		{
			ID:                 "storm-reader",
			Name:               "Lector de tormentas",
			Description:        "Perfil bloqueado de ejemplo para futuras runs mas agresivas.",
			Details:            []string{"Reservado para desbloqueos meta futuros."},
			DeckPresetID:       "mixed-current",
			RodPresetID:        "versatile-standard",
			AttachmentPresetID: "long-cast-kit",
			StartingThread:     3,
			UnlockedByDefault:  false,
			UnlockID:           "unlock-storm-reader",
		},
	}
}

func DefaultUnlockedProfiles() []Profile {
	profiles := DefaultProfiles()
	unlocked := make([]Profile, 0, len(profiles))
	for _, profile := range profiles {
		if profile.UnlockedByDefault {
			unlocked = append(unlocked, profile)
		}
	}

	return unlocked
}

func ResolveStart(profile Profile) (ResolvedStart, error) {
	if err := profile.Validate(); err != nil {
		return ResolvedStart{}, err
	}

	deckPreset, err := resolveDeckPreset(profile.DeckPresetID)
	if err != nil {
		return ResolvedStart{}, err
	}
	rodPreset, err := resolveRodPreset(profile.RodPresetID)
	if err != nil {
		return ResolvedStart{}, err
	}
	attachmentPreset, err := resolveAttachmentPreset(profile.AttachmentPresetID)
	if err != nil {
		return ResolvedStart{}, err
	}
	resolvedLoadout, err := rodPreset.BuildLoadoutWithAttachments(attachmentPreset.BuildAttachments())
	if err != nil {
		return ResolvedStart{}, fmt.Errorf("build angler profile loadout: %w", err)
	}

	return ResolvedStart{
		Profile:          profile,
		DeckPreset:       deckPreset,
		RodPreset:        rodPreset,
		AttachmentPreset: attachmentPreset,
		Loadout:          resolvedLoadout,
		StartingThread:   profile.StartingThread,
	}, nil
}

func resolveDeckPreset(id string) (playerprofiles.DeckPreset, error) {
	for _, preset := range playerprofiles.DefaultPresets() {
		if preset.ID == id {
			return preset, nil
		}
	}

	return playerprofiles.DeckPreset{}, fmt.Errorf("unknown player deck preset %q", id)
}

func resolveRodPreset(id string) (rodpresets.Preset, error) {
	for _, preset := range rodpresets.DefaultPresets() {
		if preset.ID == id {
			return preset, nil
		}
	}

	return rodpresets.Preset{}, fmt.Errorf("unknown rod preset %q", id)
}

func resolveAttachmentPreset(id string) (attachmentpresets.Preset, error) {
	for _, preset := range attachmentpresets.DefaultPresets() {
		if preset.ID == id {
			return preset, nil
		}
	}

	return attachmentpresets.Preset{}, fmt.Errorf("unknown attachment preset %q", id)
}
