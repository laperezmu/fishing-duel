package app

import (
	"pesca/internal/cards"
	"pesca/internal/deck"
	"pesca/internal/game"
)

type recycleCountOverrideDeck struct {
	base   game.FishDeck
	offset int
}

func applyRecycleCountOverride(base game.FishDeck, recycleCount int) *recycleCountOverrideDeck {
	return &recycleCountOverrideDeck{base: base, offset: recycleCount}
}

func (deckOverride *recycleCountOverrideDeck) Draw() (cards.FishCard, error) {
	return deckOverride.base.Draw()
}

func (deckOverride *recycleCountOverrideDeck) PrepareNextRound() {
	deckOverride.base.PrepareNextRound()
}

func (deckOverride *recycleCountOverrideDeck) ActiveCount() int {
	return deckOverride.base.ActiveCount()
}

func (deckOverride *recycleCountOverrideDeck) DiscardCount() int {
	return deckOverride.base.DiscardCount()
}

func (deckOverride *recycleCountOverrideDeck) RecycleCount() int {
	return deckOverride.base.RecycleCount() + deckOverride.offset
}

func (deckOverride *recycleCountOverrideDeck) Exhausted() bool {
	return deckOverride.base.Exhausted()
}

func (deckOverride *recycleCountOverrideDeck) VisibilitySnapshot() deck.VisibilitySnapshot {
	snapshot := deckOverride.base.VisibilitySnapshot()
	snapshot.RecycleCount += deckOverride.offset
	snapshot.CurrentCycle.Number += deckOverride.offset
	for index := range snapshot.PreviousCycles {
		snapshot.PreviousCycles[index].Number += deckOverride.offset
	}

	return snapshot
}
