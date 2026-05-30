package controller

import (
	"errors"
	"net/http"

	"github.com/llanuzo/card-game/internal/service"
	"github.com/llanuzo/card-game/internal/service/svcmodel"
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

	cardsMap, err := c.game.GetCardsBySuit(gameId)
	if err != nil {
		if errors.Is(err, service.ErrGameNotFound) {
			return newErrApiResponse(http.StatusNotFound, "game id %s does not exist", gameId)
		}

		return err
	}

	return writeJson(w, http.StatusOK, &httpapi.GameCardsBySuite{
		Hearts:   cardsMap[svcmodel.CardSuit_Hearts],
		Diamonds: cardsMap[svcmodel.CardSuit_Diamonds],
		Clubs:    cardsMap[svcmodel.CardSuit_Clubs],
		Spades:   cardsMap[svcmodel.CardSuit_Spades],
	})
}

func (c Games) ListCardCounts(w http.ResponseWriter, r *http.Request) error {
	gameId, err := loadUuidFromPath(r, PathId1)
	if err != nil {
		return err
	}

	cardCounts, err := c.game.ListCardCounts(gameId)
	if err != nil {
		if errors.Is(err, service.ErrGameNotFound) {
			return newErrApiResponse(http.StatusNotFound, "game id %s does not exist", gameId)
		}

		return err
	}

	listResp := httpapi.NewListResponse(len(cardCounts), func(i int) httpapi.ListResponseItem[httpapi.CardCount] {
		return cardCounts[i]
	})

	return writeJson(w, http.StatusOK, &listResp)
}

func (c Games) AddDeck(w http.ResponseWriter, r *http.Request) error {
	gameId, err := loadUuidFromPath(r, PathId1)
	if err != nil {
		return err
	}

	err = c.game.AddDeck(gameId)
	if err != nil {
		if errors.Is(err, service.ErrGameNotFound) {
			return newErrApiResponse(http.StatusNotFound, "game id %s does not exist", gameId)
		}

		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (c Games) Shuffle(w http.ResponseWriter, r *http.Request) error {
	gameId, err := loadUuidFromPath(r, PathId1)
	if err != nil {
		return err
	}

	err = c.game.Shuffle(gameId)
	if err != nil {
		if errors.Is(err, service.ErrGameNotFound) {
			return newErrApiResponse(http.StatusNotFound, "game id %s does not exist", gameId)
		}

		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (c Games) AddPlayer(w http.ResponseWriter, r *http.Request) error {
	gameId, err := loadUuidFromPath(r, PathId1)
	if err != nil {
		return err
	}

	player, err := c.game.AddPlayer(gameId)
	if err != nil {
		if errors.Is(err, service.ErrGameNotFound) {
			return newErrApiResponse(http.StatusNotFound, "game id %s does not exist", gameId)
		}

		return err
	}

	return writeJson(w, http.StatusOK, player.ToHttp())
}

func (c Games) DeletePlayer(w http.ResponseWriter, r *http.Request) error {
	gameId, err := loadUuidFromPath(r, PathId1)
	if err != nil {
		return err
	}

	playerId, err := loadUuidFromPath(r, PathId2)
	if err != nil {
		return err
	}

	err = c.game.DeletePlayer(gameId, playerId)
	if err != nil {
		if errors.Is(err, service.ErrGameNotFound) {
			return newErrApiResponse(http.StatusNotFound, "game id %s does not exist", gameId)
		}

		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (c Games) ListPlayers(w http.ResponseWriter, r *http.Request) error {
	gameId, err := loadUuidFromPath(r, PathId1)
	if err != nil {
		return err
	}

	players, err := c.game.ListPlayers(gameId)
	if err != nil {
		if errors.Is(err, service.ErrGameNotFound) {
			return newErrApiResponse(http.StatusNotFound, "game id %s does not exist", gameId)
		}

		return err
	}

	listResp := httpapi.NewListResponse(len(players), func(i int) httpapi.ListResponseItem[*httpapi.Player] {
		return players[i]
	})

	return writeJson(w, http.StatusOK, &listResp)
}
