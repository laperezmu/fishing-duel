package deck

import (
	"testing"

	"pesca/internal/domain"
)

func TestManagerRecyclesAndRemovesThreeCards(t *testing.T) {
	manager := NewManager(
		NewStandardFishDeck(),
		func([]domain.Move) {},
		RemoveCardsRecyclePolicy{CardsToRemove: 3},
	)

	for i := 0; i < 9; i++ {
		if _, err := manager.Draw(); err != nil {
			t.Fatalf("draw %d failed: %v", i+1, err)
		}
	}
	if manager.ActiveCount() != 0 || manager.DiscardCount() != 9 {
		t.Fatalf("after first pass active=%d discard=%d, want 0 and 9", manager.ActiveCount(), manager.DiscardCount())
	}

	manager.PrepareNextRound()
	if manager.ActiveCount() != 6 || manager.DiscardCount() != 0 || manager.RecycleCount() != 1 {
		t.Fatalf("after first recycle active=%d discard=%d recycles=%d, want 6 0 1", manager.ActiveCount(), manager.DiscardCount(), manager.RecycleCount())
	}

	for i := 0; i < 6; i++ {
		if _, err := manager.Draw(); err != nil {
			t.Fatalf("second cycle draw %d failed: %v", i+1, err)
		}
	}
	manager.PrepareNextRound()
	if manager.ActiveCount() != 3 || manager.RecycleCount() != 2 {
		t.Fatalf("after second recycle active=%d recycles=%d, want 3 and 2", manager.ActiveCount(), manager.RecycleCount())
	}

	for i := 0; i < 3; i++ {
		if _, err := manager.Draw(); err != nil {
			t.Fatalf("third cycle draw %d failed: %v", i+1, err)
		}
	}
	manager.PrepareNextRound()
	if manager.ActiveCount() != 0 || !manager.Exhausted() || manager.RecycleCount() != 3 {
		t.Fatalf("after final recycle active=%d exhausted=%t recycles=%d, want 0 true 3", manager.ActiveCount(), manager.Exhausted(), manager.RecycleCount())
	}
	if _, err := manager.Draw(); err != ErrNoCardsAvailable {
		t.Fatalf("draw on exhausted deck = %v, want %v", err, ErrNoCardsAvailable)
	}
}
