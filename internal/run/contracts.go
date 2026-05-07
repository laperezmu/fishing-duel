package run

import (
	"fmt"
	"pesca/internal/content/watercontexts"
	"pesca/internal/encounter"
	"pesca/internal/player/loadout"
)

type Status string

const (
	StatusInProgress Status = "in_progress"
	StatusVictory    Status = "victory"
	StatusDefeat     Status = "defeat"
	StatusRetired    Status = "retired"
)

func (status Status) Validate() error {
	switch status {
	case StatusInProgress, StatusVictory, StatusDefeat, StatusRetired:
		return nil
	default:
		return fmt.Errorf("invalid run status %q", status)
	}
}

type NodeKind string

const (
	NodeKindStart      NodeKind = "start"
	NodeKindFishing    NodeKind = "fishing"
	NodeKindService    NodeKind = "service"
	NodeKindCheckpoint NodeKind = "checkpoint"
	NodeKindBoss       NodeKind = "boss"
	NodeKindEnd        NodeKind = "end"
)

func (kind NodeKind) Validate() error {
	switch kind {
	case NodeKindStart, NodeKindFishing, NodeKindService, NodeKindCheckpoint, NodeKindBoss, NodeKindEnd:
		return nil
	default:
		return fmt.Errorf("invalid run node kind %q", kind)
	}
}

type NodeState struct {
	ZoneID        watercontexts.ID
	NodeID        string
	Kind          NodeKind
	WaterPresetID watercontexts.ID
}

func (node NodeState) Validate() error {
	if node.ZoneID == "" {
		return fmt.Errorf("run node zone id is required")
	}
	if node.NodeID == "" {
		return fmt.Errorf("run node id is required")
	}
	if node.WaterPresetID != "" {
		if _, err := ResolveWaterPreset(node.WaterPresetID); err != nil {
			return err
		}
	}
	return node.Kind.Validate()
}

func ResolveWaterPreset(id watercontexts.ID) (watercontexts.Preset, error) {
	return watercontexts.ResolveDefaultPreset(id)
}

type ProgressState struct {
	ZoneIndex int
	NodeIndex int
	Current   NodeState
	Next      *NodeState
}

func (progress ProgressState) Validate() error {
	if progress.ZoneIndex < 0 {
		return fmt.Errorf("run progress zone index must be greater than or equal to 0")
	}
	if progress.NodeIndex < 0 {
		return fmt.Errorf("run progress node index must be greater than or equal to 0")
	}
	if err := progress.Current.Validate(); err != nil {
		return fmt.Errorf("current node: %w", err)
	}
	if progress.Next != nil {
		next := *progress.Next
		if err := next.Validate(); err != nil {
			return fmt.Errorf("next node: %w", err)
		}
	}

	return nil
}

type ThreadState struct {
	Current int
	Maximum int
}

func (thread ThreadState) Validate() error {
	if thread.Maximum <= 0 {
		return fmt.Errorf("run thread maximum must be greater than 0")
	}
	if thread.Current < 0 {
		return fmt.Errorf("run thread current must be greater than or equal to 0")
	}
	if thread.Current > thread.Maximum {
		return fmt.Errorf("run thread current must be less than or equal to maximum")
	}

	return nil
}

type CaptureRecord struct {
	FishID      string
	FishName    string
	EncounterID string
}

func (capture CaptureRecord) Validate() error {
	if capture.FishID == "" {
		return fmt.Errorf("capture fish id is required")
	}
	if capture.FishName == "" {
		return fmt.Errorf("capture fish name is required")
	}

	return nil
}

type Modifier struct {
	ID          string
	Name        string
	Description string
}

type State struct {
	Status    Status
	Progress  ProgressState
	Thread    ThreadState
	Loadout   loadout.State
	Currency  int
	Captures  []CaptureRecord
	Modifiers []Modifier
}

func (state State) Validate() error {
	if err := state.Status.Validate(); err != nil {
		return err
	}
	if err := state.Progress.Validate(); err != nil {
		return fmt.Errorf("progress: %w", err)
	}
	if err := state.Thread.Validate(); err != nil {
		return fmt.Errorf("thread: %w", err)
	}
	if err := state.Loadout.Validate(); err != nil {
		return fmt.Errorf("loadout: %w", err)
	}
	if state.Currency < 0 {
		return fmt.Errorf("run currency must be greater than or equal to 0")
	}
	for _, capture := range state.Captures {
		if err := capture.Validate(); err != nil {
			return fmt.Errorf("capture: %w", err)
		}
	}

	return nil
}

type EncounterOutcome string

const (
	EncounterOutcomeCaptured EncounterOutcome = "captured"
	EncounterOutcomeEscaped  EncounterOutcome = "escaped"
	EncounterOutcomeDefeated EncounterOutcome = "defeated"
	EncounterOutcomeRetired  EncounterOutcome = "retired"
)

func (outcome EncounterOutcome) Validate() error {
	switch outcome {
	case EncounterOutcomeCaptured, EncounterOutcomeEscaped, EncounterOutcomeDefeated, EncounterOutcomeRetired:
		return nil
	default:
		return fmt.Errorf("invalid encounter outcome %q", outcome)
	}
}

type EncounterResult struct {
	Outcome       EncounterOutcome
	Status        encounter.Status
	EndReason     encounter.EndReason
	ThreadDamage  int
	Capture       *CaptureRecord
	NodeResolved  bool
	Retryable     bool
	FinishedMatch bool
}

func (result EncounterResult) Validate() error {
	if err := result.Outcome.Validate(); err != nil {
		return err
	}
	if result.ThreadDamage < 0 {
		return fmt.Errorf("encounter thread damage must be greater than or equal to 0")
	}
	if result.Status == encounter.StatusOngoing {
		return fmt.Errorf("encounter result cannot use ongoing status")
	}
	if result.Capture != nil {
		if err := result.Capture.Validate(); err != nil {
			return fmt.Errorf("capture: %w", err)
		}
	}
	if result.Outcome == EncounterOutcomeCaptured && result.Capture == nil {
		return fmt.Errorf("captured encounter result requires capture data")
	}

	return nil
}
