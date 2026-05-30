//go:build unit

package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/llanuzo/card-game/internal/service/svcmodel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newGameWithDeck(t *testing.T) (Game, uuid.UUID) {
	t.Helper()
	svc := NewGame()
	g := svc.Create()
	require.NoError(t, svc.AddDeck(g.Id))
	return svc, g.Id
}

func TestGame_List(t *testing.T) {
	t.Parallel()

	t.Run("empty returns empty slice", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		assert.Empty(t, svc.List())
	})

	t.Run("returns all created games sorted by id", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		g1 := svc.Create()
		g2 := svc.Create()

		games := svc.List()
		require.Len(t, games, 2)

		ids := []string{games[0].Id.String(), games[1].Id.String()}
		assert.Contains(t, ids, g1.Id.String())
		assert.Contains(t, ids, g2.Id.String())
		assert.True(t, ids[0] <= ids[1], "games should be sorted by id")
	})
}

func TestGame_Create(t *testing.T) {
	t.Parallel()

	svc := NewGame()
	g := svc.Create()

	require.NotNil(t, g)
	assert.NotEqual(t, uuid.UUID{}, g.Id)

	games := svc.List()
	require.Len(t, games, 1)
	assert.Equal(t, g.Id, games[0].Id)
}

func TestGame_Delete(t *testing.T) {
	t.Parallel()

	t.Run("deletes existing game", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		g := svc.Create()

		err := svc.Delete(g.Id)
		require.NoError(t, err)
		assert.Empty(t, svc.List())
	})

	t.Run("returns ErrGameNotFound for unknown game", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		err := svc.Delete(uuid.New())
		assert.ErrorIs(t, err, ErrGameNotFound)
	})
}

func TestGame_AddDeck(t *testing.T) {
	t.Parallel()

	t.Run("adds 52 cards to game deck", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		g := svc.Create()

		err := svc.AddDeck(g.Id)
		require.NoError(t, err)

		counts, err := svc.GetCardsBySuit(g.Id)
		require.NoError(t, err)
		total := counts[svcmodel.CardSuit_Hearts] + counts[svcmodel.CardSuit_Diamonds] +
			counts[svcmodel.CardSuit_Clubs] + counts[svcmodel.CardSuit_Spades]
		assert.Equal(t, 52, total)
	})

	t.Run("multiple decks accumulate cards", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		g := svc.Create()

		require.NoError(t, svc.AddDeck(g.Id))
		require.NoError(t, svc.AddDeck(g.Id))

		counts, err := svc.GetCardsBySuit(g.Id)
		require.NoError(t, err)
		total := counts[svcmodel.CardSuit_Hearts] + counts[svcmodel.CardSuit_Diamonds] +
			counts[svcmodel.CardSuit_Clubs] + counts[svcmodel.CardSuit_Spades]
		assert.Equal(t, 104, total)
	})

	t.Run("returns ErrGameNotFound for unknown game", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		err := svc.AddDeck(uuid.New())
		assert.ErrorIs(t, err, ErrGameNotFound)
	})
}

func TestGame_Shuffle(t *testing.T) {
	t.Parallel()

	t.Run("shuffles deck without error", func(t *testing.T) {
		t.Parallel()

		svc, gameId := newGameWithDeck(t)
		assert.NoError(t, svc.Shuffle(gameId))
	})

	t.Run("returns ErrGameNotFound for unknown game", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		err := svc.Shuffle(uuid.New())
		assert.ErrorIs(t, err, ErrGameNotFound)
	})
}

