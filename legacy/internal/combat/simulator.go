package combat

import "math/rand"

type CombatState struct {
	TrackPos           int
	Fish               FishProfile
	Deck               FishDeck
	FatigueCount       int
	BaitGuardAvailable bool
	RodMods            RodModifiers
}

func NewCombatState(config CombatConfig, rng *rand.Rand) CombatState {
	return CombatState{
		TrackPos:           config.InitialTrackPos,
		Fish:               config.Fish,
		Deck:               NewShuffledFishDeck(rng),
		BaitGuardAvailable: config.BaitGuardAvailable,
		RodMods:            config.RodMods,
	}
}

func NewCombatStateWithDeck(config CombatConfig, deck FishDeck) CombatState {
	return CombatState{
		TrackPos:           config.InitialTrackPos,
		Fish:               config.Fish,
		Deck:               deck.Clone(),
		BaitGuardAvailable: config.BaitGuardAvailable,
		RodMods:            config.RodMods,
	}
}

func (s CombatState) Clone() CombatState {
	s.Deck = s.Deck.Clone()
	return s
}

func (s CombatState) Observation() BeliefState {
	return BeliefState{
		TrackPos:           uint8(s.TrackPos),
		Fish:               s.Fish,
		Draw:               s.Deck.DrawCounts(),
		Discard:            s.Deck.DiscardCounts(),
		FatigueCount:       uint8(s.FatigueCount),
		BaitGuardAvailable: s.BaitGuardAvailable,
		RodMods:            s.RodMods,
	}
}

func RunCombat(initial CombatState, policy PlayerPolicy, rng *rand.Rand) CombatResult {
	state := initial.Clone()
	result := CombatResult{InitialTrackPos: state.TrackPos}

	for {
		observation := state.Observation()
		action := policy.ChooseAction(observation)
		result.ActionUsage[action]++

		fishCard, ok := state.Deck.Draw()
		if !ok {
			result.Outcome = OutcomeCapture
			result.TerminalReason = TerminalReasonTotalFatigue
			result.FinalTrackPos = state.TrackPos
			return result
		}

		result.Rounds++
		round := EvaluateRound(state.TrackPos, state.Fish, state.FatigueCount, state.BaitGuardAvailable, state.RodMods, fishCard, action)
		state.TrackPos = round.TrackPos
		if round.BaitGuardConsumed {
			state.BaitGuardAvailable = false
			result.BaitSaves++
		}
		if round.OffensiveApplied {
			result.OffensiveModTriggers++
		}

		if round.Outcome == OutcomeCapture || round.Outcome == OutcomeEscape {
			result.Outcome = round.Outcome
			result.TerminalReason = round.TerminalReason
			result.FinalTrackPos = round.TrackPos
			return result
		}

		state.Deck.DiscardCard(fishCard)
		if state.Deck.DrawLen() > 0 {
			continue
		}

		state.FatigueCount++
		result.Fatigues++
		fatigue := state.Deck.ReshuffleForFatigue(rng)
		if fatigue.TotalFatigue {
			result.Outcome = OutcomeCapture
			result.TerminalReason = TerminalReasonTotalFatigue
			result.FinalTrackPos = state.TrackPos
			return result
		}
	}
}
