//go:build unit

package svcmodel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCard_ToHttp(t *testing.T) {
	t.Parallel()

	c := Card{Value: 10, Suit: CardSuit_Hearts}
	h := c.ToHttp()

	assert.Equal(t, "10_hearts", h.Id)
	assert.Equal(t, "hearts", h.Suit)
	assert.Equal(t, 10, h.FaceValue)
}

func TestCardCount_ToHttp(t *testing.T) {
	t.Parallel()

	cc := CardCount{Card: Card{Value: 5, Suit: CardSuit_Spades}, Count: 3}
	h := cc.ToHttp()

	assert.Equal(t, "5_spades", h.Card.Id)
	assert.Equal(t, 3, h.Count)
}

func TestNewCards(t *testing.T) {
	t.Parallel()

	cards := NewCards()
	require.NotNil(t, cards)
	assert.Empty(t, cards.All())
}

func TestCards_Append_All(t *testing.T) {
	t.Parallel()

	cards := NewCards()
	cards.Append([]Card{
		{Value: 1, Suit: CardSuit_Hearts},
		{Value: 2, Suit: CardSuit_Diamonds},
	})

	all := cards.All()
	assert.Len(t, all, 2)
}

func TestCards_All_ReturnsCopy(t *testing.T) {
	t.Parallel()

	cards := NewCards()
	cards.Append([]Card{{Value: 1, Suit: CardSuit_Hearts}})

	all := cards.All()
	all[0] = Card{Value: 99, Suit: CardSuit_Clubs}

	original := cards.All()
	assert.Equal(t, Card{Value: 1, Suit: CardSuit_Hearts}, original[0])
}

func TestCards_Next(t *testing.T) {
	t.Parallel()

	t.Run("next returns tail of slice", func(t *testing.T) {
		cards := NewCards()
		cards.Append([]Card{
			{Value: 1, Suit: CardSuit_Hearts},
			{Value: 2, Suit: CardSuit_Diamonds},
		})

		got, ok := cards.Next()
		assert.True(t, ok)
		assert.Equal(t, Card{Value: 2, Suit: CardSuit_Diamonds}, got)
		assert.Len(t, cards.All(), 1)
	})

	t.Run("next with no cards returns not ok", func(t *testing.T) {
		cards := NewCards()
		cards.Append([]Card{})

		_, ok := cards.Next()
		assert.False(t, ok)
	})

}

func TestCards_Shuffle(t *testing.T) {
	t.Parallel()

	t.Run("shuffle on empty deck", func(t *testing.T) {
		t.Parallel()
		cards := NewCards()
		cards.Shuffle()
	})

	t.Run("next with no cards returns not ok", func(t *testing.T) {
		t.Parallel()
		cards := NewCards()
		input := []Card{
			{Value: 1, Suit: CardSuit_Hearts},
			{Value: 2, Suit: CardSuit_Diamonds},
			{Value: 3, Suit: CardSuit_Spades},
			{Value: 4, Suit: CardSuit_Clubs},
			{Value: 5, Suit: CardSuit_Hearts},
		}
		cards.Append(input)
		cards.Shuffle()

		all := cards.All()
		assert.Len(t, all, len(input))
		assert.ElementsMatch(t, input, all)
	})
}
