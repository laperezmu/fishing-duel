package deck

import (
	"errors"
	"pesca/internal/cards"
	"pesca/internal/domain"
)

var ErrNoCardsAvailable = errors.New("no cards available")

type Shuffler func([]cards.FishCard)

type RecyclePolicy interface {
	Recycle(discardedCards []cards.FishCard, shuffler Shuffler) []cards.FishCard
}

type RemoveCardsRecyclePolicy struct {
	CardsToRemove int
}

func (policy RemoveCardsRecyclePolicy) Recycle(discardedCards []cards.FishCard, shuffler Shuffler) []cards.FishCard {
	refreshedCards := append([]cards.FishCard(nil), discardedCards...)
	if shuffler != nil {
		shuffler(refreshedCards)
	}

	removeCount := policy.CardsToRemove
	if removeCount < 0 {
		removeCount = 0
	}
	if removeCount > len(refreshedCards) {
		removeCount = len(refreshedCards)
	}

	return refreshedCards[:len(refreshedCards)-removeCount]
}

type Deck struct {
	activeCards    []cards.FishCard
	discardedCards []cards.FishCard
	shuffler       Shuffler
	recyclePolicy  RecyclePolicy
	recycleCount   int
	exhausted      bool
}

func New(initialCards []cards.FishCard, shuffler Shuffler, recyclePolicy RecyclePolicy) *Deck {
	activeCards := append([]cards.FishCard(nil), initialCards...)
	if shuffler != nil {
		shuffler(activeCards)
	}
	if recyclePolicy == nil {
		recyclePolicy = RemoveCardsRecyclePolicy{CardsToRemove: 3}
	}

	return &Deck{
		activeCards:   activeCards,
		shuffler:      shuffler,
		recyclePolicy: recyclePolicy,
		exhausted:     len(activeCards) == 0,
	}
}

func NewStandardFishCards() []cards.FishCard {
	return []cards.FishCard{
		cards.NewFishCard(domain.Blue),
		cards.NewFishCard(domain.Blue),
		cards.NewFishCard(domain.Blue),
		cards.NewFishCard(domain.Red),
		cards.NewFishCard(domain.Red),
		cards.NewFishCard(domain.Red),
		cards.NewFishCard(domain.Yellow),
		cards.NewFishCard(domain.Yellow),
		cards.NewFishCard(domain.Yellow),
	}
}

func (deck *Deck) Draw() (cards.FishCard, error) {
	if err := deck.ensureCardsAvailable(); err != nil {
		return cards.FishCard{}, err
	}

	lastCardIndex := len(deck.activeCards) - 1
	drawnMove := deck.activeCards[lastCardIndex]
	deck.activeCards = deck.activeCards[:lastCardIndex]
	deck.discardedCards = append(deck.discardedCards, drawnMove)
	deck.exhausted = false

	return drawnMove, nil
}

func (deck *Deck) PrepareNextRound() {
	_ = deck.ensureCardsAvailable()
}

func (deck *Deck) ActiveCount() int {
	return len(deck.activeCards)
}

func (deck *Deck) DiscardCount() int {
	return len(deck.discardedCards)
}

func (deck *Deck) RecycleCount() int {
	return deck.recycleCount
}

func (deck *Deck) Exhausted() bool {
	return deck.exhausted
}

func (deck *Deck) ensureCardsAvailable() error {
	if len(deck.activeCards) > 0 {
		return nil
	}

	deck.recycleDiscardPileIfNeeded()
	if len(deck.activeCards) == 0 {
		deck.exhausted = true
		return ErrNoCardsAvailable
	}

	return nil
}

func (deck *Deck) recycleDiscardPileIfNeeded() {
	if len(deck.activeCards) > 0 {
		return
	}
	if len(deck.discardedCards) == 0 {
		deck.exhausted = true
		return
	}

	deck.activeCards = deck.recyclePolicy.Recycle(deck.discardedCards, deck.shuffler)
	deck.discardedCards = nil
	deck.recycleCount++
	deck.exhausted = len(deck.activeCards) == 0
}
