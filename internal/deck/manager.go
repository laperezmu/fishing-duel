package deck

import (
	"errors"

	"pesca/internal/domain"
)

var ErrNoCardsAvailable = errors.New("no cards available")

type Shuffler func([]domain.Move)

type RecyclePolicy interface {
	Recycle(discard []domain.Move, shuffler Shuffler) []domain.Move
}

type RemoveCardsRecyclePolicy struct {
	CardsToRemove int
}

func (p RemoveCardsRecyclePolicy) Recycle(discard []domain.Move, shuffler Shuffler) []domain.Move {
	refreshed := append([]domain.Move(nil), discard...)
	if shuffler != nil {
		shuffler(refreshed)
	}

	removeCount := p.CardsToRemove
	if removeCount < 0 {
		removeCount = 0
	}
	if removeCount > len(refreshed) {
		removeCount = len(refreshed)
	}

	return refreshed[:len(refreshed)-removeCount]
}

type Manager struct {
	active        []domain.Move
	discard       []domain.Move
	shuffler      Shuffler
	recyclePolicy RecyclePolicy
	recycleCount  int
	exhausted     bool
}

func NewManager(initial []domain.Move, shuffler Shuffler, recyclePolicy RecyclePolicy) *Manager {
	active := append([]domain.Move(nil), initial...)
	if shuffler != nil {
		shuffler(active)
	}
	if recyclePolicy == nil {
		recyclePolicy = RemoveCardsRecyclePolicy{CardsToRemove: 3}
	}

	return &Manager{
		active:        active,
		shuffler:      shuffler,
		recyclePolicy: recyclePolicy,
		exhausted:     len(active) == 0,
	}
}

func NewStandardFishDeck() []domain.Move {
	return []domain.Move{
		domain.Blue,
		domain.Blue,
		domain.Blue,
		domain.Red,
		domain.Red,
		domain.Red,
		domain.Yellow,
		domain.Yellow,
		domain.Yellow,
	}
}

func (m *Manager) Draw() (domain.Move, error) {
	if err := m.ensureActive(); err != nil {
		return 0, err
	}

	lastIndex := len(m.active) - 1
	move := m.active[lastIndex]
	m.active = m.active[:lastIndex]
	m.discard = append(m.discard, move)
	m.exhausted = false

	return move, nil
}

func (m *Manager) PrepareNextRound() {
	_ = m.ensureActive()
}

func (m *Manager) ActiveCount() int {
	return len(m.active)
}

func (m *Manager) DiscardCount() int {
	return len(m.discard)
}

func (m *Manager) RecycleCount() int {
	return m.recycleCount
}

func (m *Manager) Exhausted() bool {
	return m.exhausted
}

func (m *Manager) ensureActive() error {
	if len(m.active) > 0 {
		return nil
	}

	m.recycleIfNeeded()
	if len(m.active) == 0 {
		m.exhausted = true
		return ErrNoCardsAvailable
	}

	return nil
}

func (m *Manager) recycleIfNeeded() {
	if len(m.active) > 0 {
		return
	}
	if len(m.discard) == 0 {
		m.exhausted = true
		return
	}

	m.active = m.recyclePolicy.Recycle(m.discard, m.shuffler)
	m.discard = nil
	m.recycleCount++
	m.exhausted = len(m.active) == 0
}
