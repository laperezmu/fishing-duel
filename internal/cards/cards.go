package cards

import "pesca/internal/domain"

type Trigger int

const (
	TriggerOnDraw Trigger = iota
	TriggerOnPlayerWin
	TriggerOnFishWin
)

type EncounterModifier struct {
	Trigger       Trigger
	DistanceShift int
	DepthShift    int
}

func (modifier EncounterModifier) Applies(outcome domain.RoundOutcome) bool {
	switch modifier.Trigger {
	case TriggerOnPlayerWin:
		return outcome == domain.PlayerWin
	case TriggerOnFishWin:
		return outcome == domain.FishWin
	default:
		return outcome == domain.Draw
	}
}

type FishCard struct {
	Move               domain.Move
	EncounterModifiers []EncounterModifier
}

func NewFishCard(move domain.Move, encounterModifiers ...EncounterModifier) FishCard {
	return FishCard{
		Move:               move,
		EncounterModifiers: append([]EncounterModifier(nil), encounterModifiers...),
	}
}