func TestGame_GetCardsBySuit(t *testing.T) {
	t.Parallel()

	t.Run("returns 13 cards per suit after one deck", func(t *testing.T) {
		t.Parallel()

		svc, gameId := newGameWithDeck(t)

		counts, err := svc.GetCardsBySuit(gameId)
		require.NoError(t, err)
		assert.Equal(t, 13, counts[svcmodel.CardSuit_Hearts])
		assert.Equal(t, 13, counts[svcmodel.CardSuit_Diamonds])
		assert.Equal(t, 13, counts[svcmodel.CardSuit_Clubs])
		assert.Equal(t, 13, counts[svcmodel.CardSuit_Spades])
	})

	t.Run("returns zeros for empty deck", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		g := svc.Create()

		counts, err := svc.GetCardsBySuit(g.Id)
		require.NoError(t, err)
		assert.Equal(t, 0, counts[svcmodel.CardSuit_Hearts])
		assert.Equal(t, 0, counts[svcmodel.CardSuit_Diamonds])
		assert.Equal(t, 0, counts[svcmodel.CardSuit_Clubs])
		assert.Equal(t, 0, counts[svcmodel.CardSuit_Spades])
	})

	t.Run("returns ErrGameNotFound for unknown game", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		_, err := svc.GetCardsBySuit(uuid.New())
		assert.ErrorIs(t, err, ErrGameNotFound)
	})
}

func TestGame_ListCardCounts(t *testing.T) {
	t.Parallel()

	t.Run("returns 52 entries grouped by suit and sorted by value descending within each suit", func(t *testing.T) {
		t.Parallel()

		svc, gameId := newGameWithDeck(t)

		cardCounts, err := svc.ListCardCounts(gameId)
		require.NoError(t, err)
		assert.Len(t, cardCounts, 52)

		// results are grouped by suit in order: Hearts, Spades, Clubs, Diamonds
		expectedSuitOrder := []svcmodel.CardSuit{
			svcmodel.CardSuit_Hearts,
			svcmodel.CardSuit_Spades,
			svcmodel.CardSuit_Clubs,
			svcmodel.CardSuit_Diamonds,
		}
		for suitIdx, suit := range expectedSuitOrder {
			group := cardCounts[suitIdx*13 : (suitIdx+1)*13]
			for _, cc := range group {
				assert.Equal(t, suit, cc.Suit)
			}
			for i := 1; i < len(group); i++ {
				assert.GreaterOrEqual(t, int(group[i-1].Value), int(group[i].Value),
					"cards within suit %v should be sorted by value descending", suit)
			}
		}
	})

	t.Run("two decks doubles each card count", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		g := svc.Create()
		require.NoError(t, svc.AddDeck(g.Id))
		require.NoError(t, svc.AddDeck(g.Id))

		cardCounts, err := svc.ListCardCounts(g.Id)
		require.NoError(t, err)
		assert.Len(t, cardCounts, 52)
		for _, cc := range cardCounts {
			assert.Equal(t, 2, cc.Count)
		}
	})

	t.Run("returns ErrGameNotFound for unknown game", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		_, err := svc.ListCardCounts(uuid.New())
		assert.ErrorIs(t, err, ErrGameNotFound)
	})
}

func TestGame_AddPlayer(t *testing.T) {
	t.Parallel()

	t.Run("adds player and returns player with id", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		g := svc.Create()

		player, err := svc.AddPlayer(g.Id)
		require.NoError(t, err)
		require.NotNil(t, player)
		assert.NotEqual(t, uuid.UUID{}, player.Id)
	})

	t.Run("returns ErrGameNotFound for unknown game", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		_, err := svc.AddPlayer(uuid.New())
		assert.ErrorIs(t, err, ErrGameNotFound)
	})
}

func TestGame_DeletePlayer(t *testing.T) {
	t.Parallel()

	t.Run("deletes existing player", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		g := svc.Create()
		player, err := svc.AddPlayer(g.Id)
		require.NoError(t, err)

		err = svc.DeletePlayer(g.Id, player.Id)
		require.NoError(t, err)

		players, err := svc.ListPlayers(g.Id)
		require.NoError(t, err)
		assert.Empty(t, players)
	})

	t.Run("returns ErrGameNotFound for unknown game", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		err := svc.DeletePlayer(uuid.New(), uuid.New())
		assert.ErrorIs(t, err, ErrGameNotFound)
	})

	t.Run("returns ErrPlayerNotFound for unknown player", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		g := svc.Create()

		err := svc.DeletePlayer(g.Id, uuid.New())
		assert.ErrorIs(t, err, ErrPlayerNotFound)
	})
}

