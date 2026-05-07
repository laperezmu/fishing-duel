package fishprofiles

import (
	"encoding/json"
	"fmt"
	"pesca/internal/cards"
	"pesca/internal/content/habitats"
	"pesca/internal/content/waterpools"
	"pesca/internal/domain"
)

type profileCatalogDocument struct {
	Profiles []profileRecord `json:"profiles"`
}

type poolCatalogDocument struct {
	Pools []poolRecord `json:"pools"`
}

type profileRecord struct {
	ID            ProfileID           `json:"id"`
	ArchetypeID   ArchetypeID         `json:"archetype_id"`
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	Details       []string            `json:"details"`
	Appearance    appearanceRecord    `json:"appearance"`
	Cards         []cardPatternRecord `json:"cards"`
	CardsToRemove int                 `json:"cards_to_remove"`
	Shuffle       bool                `json:"shuffle"`
}

type appearanceRecord struct {
	WaterPoolTags       []waterpools.ID `json:"water_pool_tags"`
	MinInitialDistance  int             `json:"min_initial_distance"`
	MaxInitialDistance  int             `json:"max_initial_distance"`
	MinInitialDepth     int             `json:"min_initial_depth"`
	MaxInitialDepth     int             `json:"max_initial_depth"`
	RequiredHabitatTags []habitats.Tag  `json:"required_habitat_tags"`
}

type cardPatternRecord struct {
	Name              string                  `json:"name,omitempty"`
	Summary           string                  `json:"summary,omitempty"`
	Move              string                  `json:"move"`
	Effects           []cardEffectRecord      `json:"effects,omitempty"`
	DiscardVisibility cards.DiscardVisibility `json:"discard_visibility,omitempty"`
}

type cardEffectRecord struct {
	Trigger                        string `json:"trigger"`
	DistanceShift                  int    `json:"distance_shift,omitempty"`
	DepthShift                     int    `json:"depth_shift,omitempty"`
	CaptureDistanceBonus           int    `json:"capture_distance_bonus,omitempty"`
	ExhaustionCaptureDistanceBonus int    `json:"exhaustion_capture_distance_bonus,omitempty"`
	SurfaceDepthBonus              int    `json:"surface_depth_bonus,omitempty"`
}

type poolRecord struct {
	ID          PoolID            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	ProfileIDs  []ProfileID       `json:"profile_ids,omitempty"`
	Entries     []poolEntryRecord `json:"entries,omitempty"`
}

type poolEntryRecord struct {
	ProfileID ProfileID `json:"profile_id"`
	Weight    *int      `json:"weight,omitempty"`
}

func LoadCatalog(profileData []byte, poolData []byte) (Catalog, error) {
	var profileDocument profileCatalogDocument
	if err := json.Unmarshal(profileData, &profileDocument); err != nil {
		return Catalog{}, fmt.Errorf("parse fish profiles catalog: %w", err)
	}
	var poolDocument poolCatalogDocument
	if err := json.Unmarshal(poolData, &poolDocument); err != nil {
		return Catalog{}, fmt.Errorf("parse fish pools catalog: %w", err)
	}

	profiles := make([]Profile, 0, len(profileDocument.Profiles))
	for _, profileRecord := range profileDocument.Profiles {
		profile, err := profileRecord.toDomain()
		if err != nil {
			return Catalog{}, fmt.Errorf("fish profile %s: %w", profileRecord.ID, err)
		}
		profiles = append(profiles, profile)
	}

	pools := make([]Pool, 0, len(poolDocument.Pools))
	for _, poolRecord := range poolDocument.Pools {
		pools = append(pools, poolRecord.toDomain())
	}

	return NewCatalog(profiles, pools)
}

func (record profileRecord) toDomain() (Profile, error) {
	profile := Profile{
		ID:          record.ID,
		ArchetypeID: record.ArchetypeID,
		Name:        record.Name,
		Description: record.Description,
		Details:     append([]string(nil), record.Details...),
		Appearance: Appearance{
			WaterPoolTags:       append([]waterpools.ID(nil), record.Appearance.WaterPoolTags...),
			MinInitialDistance:  record.Appearance.MinInitialDistance,
			MaxInitialDistance:  record.Appearance.MaxInitialDistance,
			MinInitialDepth:     record.Appearance.MinInitialDepth,
			MaxInitialDepth:     record.Appearance.MaxInitialDepth,
			RequiredHabitatTags: append([]habitats.Tag(nil), record.Appearance.RequiredHabitatTags...),
		},
		CardsToRemove: record.CardsToRemove,
		Shuffle:       record.Shuffle,
	}
	for _, cardRecord := range record.Cards {
		move, err := parseMove(cardRecord.Move)
		if err != nil {
			return Profile{}, err
		}
		effects := make([]cards.CardEffect, 0, len(cardRecord.Effects))
		for _, effectRecord := range cardRecord.Effects {
			effect, err := effectRecord.toDomain()
			if err != nil {
				return Profile{}, err
			}
			effects = append(effects, effect)
		}
		profile.Cards = append(profile.Cards, CardPattern{
			Name:              cardRecord.Name,
			Summary:           cardRecord.Summary,
			Move:              move,
			Effects:           effects,
			DiscardVisibility: cardRecord.DiscardVisibility,
		})
	}

	return profile, profile.Validate()
}

func (record poolRecord) toDomain() Pool {
	entries := make([]PoolEntry, 0, len(record.Entries)+len(record.ProfileIDs))
	if len(record.Entries) > 0 {
		for _, entry := range record.Entries {
			weight := 1
			if entry.Weight != nil {
				weight = *entry.Weight
			}
			entries = append(entries, PoolEntry{ProfileID: entry.ProfileID, Weight: weight})
		}
	} else {
		for _, profileID := range record.ProfileIDs {
			entries = append(entries, PoolEntry{ProfileID: profileID, Weight: 1})
		}
	}

	return Pool{
		ID:          record.ID,
		Name:        record.Name,
		Description: record.Description,
		Entries:     entries,
	}
}

func parseMove(value string) (domain.Move, error) {
	switch value {
	case domain.Blue.String():
		return domain.Blue, nil
	case domain.Red.String():
		return domain.Red, nil
	case domain.Yellow.String():
		return domain.Yellow, nil
	default:
		return 0, fmt.Errorf("unknown move %s", value)
	}
}

func (record cardEffectRecord) toDomain() (cards.CardEffect, error) {
	effect := cards.CardEffect{
		DistanceShift:                  record.DistanceShift,
		DepthShift:                     record.DepthShift,
		CaptureDistanceBonus:           record.CaptureDistanceBonus,
		ExhaustionCaptureDistanceBonus: record.ExhaustionCaptureDistanceBonus,
		SurfaceDepthBonus:              record.SurfaceDepthBonus,
	}

	switch record.Trigger {
	case "on_draw":
		effect.Trigger = cards.TriggerOnDraw
	case "on_owner_win":
		effect.Trigger = cards.TriggerOnOwnerWin
	case "on_owner_lose":
		effect.Trigger = cards.TriggerOnOwnerLose
	case "on_round_draw":
		effect.Trigger = cards.TriggerOnRoundDraw
	default:
		return cards.CardEffect{}, fmt.Errorf("unknown trigger %s", record.Trigger)
	}

	return effect, nil
}
