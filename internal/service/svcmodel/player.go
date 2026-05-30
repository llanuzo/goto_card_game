package svcmodel

import (
	"maps"
	"sync"

	"github.com/google/uuid"
	"github.com/llanuzo/card-game/pkg/httpapi"
)

type Player struct {
	Id    uuid.UUID
	Cards *Cards
}

func NewPlayer() *Player {
	return &Player{
		Id:    uuid.New(),
		Cards: NewCards(),
	}
}

func (m *Player) CardTotal() int {
	var cardsTotal int
	for _, val := range m.Cards.All() {
		cardsTotal += int(val.Value)
	}

	return cardsTotal
}

func (m *Player) ToHttp() *httpapi.Player {
	return &httpapi.Player{
		Id:         httpapi.NewUuid(m.Id),
		CardsTotal: m.CardTotal(),
	}
}

type Players struct {
	mu      sync.Mutex
	players map[uuid.UUID]*Player
}

func NewPlayers() *Players {
	return &Players{
		players: make(map[uuid.UUID]*Player),
	}
}

func (m *Players) Add(player *Player) {
	m.mu.Lock()
	defer m.mu.Unlock() // Ensures unlock even if a panic occurs
	m.players[player.Id] = player
}

func (m *Players) Delete(id uuid.UUID) bool {
	m.mu.Lock()
	defer m.mu.Unlock() // Ensures unlock even if a panic occurs
	if _, ok := m.players[id]; ok {
		delete(m.players, id)
		return true
	}

	return false
}

func (m *Players) All() map[uuid.UUID]*Player {
	m.mu.Lock()
	defer m.mu.Unlock()

	copy := make(map[uuid.UUID]*Player, len(m.players))
	maps.Copy(copy, m.players)

	return copy
}
