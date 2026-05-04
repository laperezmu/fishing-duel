package run

import (
	"fmt"
	"pesca/internal/player/loadout"
)

func DefaultRoute() []NodeState {
	return []NodeState{
		{ZoneID: "shoreline-cove", NodeID: "start", Kind: NodeKindStart},
		{ZoneID: "shoreline-cove", NodeID: "fishing-1", Kind: NodeKindFishing, WaterPresetID: "shoreline-cove"},
		{ZoneID: "shoreline-cove", NodeID: "fishing-2", Kind: NodeKindFishing, WaterPresetID: "shoreline-cove"},
		{ZoneID: "shoreline-cove", NodeID: "boss-1", Kind: NodeKindBoss, WaterPresetID: "shoreline-cove"},
		{ZoneID: "open-channel", NodeID: "checkpoint-1", Kind: NodeKindCheckpoint},
		{ZoneID: "open-channel", NodeID: "fishing-3", Kind: NodeKindFishing, WaterPresetID: "open-channel"},
		{ZoneID: "open-channel", NodeID: "fishing-4", Kind: NodeKindFishing, WaterPresetID: "open-channel"},
		{ZoneID: "open-channel", NodeID: "boss-2", Kind: NodeKindBoss, WaterPresetID: "open-channel"},
		{ZoneID: "broken-current", NodeID: "service-1", Kind: NodeKindService},
		{ZoneID: "broken-current", NodeID: "fishing-5", Kind: NodeKindFishing, WaterPresetID: "broken-current"},
		{ZoneID: "broken-current", NodeID: "fishing-6", Kind: NodeKindFishing, WaterPresetID: "broken-current"},
		{ZoneID: "broken-current", NodeID: "boss-3", Kind: NodeKindBoss, WaterPresetID: "broken-current"},
		{ZoneID: "reef-shadow", NodeID: "checkpoint-2", Kind: NodeKindCheckpoint},
		{ZoneID: "reef-shadow", NodeID: "fishing-7", Kind: NodeKindFishing, WaterPresetID: "reef-shadow"},
		{ZoneID: "reef-shadow", NodeID: "fishing-8", Kind: NodeKindFishing, WaterPresetID: "reef-shadow"},
		{ZoneID: "reef-shadow", NodeID: "boss-4", Kind: NodeKindBoss, WaterPresetID: "reef-shadow"},
		{ZoneID: "tidal-gate", NodeID: "service-2", Kind: NodeKindService},
		{ZoneID: "tidal-gate", NodeID: "fishing-9", Kind: NodeKindFishing, WaterPresetID: "tidal-gate"},
		{ZoneID: "tidal-gate", NodeID: "fishing-10", Kind: NodeKindFishing, WaterPresetID: "tidal-gate"},
		{ZoneID: "tidal-gate", NodeID: "boss-5", Kind: NodeKindBoss, WaterPresetID: "tidal-gate"},
		{ZoneID: "weed-pocket", NodeID: "checkpoint-3", Kind: NodeKindCheckpoint},
		{ZoneID: "weed-pocket", NodeID: "fishing-11", Kind: NodeKindFishing, WaterPresetID: "weed-pocket"},
		{ZoneID: "weed-pocket", NodeID: "fishing-12", Kind: NodeKindFishing, WaterPresetID: "weed-pocket"},
		{ZoneID: "weed-pocket", NodeID: "boss-6", Kind: NodeKindBoss, WaterPresetID: "weed-pocket"},
		{ZoneID: "stone-drop", NodeID: "service-3", Kind: NodeKindService},
		{ZoneID: "stone-drop", NodeID: "fishing-13", Kind: NodeKindFishing, WaterPresetID: "stone-drop"},
		{ZoneID: "stone-drop", NodeID: "fishing-14", Kind: NodeKindFishing, WaterPresetID: "stone-drop"},
		{ZoneID: "stone-drop", NodeID: "boss-7", Kind: NodeKindBoss, WaterPresetID: "stone-drop"},
		{ZoneID: "deep-ledge", NodeID: "checkpoint-4", Kind: NodeKindCheckpoint},
		{ZoneID: "deep-ledge", NodeID: "fishing-15", Kind: NodeKindFishing, WaterPresetID: "deep-ledge"},
		{ZoneID: "deep-ledge", NodeID: "fishing-16", Kind: NodeKindFishing, WaterPresetID: "deep-ledge"},
		{ZoneID: "deep-ledge", NodeID: "boss-8", Kind: NodeKindBoss, WaterPresetID: "deep-ledge"},
		{ZoneID: "deep-ledge", NodeID: "end", Kind: NodeKindEnd},
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
