//go:build unit

package svcmodel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardSuit_String(t *testing.T) {
	t.Parallel()

	cases := []struct {
		suit     CardSuit
		expected string
	}{
		{CardSuit_Hearts, "hearts"},
		{CardSuit_Diamonds, "diamonds"},
		{CardSuit_Spades, "spades"},
		{CardSuit_Clubs, "clubs"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.expected, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, tc.suit.String())
		})
	}
}

func TestCardSuitValues(t *testing.T) {
	t.Parallel()

	suits := CardSuitValues()

	assert.Len(t, suits, 4)
	assert.ElementsMatch(t, []CardSuit{
		CardSuit_Hearts,
		CardSuit_Diamonds,
		CardSuit_Spades,
		CardSuit_Clubs,
	}, suits)
}
