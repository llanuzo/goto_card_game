package svcmodel

import (
	"fmt"
	"math/rand/v2"
	"sync"

	"github.com/llanuzo/card-game/pkg/httpapi"
)

type Card struct {
	Value uint8
	Suit  CardSuit
}

func (m Card) ToHttp() httpapi.Card {
	return httpapi.Card{
		Id:        fmt.Sprintf("%d_%s", m.Value, m.Suit),
		Suit:      m.Suit.String(),
		FaceValue: int(m.Value),
	}
}

type CardCount struct {
	Card
	Count int
}

func (m CardCount) ToHttp() httpapi.CardCount {
	return httpapi.CardCount{
		Card:  m.Card.ToHttp(),
		Count: m.Count,
	}
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

func (m *Cards) Next() (Card, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.cards) == 0 {
		return Card{}, false
	}

	card := m.cards[len(m.cards)-1]
	m.cards = m.cards[:len(m.cards)-1]

	return card, true
}

func (m *Cards) All() []Card {
	m.mu.Lock()
	defer m.mu.Unlock()

	copySlice := make([]Card, len(m.cards))
	copy(copySlice, m.cards)
	return copySlice
}

func (m *Cards) Shuffle() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Shuffle by randomly selecting 2 indexes and swap
	// Do it by a factor of the cardLen just to make sure it shuffles enough
	cardLen := len(m.cards)
	for range cardLen * 10 {
		p1 := rand.IntN(cardLen)
		p2 := rand.IntN(cardLen)

		m.cards[p1], m.cards[p2] = m.cards[p2], m.cards[p1]
	}
}
