package svcmodel

type CardSuit int

const (
	CardSuit_Hearts CardSuit = iota
	CardSuit_Diamonds
	CardSuit_Spades
	CardSuit_Clubs
)

func (e CardSuit) String() string {
	return map[CardSuit]string{
		CardSuit_Hearts:   "hearts",
		CardSuit_Diamonds: "diamonds",
		CardSuit_Spades:   "spades",
		CardSuit_Clubs:    "clubs",
	}[e]
}

func CardSuitValues() []CardSuit {
	return []CardSuit{
		CardSuit_Hearts,
		CardSuit_Diamonds,
		CardSuit_Clubs,
		CardSuit_Spades,
	}
}
