package controller

import (
	"net/http"
)

type Games struct {
}

func NewGames() Games {
	return Games{}
}

func (c Games) PostGames(w http.ResponseWriter, r *http.Request) error {
	return newErrApiResponse(http.StatusTeapot, "teapot!")
}
