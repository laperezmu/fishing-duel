package combat

import "math"

type RoundEvaluation struct {
	Result            int
	TrackPos          int
	Outcome           Outcome
	TerminalReason    TerminalReason
	BaitGuardConsumed bool
	OffensiveApplied  bool
}

func ActionColor(action Action) Color {
	switch action {
	case Forzar:
		return ColorRed
	case Tensar:
		return ColorBlue
	case Soltar:
		return ColorYellow
	default:
		return ColorNone
	}
}

func FamilyColor(f FishFamily) Color {
	switch f {
	case Embiste:
		return ColorRed
	case Aguante:
		return ColorBlue
	case Quiebre:
		return ColorYellow
	default:
		return ColorNone
	}
}

func BaseResult(fishCard FishFamily, action Action) int {
	switch {
	case fishCard == Embiste && action == Forzar:
		return 0
	case fishCard == Embiste && action == Tensar:
		return -1
	case fishCard == Embiste && action == Soltar:
		return 1
	case fishCard == Aguante && action == Forzar:
		return 1
	case fishCard == Aguante && action == Tensar:
		return 0
	case fishCard == Aguante && action == Soltar:
		return -1
	case fishCard == Quiebre && action == Forzar:
		return -1
	case fishCard == Quiebre && action == Tensar:
		return 1
	case fishCard == Quiebre && action == Soltar:
		return 0
	default:
		return 0
	}
}

func ApplyArchetype(result int, fish FishProfile, fishCard FishFamily, action Action) int {
	if result != 0 {
		return 0
	}
	color := FamilyColor(fishCard)
	if fish.HasColor(color) && ActionColor(action) == color {
		return -1
	}
	return 0
}

func ApplyFatigueModifier(baseResult int, fatigueCount int) int {
	if baseResult == 0 || fatigueCount < 1 {
		return 0
	}
	return int(math.Copysign(1, float64(baseResult)))
}

func ApplyRodModifier(baseResult int, action Action, mods RodModifiers) (modifier int, offensiveApplied bool) {
	color := ActionColor(action)
	if baseResult == 0 && mods.HasOffensive(color) {
		return 1, true
	}
	return 0, false
}

func EvaluateRound(trackPos int, fish FishProfile, fatigueCount int, baitGuardAvailable bool, mods RodModifiers, fishCard FishFamily, action Action) RoundEvaluation {
	base := BaseResult(fishCard, action)
	archetypeModifier := ApplyArchetype(base, fish, fishCard, action)
	fatigueModifier := ApplyFatigueModifier(base, fatigueCount)
	rodModifier, offensiveApplied := ApplyRodModifier(base, action, mods)
	result := base + archetypeModifier + fatigueModifier + rodModifier

	nextTrack := trackPos - result
	if nextTrack < 1 {
		return RoundEvaluation{
			Result:           result,
			TrackPos:         nextTrack,
			Outcome:          OutcomeCapture,
			TerminalReason:   TerminalReasonTrackCapture,
			OffensiveApplied: offensiveApplied,
		}
	}

	if nextTrack > 5 {
		if baitGuardAvailable {
			return RoundEvaluation{
				Result:            result,
				TrackPos:          5,
				Outcome:           OutcomeOngoing,
				BaitGuardConsumed: true,
				OffensiveApplied:  offensiveApplied,
			}
		}
		return RoundEvaluation{
			Result:           result,
			TrackPos:         nextTrack,
			Outcome:          OutcomeEscape,
			TerminalReason:   TerminalReasonTrackEscape,
			OffensiveApplied: offensiveApplied,
		}
	}

	return RoundEvaluation{
		Result:           result,
		TrackPos:         nextTrack,
		Outcome:          OutcomeOngoing,
		OffensiveApplied: offensiveApplied,
	}
}
