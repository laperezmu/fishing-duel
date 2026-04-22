package encounter

type SplashEscapeDecider interface {
	ShouldEscape(chance float64) bool
}

func ApplyDelta(state *State, delta Delta, splashEscapeDecider SplashEscapeDecider) {
	state.Distance += delta.DistanceShift
	state.LastEvent = Event{}

	if delta.DepthShift == 0 {
		return
	}

	targetDepth := state.Depth + delta.DepthShift
	if targetDepth < state.Config.SurfaceDepth {
		state.Depth = state.Config.SurfaceDepth
		escaped := false
		if splashEscapeDecider != nil {
			escaped = splashEscapeDecider.ShouldEscape(state.Config.SplashEscapeChance)
		}

		state.LastEvent = Event{
			Kind:    EventKindSplash,
			Escaped: escaped,
		}
		if escaped {
			state.Status = StatusEscaped
			state.EndReason = EndReasonSplashEscape
		}

		return
	}

	state.Depth = targetDepth
}
