package encounter

import "testing"

func TestDefaultConfigIsValid(t *testing.T) {
	state, err := NewState(DefaultConfig())
	if err != nil {
		t.Fatalf("NewState(DefaultConfig()) error = %v", err)
	}
	if state.Distance != 3 {
		t.Fatalf("distance = %d, want 3", state.Distance)
	}
	if state.Status != StatusOngoing {
		t.Fatalf("status = %q, want %q", state.Status, StatusOngoing)
	}
}

func TestConfigRejectsInitialDistanceAtOrBelowTwo(t *testing.T) {
	_, err := NewState(Config{
		InitialDistance:           2,
		CaptureDistance:           0,
		EscapeDistance:            5,
		ExhaustionCaptureDistance: 2,
		PlayerWinStep:             1,
		FishWinStep:               1,
	})
	if err == nil {
		t.Fatalf("expected validation error for initial distance <= 2")
	}
}
