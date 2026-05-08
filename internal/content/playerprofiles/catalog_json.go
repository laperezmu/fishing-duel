package playerprofiles

import (
	"encoding/json"
	"fmt"
	"pesca/internal/cards"
	"pesca/internal/domain"
	"pesca/internal/player/playermoves"
)

type presetCatalogDocument struct {
	Presets []presetRecord `json:"presets"`
}

type presetRecord struct {
	ID                  string                  `json:"id"`
	Name                string                  `json:"name"`
	Description         string                  `json:"description"`
	Details             []string                `json:"details"`
	RecoveryDelayRounds int                     `json:"recovery_delay_rounds"`
	Cards               map[string][]cardRecord `json:"cards"`
}

type cardRecord struct {
	Name    string             `json:"name,omitempty"`
	Summary string             `json:"summary,omitempty"`
	Move    string             `json:"move"`
	Effects []cardEffectRecord `json:"effects,omitempty"`
}

type cardEffectRecord struct {
	Trigger                        string           `json:"trigger"`
	Type                           cards.EffectType `json:"effect_type,omitempty"`
	Priority                       int              `json:"priority,omitempty"`
	TargetMove                     string           `json:"target_move,omitempty"`
	DistanceShift                  int              `json:"distance_shift,omitempty"`
	DepthShift                     int              `json:"depth_shift,omitempty"`
	CaptureDistanceBonus           int              `json:"capture_distance_bonus,omitempty"`
	ExhaustionCaptureDistanceBonus int              `json:"exhaustion_capture_distance_bonus,omitempty"`
	SurfaceDepthBonus              int              `json:"surface_depth_bonus,omitempty"`
}

func LoadPresets(data []byte) ([]DeckPreset, error) {
	var document presetCatalogDocument
	if err := json.Unmarshal(data, &document); err != nil {
		return nil, fmt.Errorf("parse player presets catalog: %w", err)
	}

	presets := make([]DeckPreset, 0, len(document.Presets))
	seenIDs := make(map[string]struct{}, len(document.Presets))
	for _, record := range document.Presets {
		preset, err := record.toDomain()
		if err != nil {
			return nil, fmt.Errorf("player preset %s: %w", record.ID, err)
		}
		if _, exists := seenIDs[preset.ID]; exists {
			return nil, fmt.Errorf("duplicated player preset id %q", preset.ID)
		}
		seenIDs[preset.ID] = struct{}{}
		presets = append(presets, preset)
	}

	return presets, nil
}

func (record presetRecord) toDomain() (DeckPreset, error) {
	preset := DeckPreset{
		ID:          record.ID,
		Name:        record.Name,
		Description: record.Description,
		Details:     append([]string(nil), record.Details...),
		Config: playermoves.Config{
			InitialDecks:        make(map[domain.Move][]cards.PlayerCard, 3),
			RecoveryDelayRounds: record.RecoveryDelayRounds,
		},
	}

	for _, move := range []domain.Move{domain.Blue, domain.Red, domain.Yellow} {
		key := move.String()
		cardRecords, ok := record.Cards[key]
		if !ok {
			return DeckPreset{}, fmt.Errorf("cards for move %s are required", key)
		}
		builtCards := make([]cards.PlayerCard, 0, len(cardRecords))
		for _, cardRecord := range cardRecords {
			card, err := cardRecord.toDomain()
			if err != nil {
				return DeckPreset{}, err
			}
			if card.Move != move {
				return DeckPreset{}, fmt.Errorf("card move %s does not match deck move %s", card.Move, move)
			}
			builtCards = append(builtCards, card)
		}
		preset.Config.InitialDecks[move] = builtCards
	}

	return preset, preset.Config.Validate()
}

func (record cardRecord) toDomain() (cards.PlayerCard, error) {
	move, err := parseMove(record.Move)
	if err != nil {
		return cards.PlayerCard{}, err
	}
	effects := make([]cards.CardEffect, 0, len(record.Effects))
	for _, effectRecord := range record.Effects {
		effect, err := effectRecord.toDomain()
		if err != nil {
			return cards.PlayerCard{}, err
		}
		effects = append(effects, effect)
	}
	if record.Name == "" && record.Summary == "" && len(effects) == 0 {
		return cards.NewPlayerCard(move), nil
	}

	return cards.NewNamedPlayerCard(record.Name, record.Summary, move, effects...), nil
}

func (record cardEffectRecord) toDomain() (cards.CardEffect, error) {
	effect := cards.CardEffect{
		Type:                           record.Type,
		Priority:                       record.Priority,
		DistanceShift:                  record.DistanceShift,
		DepthShift:                     record.DepthShift,
		CaptureDistanceBonus:           record.CaptureDistanceBonus,
		ExhaustionCaptureDistanceBonus: record.ExhaustionCaptureDistanceBonus,
		SurfaceDepthBonus:              record.SurfaceDepthBonus,
	}
	if record.TargetMove != "" {
		targetMove, err := parseMove(record.TargetMove)
		if err != nil {
			return cards.CardEffect{}, fmt.Errorf("target move: %w", err)
		}
		effect.TargetMove = targetMove
	}
	trigger, err := parseTrigger(record.Trigger)
	if err != nil {
		return cards.CardEffect{}, err
	}
	effect.Trigger = trigger
	if err := effect.Validate(); err != nil {
		return cards.CardEffect{}, err
	}

	return effect.Normalize(), nil
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

func parseTrigger(value string) (cards.Trigger, error) {
	trigger, ok := triggerRecords[value]
	if !ok {
		return 0, fmt.Errorf("unknown trigger %s", value)
	}

	return trigger, nil
}

var triggerRecords = map[string]cards.Trigger{
	"on_draw":             cards.TriggerOnDraw,
	"on_card_used":        cards.TriggerOnCardUsed,
	"on_owner_win":        cards.TriggerOnOwnerWin,
	"on_owner_lose":       cards.TriggerOnOwnerLose,
	"on_round_draw":       cards.TriggerOnRoundDraw,
	"on_fish_splash":      cards.TriggerOnFishSplash,
	"on_discard":          cards.TriggerOnDiscard,
	"on_fish_reshuffle":   cards.TriggerOnFishReshuffle,
	"on_fish_exhausted":   cards.TriggerOnFishExhausted,
	"on_color_draw":       cards.TriggerOnColorDraw,
	"on_owner_color_win":  cards.TriggerOnOwnerColorWin,
	"on_owner_color_lose": cards.TriggerOnOwnerColorLose,
}
