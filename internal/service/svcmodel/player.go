package svcmodel

import (
	"maps"
	"sync"

	"github.com/google/uuid"
)

type Player struct {
	Id    uuid.UUID
	Cards *Cards
}

func NewPlayer(id uuid.UUID) *Player {
	return &Player{
		Id:    id,
		Cards: NewCards(),
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

func (m *Players) Add(key uuid.UUID, value *Player) {
	m.mu.Lock()
	defer m.mu.Unlock() // Ensures unlock even if a panic occurs
	m.players[key] = value
}

func (m *Players) Remove(key uuid.UUID) bool {
	m.mu.Lock()
	defer m.mu.Unlock() // Ensures unlock even if a panic occurs
	if _, ok := m.players[key]; ok {
		delete(m.players, key)
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
