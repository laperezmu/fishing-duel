package rodpresets

import (
	"encoding/json"
	"fmt"
	"pesca/internal/player/rod"
)

type CatalogDocument struct {
	Presets []PresetRecord `json:"presets"`
}

type PresetRecord struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Details     []string     `json:"details"`
	Config      ConfigRecord `json:"config"`
}

type ConfigRecord struct {
	OpeningMaxDistance  int `json:"opening_max_distance"`
	OpeningMaxDepth     int `json:"opening_max_depth"`
	TrackMaxDistance    int `json:"track_max_distance"`
	TrackMaxDepth       int `json:"track_max_depth"`
	SplashBonusDistance int `json:"splash_bonus_distance"`
}

func LoadPresets(data []byte) ([]Preset, error) {
	var document CatalogDocument
	if err := json.Unmarshal(data, &document); err != nil {
		return nil, fmt.Errorf("parse rod presets catalog: %w", err)
	}

	presets := make([]Preset, 0, len(document.Presets))
	seenIDs := make(map[string]struct{}, len(document.Presets))
	for _, record := range document.Presets {
		preset, err := record.toDomain()
		if err != nil {
			return nil, fmt.Errorf("rod preset %s: %w", record.ID, err)
		}
		if _, exists := seenIDs[preset.ID]; exists {
			return nil, fmt.Errorf("duplicated rod preset id %q", preset.ID)
		}
		seenIDs[preset.ID] = struct{}{}
		presets = append(presets, preset)
	}

	return presets, nil
}

func (record PresetRecord) toDomain() (Preset, error) {
	config := rod.Config{
		OpeningMaxDistance:  record.Config.OpeningMaxDistance,
		OpeningMaxDepth:     record.Config.OpeningMaxDepth,
		TrackMaxDistance:    record.Config.TrackMaxDistance,
		TrackMaxDepth:       record.Config.TrackMaxDepth,
		SplashBonusDistance: record.Config.SplashBonusDistance,
	}
	if config.SplashBonusDistance == 0 {
		config.SplashBonusDistance = rod.DefaultConfig().SplashBonusDistance
	}
	if config.TrackMaxDistance == 0 {
		config = rod.DefaultConfig()
		config.SplashBonusDistance = record.Config.SplashBonusDistance
	}

	return Preset{
		ID:          record.ID,
		Name:        record.Name,
		Description: record.Description,
		Details:     append([]string(nil), record.Details...),
		Config:      config,
	}, nil
}
