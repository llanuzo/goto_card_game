package httpapi

type Player struct {
	Id         UUID `json:"id"`
	CardsTotal int  `json:"cards_total"`
}
