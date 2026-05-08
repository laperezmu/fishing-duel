package attachmentpresets

import (
	"encoding/json"
	"fmt"
	"pesca/internal/content/habitats"
	"pesca/internal/player/loadout"
)

type CatalogDocument struct {
	Presets []PresetRecord `json:"presets"`
}

type PresetRecord struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Details     []string           `json:"details"`
	Attachments []AttachmentRecord `json:"attachments,omitempty"`
}

type AttachmentRecord struct {
	ID                      string   `json:"id"`
	Name                    string   `json:"name"`
	Description             string   `json:"description"`
	OpeningDistanceModifier int      `json:"opening_distance_modifier,omitempty"`
	OpeningDepthModifier    int      `json:"opening_depth_modifier,omitempty"`
	TrackDistanceModifier   int      `json:"track_distance_modifier,omitempty"`
	TrackDepthModifier      int      `json:"track_depth_modifier,omitempty"`
	HabitatTags             []string `json:"habitat_tags,omitempty"`
}

func LoadPresets(data []byte) ([]Preset, error) {
	var document CatalogDocument
	if err := json.Unmarshal(data, &document); err != nil {
		return nil, fmt.Errorf("parse attachment presets catalog: %w", err)
	}

	presets := make([]Preset, 0, len(document.Presets))
	seenIDs := make(map[string]struct{}, len(document.Presets))
	for _, record := range document.Presets {
		preset, err := record.toDomain()
		if err != nil {
			return nil, fmt.Errorf("attachment preset %s: %w", record.ID, err)
		}
		if _, exists := seenIDs[preset.ID]; exists {
			return nil, fmt.Errorf("duplicated attachment preset id %q", preset.ID)
		}
		seenIDs[preset.ID] = struct{}{}
		presets = append(presets, preset)
	}

	return presets, nil
}

func (record PresetRecord) toDomain() (Preset, error) {
	attachments := make([]loadout.Attachment, 0, len(record.Attachments))
	for _, attRecord := range record.Attachments {
		habitatTags := make([]habitats.Tag, 0, len(attRecord.HabitatTags))
		for _, tagStr := range attRecord.HabitatTags {
			if tag, ok := habitatRecords[tagStr]; ok {
				habitatTags = append(habitatTags, tag)
			}
		}
		attachments = append(attachments, loadout.Attachment{
			ID:                      attRecord.ID,
			Name:                    attRecord.Name,
			Description:             attRecord.Description,
			OpeningDistanceModifier: attRecord.OpeningDistanceModifier,
			OpeningDepthModifier:    attRecord.OpeningDepthModifier,
			TrackDistanceModifier:   attRecord.TrackDistanceModifier,
			TrackDepthModifier:      attRecord.TrackDepthModifier,
			HabitatTags:             habitatTags,
		})
	}

	return Preset{
		ID:          record.ID,
		Name:        record.Name,
		Description: record.Description,
		Details:     append([]string(nil), record.Details...),
		Attachments: attachments,
	}, nil
}

var habitatRecords = map[string]habitats.Tag{
	"bottom":     habitats.Bottom,
	"channel":    habitats.Channel,
	"open-water": habitats.OpenWater,
	"surface":    habitats.Surface,
	"weed":       habitats.Weed,
	"rock":       habitats.Rock,
}
