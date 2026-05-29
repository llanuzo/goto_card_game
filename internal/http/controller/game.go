package controller

import (
	"errors"
	"net/http"

	"github.com/llanuzo/card-game/internal/service"
	"github.com/llanuzo/card-game/pkg/httpapi"
)

type Games struct {
	game service.Game
}

func NewGames(game service.Game) Games {
	return Games{
		game: game,
	}
}

func (c Games) List(w http.ResponseWriter, r *http.Request) error {
	games := c.game.List()

	listResp := httpapi.NewListResponse(len(games), func(i int) httpapi.ListResponseItem[*httpapi.Game] {
		return games[i]
	})

	return writeJson(w, http.StatusOK, &listResp)
}

func (c Games) Post(w http.ResponseWriter, r *http.Request) error {
	createdGame := c.game.Create()
	return writeJson(w, http.StatusOK, createdGame.ToHttp())
}

func (c Games) Delete(w http.ResponseWriter, r *http.Request) error {
	gameId, err := loadUuidFromPath(r, PathId1)
	if err != nil {
		return err
	}

	err = c.game.Delete(gameId)
	if err != nil {
		if errors.Is(err, service.ErrGameNotFound) {
			return newErrApiResponse(http.StatusNotFound, "game id %s does not exist", gameId)
		}

		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (c Games) GetCardsBySuit(w http.ResponseWriter, r *http.Request) error {
	gameId, err := loadUuidFromPath(r, PathId1)
	if err != nil {
		return err
	}

	resp, err := c.game.ListCardsBySuit(gameId)
	if err != nil {
		if errors.Is(err, service.ErrGameNotFound) {
			return newErrApiResponse(http.StatusNotFound, "game id %s does not exist", gameId)
		}

		return err
	}

	return writeJson(w, http.StatusOK, &resp)
}
