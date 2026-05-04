package run

import (
	"fmt"
	"pesca/internal/player/loadout"
)

func DefaultRoute() []NodeState {
	return []NodeState{
		{ZoneID: "coast", NodeID: "start", Kind: NodeKindStart},
		{ZoneID: "coast", NodeID: "fishing-1", Kind: NodeKindFishing, WaterPresetID: "shoreline-cove"},
		{ZoneID: "coast", NodeID: "fishing-2", Kind: NodeKindFishing, WaterPresetID: "broken-current"},
		{ZoneID: "coast", NodeID: "service-1", Kind: NodeKindService},
		{ZoneID: "outer-bank", NodeID: "boss-1", Kind: NodeKindBoss, WaterPresetID: "open-channel"},
		{ZoneID: "outer-bank", NodeID: "end", Kind: NodeKindEnd},
	}
}

func NewState(loadoutState loadout.State, route []NodeState, threadMaximum int) (State, error) {
	if len(route) == 0 {
		return State{}, fmt.Errorf("run route must include at least one node")
	}
	if threadMaximum <= 0 {
		return State{}, fmt.Errorf("run thread maximum must be greater than 0")
	}

	state := State{
		Status: StatusInProgress,
		Progress: ProgressState{
			ZoneIndex: 0,
			NodeIndex: 0,
			Current:   route[0],
			Next:      cloneNextNode(route, 1),
		},
		Thread:  ThreadState{Current: threadMaximum, Maximum: threadMaximum},
		Loadout: loadoutState,
	}

	if err := state.Validate(); err != nil {
		return State{}, err
	}

	return state, nil
}

func Advance(state *State, route []NodeState) error {
	if state == nil {
		return fmt.Errorf("run state is required")
	}
	if len(route) == 0 {
		return fmt.Errorf("run route must include at least one node")
	}
	if state.Progress.NodeIndex >= len(route)-1 {
		return fmt.Errorf("run route has no next node")
	}

	state.Progress.NodeIndex++
	state.Progress.Current = route[state.Progress.NodeIndex]
	state.Progress.Next = cloneNextNode(route, state.Progress.NodeIndex+1)
	state.Progress.ZoneIndex = resolveZoneIndex(route, state.Progress.NodeIndex)

	return state.Validate()
}

func ApplyEncounterResult(state *State, result EncounterResult) error {
	if state == nil {
		return fmt.Errorf("run state is required")
	}
	if err := result.Validate(); err != nil {
		return err
	}

	state.Thread.Current -= result.ThreadDamage
	if state.Thread.Current < 0 {
		state.Thread.Current = 0
	}
	if result.Capture != nil {
		state.Captures = append(state.Captures, *result.Capture)
	}

	switch result.Outcome {
	case EncounterOutcomeCaptured, EncounterOutcomeEscaped:
		if state.Thread.Current == 0 {
			state.Status = StatusDefeat
		}
	case EncounterOutcomeDefeated:
		state.Status = StatusDefeat
	case EncounterOutcomeRetired:
		state.Status = StatusRetired
	}

	return state.Validate()
}

func Complete(state *State) error {
	if state == nil {
		return fmt.Errorf("run state is required")
	}
	state.Status = StatusVictory

	return state.Validate()
}

func cloneNextNode(route []NodeState, index int) *NodeState {
	if index >= len(route) {
		return nil
	}
	next := route[index]

	return &next
}

func resolveZoneIndex(route []NodeState, currentIndex int) int {
	zoneIndex := 0
	for index := 1; index <= currentIndex; index++ {
		if route[index].ZoneID != route[index-1].ZoneID {
			zoneIndex++
		}
	}

	return zoneIndex
}
