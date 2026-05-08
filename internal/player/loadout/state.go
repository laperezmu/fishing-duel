package loadout

import (
	"fmt"
	"pesca/internal/content/habitats"
	"pesca/internal/encounter"
	"pesca/internal/player/rod"
)

type Attachment struct {
	ID                      string
	Name                    string
	Description             string
	OpeningDistanceModifier int
	OpeningDepthModifier    int
	TrackDistanceModifier   int
	TrackDepthModifier      int
	SplashBonusDistance     int
	HabitatTags             []habitats.Tag
}

type State struct {
	Rod         rod.State
	Attachments []Attachment
}

func NewState(playerRod rod.State, attachments []Attachment) (State, error) {
	clonedAttachments := make([]Attachment, 0, len(attachments))
	for _, attachment := range attachments {
		clonedAttachment := attachment
		clonedAttachment.HabitatTags = append([]habitats.Tag(nil), attachment.HabitatTags...)
		clonedAttachments = append(clonedAttachments, clonedAttachment)
	}

	state := State{
		Rod:         playerRod,
		Attachments: clonedAttachments,
	}
	if err := state.Validate(); err != nil {
		return State{}, err
	}

	return state, nil
}

func (state State) Validate() error {
	if err := state.Rod.Validate(); err != nil {
		return fmt.Errorf("rod: %w", err)
	}

	seenAttachmentIDs := make(map[string]struct{}, len(state.Attachments))
	for _, attachment := range state.Attachments {
		if attachment.ID == "" {
			return fmt.Errorf("attachment id is required")
		}
		if attachment.Name == "" {
			return fmt.Errorf("attachment name is required")
		}
		if _, exists := seenAttachmentIDs[attachment.ID]; exists {
			return fmt.Errorf("attachment ids must be unique")
		}
		seenAttachmentIDs[attachment.ID] = struct{}{}
		for _, habitatTag := range attachment.HabitatTags {
			if err := habitatTag.Validate(); err != nil {
				return err
			}
		}
	}

	if _, err := state.EffectiveRod(); err != nil {
		return fmt.Errorf("rod: %w", err)
	}

	return nil
}

func (state State) EffectiveRod() (rod.State, error) {
	effectiveConfig := rod.Config{
		OpeningMaxDistance:  state.Rod.OpeningMaxDistance,
		OpeningMaxDepth:     state.Rod.OpeningMaxDepth,
		TrackMaxDistance:    state.Rod.TrackMaxDistance,
		TrackMaxDepth:       state.Rod.TrackMaxDepth,
		SplashBonusDistance: state.Rod.SplashBonusDistance,
	}

	for _, attachment := range state.Attachments {
		effectiveConfig.OpeningMaxDistance += attachment.OpeningDistanceModifier
		effectiveConfig.OpeningMaxDepth += attachment.OpeningDepthModifier
		effectiveConfig.TrackMaxDistance += attachment.TrackDistanceModifier
		effectiveConfig.TrackMaxDepth += attachment.TrackDepthModifier
		effectiveConfig.SplashBonusDistance += attachment.SplashBonusDistance
	}

	return rod.NewState(effectiveConfig)
}

func (state State) SplashSuccessDistanceBonus() int {
	effectiveRod, err := state.EffectiveRod()
	if err != nil {
		return 0
	}

	return effectiveRod.SplashBonusDistance
}

func (state State) OpeningLimits() encounter.OpeningLimits {
	effectiveRod, err := state.EffectiveRod()
	if err != nil {
		return encounter.OpeningLimits{}
	}

	return encounter.OpeningLimits{
		MaxInitialDistance: effectiveRod.OpeningMaxDistance,
		MaxInitialDepth:    effectiveRod.OpeningMaxDepth,
	}
}

func (state State) TrackMaxDistance() int {
	effectiveRod, err := state.EffectiveRod()
	if err != nil {
		return 0
	}

	return effectiveRod.TrackMaxDistance
}

func (state State) TrackMaxDepth() int {
	effectiveRod, err := state.EffectiveRod()
	if err != nil {
		return 0
	}

	return effectiveRod.TrackMaxDepth
}

func (state State) HabitatTags() []habitats.Tag {
	tags := make([]habitats.Tag, 0, len(state.Attachments))
	seenTags := make(map[habitats.Tag]struct{}, len(state.Attachments))
	for _, attachment := range state.Attachments {
		for _, habitatTag := range attachment.HabitatTags {
			if _, exists := seenTags[habitatTag]; exists {
				continue
			}
			seenTags[habitatTag] = struct{}{}
			tags = append(tags, habitatTag)
		}
	}

	return tags
}
