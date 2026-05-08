package match

type ResolutionTraceSnapshot struct {
	Before          StatusSnapshot
	After           StatusSnapshot
	ResolvedEffects []ResolvedEffectState
}

func NewResolutionTraceSnapshot(before State, after State, resolved []ResolvedEffectState) ResolutionTraceSnapshot {
	return ResolutionTraceSnapshot{
		Before:          NewStatusSnapshot(before),
		After:           NewStatusSnapshot(after),
		ResolvedEffects: append([]ResolvedEffectState(nil), resolved...),
	}
}
