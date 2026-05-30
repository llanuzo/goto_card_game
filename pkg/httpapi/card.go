package httpapi

type Card struct {
	Id        string `json:"id"`
	Suit      string `json:"suit"`
	FaceValue int    `json:"face_value"`
}

type CardCount struct {
	Card  Card `json:"card"`
	Count int  `json:"count"`
}
