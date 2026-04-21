package combat

import (
	"math"
	"math/rand"
)

type PlayerPolicy interface {
	ChooseAction(state BeliefState) Action
	Name() string
}

type RandomPolicy struct {
	rng *rand.Rand
}

func NewRandomPolicy(rng *rand.Rand) *RandomPolicy {
	return &RandomPolicy{rng: rng}
}

func (p *RandomPolicy) ChooseAction(BeliefState) Action {
	return allActions[p.rng.Intn(len(allActions))]
}

func (p *RandomPolicy) Name() string {
	return "random"
}

type HeuristicPolicy struct{}

func NewHeuristicPolicy() *HeuristicPolicy {
	return &HeuristicPolicy{}
}

func (p *HeuristicPolicy) ChooseAction(state BeliefState) Action {
	total := float64(state.Draw.Total())
	bestAction := Forzar
	bestScore := math.Inf(-1)

	for _, action := range allActions {
		score := 0.0
		for family := FishFamily(0); family < fishFamilyCount; family++ {
			count := state.Draw.Count(family)
			if count == 0 {
				continue
			}
			prob := float64(count) / total
			round := EvaluateRound(int(state.TrackPos), state.Fish, int(state.FatigueCount), state.BaitGuardAvailable, state.RodMods, family, action)
			switch round.Outcome {
			case OutcomeCapture:
				score += prob * 100
			case OutcomeEscape:
				score -= prob * 100
			default:
				score += prob * float64(6-round.TrackPos)
			}
		}
		if score > bestScore {
			bestScore = score
			bestAction = action
		}
	}

	return bestAction
}

func (p *HeuristicPolicy) Name() string {
	return "heuristic"
}

type FixedPolicy struct {
	action Action
}

func NewFixedPolicy(action Action) *FixedPolicy {
	return &FixedPolicy{action: action}
}

func (p *FixedPolicy) ChooseAction(BeliefState) Action {
	return p.action
}

func (p *FixedPolicy) Name() string {
	return p.action.String()
}

type OptimalPolicy struct {
	cache map[BeliefState]ActionValues
}

func NewOptimalPolicy() *OptimalPolicy {
	return &OptimalPolicy{cache: make(map[BeliefState]ActionValues)}
}

func (p *OptimalPolicy) Name() string {
	return "optimal"
}

func (p *OptimalPolicy) ChooseAction(state BeliefState) Action {
	return p.ActionValues(state).BestAction
}

func (p *OptimalPolicy) CaptureProbability(state BeliefState) float64 {
	return p.ActionValues(state).BestValue
}

func (p *OptimalPolicy) ActionValues(state BeliefState) ActionValues {
	if cached, ok := p.cache[state]; ok {
		return cached
	}

	if state.TrackPos < 1 {
		result := ActionValues{BestAction: Forzar, BestValue: 1}
		for i := range result.Values {
			result.Values[i] = 1
		}
		p.cache[state] = result
		return result
	}

	if state.TrackPos > 5 || state.Draw.Total() == 0 {
		result := ActionValues{BestAction: Forzar, BestValue: 0}
		p.cache[state] = result
		return result
	}

	total := float64(state.Draw.Total())
	bestAction := Forzar
	bestValue := math.Inf(-1)
	var values [actionCount]float64

	for _, action := range allActions {
		value := 0.0
		for family := FishFamily(0); family < fishFamilyCount; family++ {
			count := state.Draw.Count(family)
			if count == 0 {
				continue
			}
			prob := float64(count) / total
			transition := p.advanceState(state, action, family)
			switch transition.outcome {
			case OutcomeCapture:
				value += prob
			case OutcomeEscape:
				// no-op
			default:
				for _, next := range transition.nextStates {
					value += prob * next.Probability * p.CaptureProbability(next.State)
				}
			}
		}
		values[action] = value
		if value > bestValue {
			bestValue = value
			bestAction = action
		}
	}

	result := ActionValues{Values: values, BestAction: bestAction, BestValue: bestValue}
	p.cache[state] = result
	return result
}

type beliefTransition struct {
	outcome    Outcome
	nextStates []weightedBeliefState
}

func (p *OptimalPolicy) advanceState(state BeliefState, action Action, fishCard FishFamily) beliefTransition {
	round := EvaluateRound(int(state.TrackPos), state.Fish, int(state.FatigueCount), state.BaitGuardAvailable, state.RodMods, fishCard, action)
	if round.Outcome == OutcomeCapture {
		return beliefTransition{outcome: OutcomeCapture}
	}
	if round.Outcome == OutcomeEscape {
		return beliefTransition{outcome: OutcomeEscape}
	}

	nextDraw := state.Draw
	nextDraw.Decrement(fishCard)
	nextDiscard := state.Discard
	nextDiscard.Increment(fishCard)
	nextState := BeliefState{
		TrackPos:           uint8(round.TrackPos),
		Fish:               state.Fish,
		Draw:               nextDraw,
		Discard:            nextDiscard,
		FatigueCount:       state.FatigueCount,
		BaitGuardAvailable: state.BaitGuardAvailable && !round.BaitGuardConsumed,
		RodMods:            state.RodMods,
	}

	if nextState.Draw.Total() > 0 {
		return beliefTransition{outcome: OutcomeOngoing, nextStates: []weightedBeliefState{{State: nextState, Probability: 1}}}
	}

	nextState.FatigueCount++
	if nextState.Discard.Total() <= 3 {
		return beliefTransition{outcome: OutcomeCapture}
	}

	trimmedStates := enumerateFatigueStates(nextState)
	return beliefTransition{outcome: OutcomeOngoing, nextStates: trimmedStates}
}

func enumerateFatigueStates(state BeliefState) []weightedBeliefState {
	total := state.Discard.Total()
	denominator := combination(total, 3)
	states := make([]weightedBeliefState, 0, 10)

	for embiste := 0; embiste <= min(3, int(state.Discard[Embiste])); embiste++ {
		for aguante := 0; aguante <= min(3-embiste, int(state.Discard[Aguante])); aguante++ {
			quiebre := 3 - embiste - aguante
			if quiebre < 0 || quiebre > int(state.Discard[Quiebre]) {
				continue
			}
			removed := NewFamilyCounts(uint8(embiste), uint8(aguante), uint8(quiebre))
			weight := combination(int(state.Discard[Embiste]), embiste) * combination(int(state.Discard[Aguante]), aguante) * combination(int(state.Discard[Quiebre]), quiebre) / denominator

			nextDraw := state.Discard
			nextDraw[Embiste] -= uint8(embiste)
			nextDraw[Aguante] -= uint8(aguante)
			nextDraw[Quiebre] -= uint8(quiebre)

			_ = removed
			states = append(states, weightedBeliefState{
				State: BeliefState{
					TrackPos:           state.TrackPos,
					Fish:               state.Fish,
					Draw:               nextDraw,
					Discard:            FamilyCounts{},
					FatigueCount:       state.FatigueCount,
					BaitGuardAvailable: state.BaitGuardAvailable,
					RodMods:            state.RodMods,
				},
				Probability: weight,
			})
		}
	}

	return states
}

func combination(n, k int) float64 {
	if k < 0 || k > n {
		return 0
	}
	if k == 0 || k == n {
		return 1
	}
	if k > n-k {
		k = n - k
	}
	result := 1.0
	for i := 1; i <= k; i++ {
		result *= float64(n-k+i) / float64(i)
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
