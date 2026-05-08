package encounter

import "pesca/internal/cards"

func ApplyDelta(state *State, delta Delta) {
	state.Distance += delta.DistanceShift
	if state.Splash == nil {
		state.LastEvent = Event{}
	}

	if delta.DepthShift == 0 {
		return
	}

	targetDepth := state.Depth + delta.DepthShift
	if targetDepth < state.Config.SurfaceDepth {
		state.Depth = state.Config.SurfaceDepth
		state.LastEvent = Event{
			Kind:    EventKindSplash,
			Escaped: false,
		}
		profile := state.Config.SplashProfile.WithDefaults()
		state.Splash = &SplashState{
			TotalJumps:    profile.JumpCount,
			ResolvedJumps: 0,
			TimeLimit:     profile.TimeLimit,
		}

		return
	}

	state.Depth = targetDepth
}

func ApplyMovementEffects(state *State, effects []cards.CardEffect) {
	for _, effect := range effects {
		effect = effect.Normalize()
		if effect.DistanceShift == 0 && effect.DepthShift == 0 {
			continue
		}
		ApplyDelta(state, Delta{
			DistanceShift: effect.DistanceShift,
			DepthShift:    effect.DepthShift,
		})
	}
}

func ApplySplashResolution(state *State, resolution SplashResolution) {
	if state.Splash == nil {
		return
	}

	if resolution.Escaped {
		state.LastEvent = Event{Kind: EventKindSplash, Escaped: true}
		state.Status = StatusEscaped
		state.EndReason = EndReasonSplashEscape
		state.Splash = nil
		return
	}

	if resolution.SuccessfulJumps > 0 {
		state.Splash.ResolvedJumps += resolution.SuccessfulJumps
		state.Distance -= resolution.DistanceRewardApplied
		if state.Distance < 0 {
			state.Distance = 0
		}
	}

	state.LastEvent = Event{Kind: EventKindSplash, Escaped: false}
	if !state.Splash.Pending() {
		state.Splash = nil
	}
}
