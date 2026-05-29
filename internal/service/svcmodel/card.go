package svcmodel

import "sync"

type Card struct {
	Value uint8
	Suit  CardSuit
}

type Cards struct {
	mu    sync.Mutex
	cards []Card
}

func NewCards() *Cards {
	return &Cards{}
}

func (m *Cards) Append(card []Card) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cards = append(m.cards, card...)
}

func (m *Cards) Next() Card {
	m.mu.Lock()
	defer m.mu.Unlock()

	card := m.cards[len(m.cards)-1]
	m.cards = m.cards[:len(m.cards)-1]

	return card
}

func (m *Cards) All() []Card {
	m.mu.Lock()
	defer m.mu.Unlock()

	copySlice := make([]Card, len(m.cards))
	copy(copySlice, m.cards)
	return copySlice
}
