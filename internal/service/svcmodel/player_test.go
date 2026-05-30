//go:build unit

package svcmodel

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPlayer(t *testing.T) {
	t.Parallel()
	p := NewPlayer()

	require.NotNil(t, p)
	assert.NotEqual(t, [16]byte{}, p.Id)
	assert.NotNil(t, p.Cards)
}

func TestNewPlayer_UniqueIds(t *testing.T) {
	t.Parallel()
	p1 := NewPlayer()
	p2 := NewPlayer()

	assert.NotEqual(t, p1.Id, p2.Id)
}

func TestPlayer_CardTotal(t *testing.T) {
	t.Parallel()

	t.Run("no cards returns zero", func(t *testing.T) {
		t.Parallel()

		p := NewPlayer()
		assert.Equal(t, 0, p.CardTotal())
	})

	t.Run("sums all card values", func(t *testing.T) {
		t.Parallel()

		p := NewPlayer()
		p.Cards.Append([]Card{
			{Value: 5, Suit: CardSuit_Hearts},
			{Value: 10, Suit: CardSuit_Spades},
			{Value: 3, Suit: CardSuit_Clubs},
		})
		assert.Equal(t, 18, p.CardTotal())
	})

	t.Run("single card", func(t *testing.T) {
		t.Parallel()

		p := NewPlayer()
		p.Cards.Append([]Card{{Value: 7, Suit: CardSuit_Diamonds}})
		assert.Equal(t, 7, p.CardTotal())
	})
}

func TestPlayer_ToHttp(t *testing.T) {
	t.Parallel()

	p := NewPlayer()
	p.Cards.Append([]Card{
		{Value: 4, Suit: CardSuit_Hearts},
		{Value: 6, Suit: CardSuit_Clubs},
	})

	h := p.ToHttp()

	require.NotNil(t, h)
	assert.Equal(t, p.Id, h.Id.UUID)
	assert.Equal(t, 10, h.CardsTotal)
}

func TestNewPlayers(t *testing.T) {
	t.Parallel()

	ps := NewPlayers()

	require.NotNil(t, ps)
	assert.Empty(t, ps.All())
}

func TestPlayers_Add_Load(t *testing.T) {
	t.Parallel()

	ps := NewPlayers()
	p := NewPlayer()

	ps.Add(p)

	got, found := ps.Load(p.Id)
	assert.True(t, found)
	assert.Equal(t, p, got)
}

func TestPlayers_Load_NotFound(t *testing.T) {
	t.Parallel()

	ps := NewPlayers()
	p := NewPlayer()

	got, found := ps.Load(p.Id)
	assert.False(t, found)
	assert.Nil(t, got)
}

func TestPlayers_Delete(t *testing.T) {
	t.Parallel()
	t.Run("existing player returns true", func(t *testing.T) {
		ps := NewPlayers()
		p := NewPlayer()
		ps.Add(p)

		ok := ps.Delete(p.Id)
		assert.True(t, ok)

		_, found := ps.Load(p.Id)
		assert.False(t, found)
	})

	t.Run("missing player returns false", func(t *testing.T) {
		ps := NewPlayers()
		p := NewPlayer()

		ok := ps.Delete(p.Id)
		assert.False(t, ok)
	})
}

func TestPlayers_All(t *testing.T) {
	t.Parallel()

	ps := NewPlayers()
	p1 := NewPlayer()
	p2 := NewPlayer()
	ps.Add(p1)
	ps.Add(p2)

	all := ps.All()

	assert.Len(t, all, 2)
	assert.Contains(t, all, p1.Id)
	assert.Contains(t, all, p2.Id)
}

func TestPlayers_All_SnapshotIsolatedFromConcurrentMutations(t *testing.T) {
	t.Parallel()

	ps := NewPlayers()
	p := NewPlayer()
	p.Cards.Append([]Card{{Value: 5, Suit: CardSuit_Hearts}})
	ps.Add(p)

	snapshot := ps.All()
	totalBefore := snapshot[p.Id].CardTotal()

	// concurrently add cards to the live player
	var wg sync.WaitGroup
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.Cards.Append([]Card{{Value: 1, Suit: CardSuit_Clubs}})
		}()
	}
	wg.Wait()

	// snapshot must be unaffected
	assert.Equal(t, totalBefore, snapshot[p.Id].CardTotal())
}

func TestPlayers_All_ReturnsCopy(t *testing.T) {
	t.Parallel()

	ps := NewPlayers()
	p := NewPlayer()
	ps.Add(p)

	all := ps.All()
	delete(all, p.Id)

	_, found := ps.Load(p.Id)
	assert.True(t, found, "deleting from returned map should not affect the original")
}
