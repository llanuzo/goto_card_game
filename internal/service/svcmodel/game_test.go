//go:build unit

package svcmodel

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGame(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	g := NewGame(id)

	require.NotNil(t, g)
	assert.Equal(t, id, g.Id)
	assert.NotNil(t, g.Players)
	assert.NotNil(t, g.Cards)
}

func TestGame_ToHttp(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	g := NewGame(id)
	h := g.ToHttp()

	require.NotNil(t, h)
	assert.Equal(t, id, h.GameId.UUID)
}

func TestNewGames(t *testing.T) {
	t.Parallel()

	gs := NewGames()
	require.NotNil(t, gs)
	assert.Empty(t, gs.All())
}

func TestGames_Add_Load(t *testing.T) {
	t.Parallel()

	gs := NewGames()
	g := NewGame(uuid.New())
	gs.Add(g)

	got, found := gs.Load(g.Id)
	assert.True(t, found)
	assert.Equal(t, g, got)
}

func TestGames_Load_NotFound(t *testing.T) {
	t.Parallel()

	gs := NewGames()

	got, found := gs.Load(uuid.New())
	assert.False(t, found)
	assert.Nil(t, got)
}

func TestGames_Delete(t *testing.T) {
	t.Parallel()

	t.Run("existing game returns true", func(t *testing.T) {
		t.Parallel()

		gs := NewGames()
		g := NewGame(uuid.New())
		gs.Add(g)

		ok := gs.Delete(g.Id)
		assert.True(t, ok)

		_, found := gs.Load(g.Id)
		assert.False(t, found)
	})

	t.Run("missing game returns false", func(t *testing.T) {
		t.Parallel()

		gs := NewGames()
		ok := gs.Delete(uuid.New())
		assert.False(t, ok)
	})
}

func TestGames_All(t *testing.T) {
	t.Parallel()

	gs := NewGames()
	g1 := NewGame(uuid.New())
	g2 := NewGame(uuid.New())
	gs.Add(g1)
	gs.Add(g2)

	all := gs.All()
	assert.Len(t, all, 2)
	assert.ElementsMatch(t, []*Game{g1, g2}, all)
}
