package svcmodel

import (
	"sync"

	"github.com/google/uuid"
	"github.com/llanuzo/card-game/pkg/httpapi"
)

type Game struct {
	Id uuid.UUID

	Players *Players
	Cards   *Cards
}

func NewGame(id uuid.UUID) *Game {
	return &Game{
		Id:      id,
		Players: NewPlayers(),
		Cards:   NewCards(),
	}
}

func (m *Game) ToHttp() *httpapi.Game {
	return &httpapi.Game{
		GameId: httpapi.NewUuid(m.Id),
	}
}

type Games struct {
	games sync.Map
}

func NewGames() *Games {
	return &Games{}
}

func (m *Games) Add(g *Game) {
	m.games.Store(g.Id, g)
}

func (m *Games) Load(id uuid.UUID) (*Game, bool) {
	val, ok := m.games.Load(id)
	if !ok {
		return nil, false
	}

	return val.(*Game), true
}

func (m *Games) Delete(id uuid.UUID) bool {
	if _, loaded := m.games.LoadAndDelete(id); loaded {
		return true
	}

	return false
}

func (m *Games) All() []*Game {
	var games []*Game
	m.games.Range(func(key, val any) bool {
		games = append(games, val.(*Game))
		return true
	})

	return games
}