func TestGame_ListPlayers(t *testing.T) {
	t.Parallel()

	t.Run("returns players sorted by card total descending", func(t *testing.T) {
		t.Parallel()

		svc, gameId := newGameWithDeck(t)

		p1, err := svc.AddPlayer(gameId)
		require.NoError(t, err)
		p2, err := svc.AddPlayer(gameId)
		require.NoError(t, err)

		// deal cards so p1 has more total value than p2
		for range 3 {
			require.NoError(t, svc.AddPlayerCard(gameId, p1.Id))
		}
		require.NoError(t, svc.AddPlayerCard(gameId, p2.Id))

		players, err := svc.ListPlayers(gameId)
		require.NoError(t, err)
		require.Len(t, players, 2)
		assert.GreaterOrEqual(t, players[0].CardTotal(), players[1].CardTotal())
	})

	t.Run("returns ErrGameNotFound for unknown game", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		_, err := svc.ListPlayers(uuid.New())
		assert.ErrorIs(t, err, ErrGameNotFound)
	})
}

func TestGame_ListPlayersCards(t *testing.T) {
	t.Parallel()

	t.Run("returns cards dealt to player", func(t *testing.T) {
		t.Parallel()

		svc, gameId := newGameWithDeck(t)
		player, err := svc.AddPlayer(gameId)
		require.NoError(t, err)

		require.NoError(t, svc.AddPlayerCard(gameId, player.Id))
		require.NoError(t, svc.AddPlayerCard(gameId, player.Id))

		cards, err := svc.ListPlayersCards(gameId, player.Id)
		require.NoError(t, err)
		assert.Len(t, cards, 2)
	})

	t.Run("returns ErrGameNotFound for unknown game", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		_, err := svc.ListPlayersCards(uuid.New(), uuid.New())
		assert.ErrorIs(t, err, ErrGameNotFound)
	})

	t.Run("returns ErrPlayerNotFound for unknown player", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		g := svc.Create()

		_, err := svc.ListPlayersCards(g.Id, uuid.New())
		assert.ErrorIs(t, err, ErrPlayerNotFound)
	})
}

func TestGame_AddPlayerCard(t *testing.T) {
	t.Parallel()

	t.Run("deals top card from deck to player", func(t *testing.T) {
		t.Parallel()

		svc, gameId := newGameWithDeck(t)
		player, err := svc.AddPlayer(gameId)
		require.NoError(t, err)

		err = svc.AddPlayerCard(gameId, player.Id)
		require.NoError(t, err)

		cards, err := svc.ListPlayersCards(gameId, player.Id)
		require.NoError(t, err)
		assert.Len(t, cards, 1)
	})

	t.Run("no-op when deck is empty", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		g := svc.Create()
		player, err := svc.AddPlayer(g.Id)
		require.NoError(t, err)

		err = svc.AddPlayerCard(g.Id, player.Id)
		require.NoError(t, err)

		cards, err := svc.ListPlayersCards(g.Id, player.Id)
		require.NoError(t, err)
		assert.Empty(t, cards)
	})

	t.Run("returns ErrGameNotFound for unknown game", func(t *testing.T) {
		t.Parallel()

		svc := NewGame()
		err := svc.AddPlayerCard(uuid.New(), uuid.New())
		assert.ErrorIs(t, err, ErrGameNotFound)
	})

	t.Run("returns ErrPlayerNotFound for unknown player", func(t *testing.T) {
		t.Parallel()

		svc, gameId := newGameWithDeck(t)
		err := svc.AddPlayerCard(gameId, uuid.New())
		assert.ErrorIs(t, err, ErrPlayerNotFound)
	})
}
