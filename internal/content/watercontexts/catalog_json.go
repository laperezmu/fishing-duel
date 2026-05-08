package watercontexts

import (
	"encoding/json"
	"fmt"
	"pesca/internal/content/waterpools"
	"pesca/internal/encounter"
)

type CatalogDocument struct {
	Presets    []PresetRecord `json:"presets"`
	PhaseOrder []string       `json:"phase_order"`
}

type PresetRecord struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Signals      []string       `json:"signals"`
	PoolTag      string         `json:"pool_tag"`
	Distances    map[string]int `json:"distances"`
	InitialDepth int            `json:"initial_depth"`
}

func LoadPresets(data []byte) ([]Preset, map[ID]int, error) {
	var document CatalogDocument
	if err := json.Unmarshal(data, &document); err != nil {
		return nil, nil, fmt.Errorf("parse water contexts catalog: %w", err)
	}

	presets := make([]Preset, 0, len(document.Presets))
	seenIDs := make(map[string]struct{}, len(document.Presets))
	for _, record := range document.Presets {
		preset, err := record.toDomain()
		if err != nil {
			return nil, nil, fmt.Errorf("water context %s: %w", record.ID, err)
		}
		if _, exists := seenIDs[string(preset.ID)]; exists {
			return nil, nil, fmt.Errorf("duplicated water context id %q", preset.ID)
		}
		seenIDs[string(preset.ID)] = struct{}{}
		presets = append(presets, preset)
	}

	phaseIndex := make(map[ID]int, len(document.PhaseOrder))
	for i, idStr := range document.PhaseOrder {
		phaseIndex[ID(idStr)] = i
	}

	return presets, phaseIndex, nil
}

func (record PresetRecord) toDomain() (Preset, error) {
	poolTag := waterpools.ID(record.PoolTag)
	if err := poolTag.Validate(); err != nil {
		return Preset{}, fmt.Errorf("pool tag: %w", err)
	}

	distances := make(map[encounter.CastBand]int, len(record.Distances))
	for bandStr, distance := range record.Distances {
		band, err := parseCastBand(bandStr)
		if err != nil {
			return Preset{}, fmt.Errorf("cast band %s: %w", bandStr, err)
		}
		distances[band] = distance
	}

	return Preset{
		ID:           ID(record.ID),
		Name:         record.Name,
		Description:  record.Description,
		Signals:      append([]string(nil), record.Signals...),
		PoolTag:      poolTag,
		Distances:    distances,
		InitialDepth: record.InitialDepth,
	}, nil
}

func parseCastBand(value string) (encounter.CastBand, error) {
	switch value {
	case "very_short":
		return encounter.CastBandVeryShort, nil
	case "short":
		return encounter.CastBandShort, nil
	case "medium":
		return encounter.CastBandMedium, nil
	case "long":
		return encounter.CastBandLong, nil
	case "very_long":
		return encounter.CastBandVeryLong, nil
	default:
		return "", fmt.Errorf("unknown cast band %s", value)
	}
}

func castBandToString(band encounter.CastBand) string {
	switch band {
	case encounter.CastBandVeryShort:
		return "very_short"
	case encounter.CastBandShort:
		return "short"
	case encounter.CastBandMedium:
		return "medium"
	case encounter.CastBandLong:
		return "long"
	case encounter.CastBandVeryLong:
		return "very_long"
	default:
		return ""
	}
}
