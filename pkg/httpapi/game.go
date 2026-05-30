package httpapi

type Game struct {
	GameId UUID `json:"game_id"`
}

type GameCardsBySuite struct {
	Hearts   int `json:"hearts"`
	Diamonds int `json:"diamonds"`
	Clubs    int `json:"clubs"`
	Spades   int `json:"spades"`
}
