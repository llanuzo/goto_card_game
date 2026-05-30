//go:build unit

package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/llanuzo/card-game/internal/service"
	svcgenmock "github.com/llanuzo/card-game/internal/service/genmock"
	"github.com/llanuzo/card-game/internal/service/svcmodel"
	"github.com/llanuzo/card-game/pkg/httpapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRequestWithVars(method, url string, vars map[string]string) *http.Request {
	req := httptest.NewRequest(method, url, nil)
	return mux.SetURLVars(req, vars)
}

type gamesSuite struct {
	controller Games

	svcGames *svcgenmock.Game
}

func newGamesSuite(t *testing.T) gamesSuite {
	t.Helper()

	s := gamesSuite{
		svcGames: svcgenmock.NewGame(t),
	}

	s.controller = NewGames(s.svcGames)

	return s
}

func assertApiErr(t *testing.T, err error, wantStatus int) {
	t.Helper()
	var apiErr ErrApiResponse
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, wantStatus, apiErr.HttpStatusCode)
}

func TestGames_List(t *testing.T) {
	t.Parallel()

	t.Run("empty list returns 200 with empty items", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)

		s.svcGames.EXPECT().List().Return(nil)

		req := httptest.NewRequest(http.MethodGet, "/games", nil)
		w := httptest.NewRecorder()

		err := s.controller.List(w, req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var resp httpapi.ListResponse[*httpapi.Game]
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Empty(t, resp.Items)
	})

	t.Run("returns all games with 200", func(t *testing.T) {
		t.Parallel()

		id1 := uuid.New()
		id2 := uuid.New()
		games := []*svcmodel.Game{
			svcmodel.NewGame(id1),
			svcmodel.NewGame(id2),
		}

		s := newGamesSuite(t)
		s.svcGames.EXPECT().List().Return(games)

		req := httptest.NewRequest(http.MethodGet, "/games", nil)
		w := httptest.NewRecorder()

		err := s.controller.List(w, req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var resp httpapi.ListResponse[*httpapi.Game]
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		require.Len(t, resp.Items, 2)

		gotIds := []string{resp.Items[0].GameId.UUID.String(), resp.Items[1].GameId.UUID.String()}
		assert.ElementsMatch(t, []string{id1.String(), id2.String()}, gotIds)
	})
}

func TestGames_Post(t *testing.T) {
	t.Parallel()

	t.Run("creates game and returns 200", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		id := uuid.New()
		s.svcGames.EXPECT().Create().Return(svcmodel.NewGame(id))

		req := httptest.NewRequest(http.MethodPost, "/games", nil)
		w := httptest.NewRecorder()

		err := s.controller.Post(w, req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var resp httpapi.Game
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, id, resp.GameId.UUID)
	})
}

func TestGames_Delete(t *testing.T) {
	t.Parallel()

	t.Run("deletes game and returns 204", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		id := uuid.New()
		s.svcGames.EXPECT().Delete(id).Return(nil)

		req := newRequestWithVars(http.MethodDelete, "/games/"+id.String(), map[string]string{PathId1: id.String()})
		w := httptest.NewRecorder()

		err := s.controller.Delete(w, req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("game not found returns 404 error", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		id := uuid.New()
		s.svcGames.EXPECT().Delete(id).Return(service.ErrGameNotFound)

		req := newRequestWithVars(http.MethodDelete, "/games/"+id.String(), map[string]string{PathId1: id.String()})
		w := httptest.NewRecorder()

		assertApiErr(t, s.controller.Delete(w, req), http.StatusNotFound)
	})

	t.Run("missing game id returns 400 error", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		req := httptest.NewRequest(http.MethodDelete, "/games/", nil)
		w := httptest.NewRecorder()

		assertApiErr(t, s.controller.Delete(w, req), http.StatusBadRequest)
	})
}

func TestGames_GetCardsBySuit(t *testing.T) {
	t.Parallel()

	t.Run("returns card counts by suit with 200", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		id := uuid.New()
		s.svcGames.EXPECT().GetCardsBySuit(id).Return(map[svcmodel.CardSuit]int{
			svcmodel.CardSuit_Hearts:   3,
			svcmodel.CardSuit_Diamonds: 5,
			svcmodel.CardSuit_Clubs:    2,
			svcmodel.CardSuit_Spades:   7,
		}, nil)

		req := newRequestWithVars(http.MethodGet, "/games/"+id.String()+"/cards", map[string]string{PathId1: id.String()})
		w := httptest.NewRecorder()

		err := s.controller.GetCardsBySuit(w, req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var resp httpapi.GameCardsBySuite
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, 3, resp.Hearts)
		assert.Equal(t, 5, resp.Diamonds)
		assert.Equal(t, 2, resp.Clubs)
		assert.Equal(t, 7, resp.Spades)
	})

	t.Run("game not found returns 404 error", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		id := uuid.New()
		s.svcGames.EXPECT().GetCardsBySuit(id).Return(nil, service.ErrGameNotFound)

		req := newRequestWithVars(http.MethodGet, "/games/"+id.String()+"/cards", map[string]string{PathId1: id.String()})
		w := httptest.NewRecorder()

		assertApiErr(t, s.controller.GetCardsBySuit(w, req), http.StatusNotFound)
	})
}

func TestGames_ListCardCounts(t *testing.T) {
	t.Parallel()

	t.Run("returns card counts list with 200", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		id := uuid.New()
		cardCounts := []svcmodel.CardCount{
			{Card: svcmodel.Card{Value: 10, Suit: svcmodel.CardSuit_Hearts}, Count: 2},
			{Card: svcmodel.Card{Value: 5, Suit: svcmodel.CardSuit_Spades}, Count: 3},
		}
		s.svcGames.EXPECT().ListCardCounts(id).Return(cardCounts, nil)

		req := newRequestWithVars(http.MethodGet, "/games/"+id.String()+"/cards/count", map[string]string{PathId1: id.String()})
		w := httptest.NewRecorder()

		err := s.controller.ListCardCounts(w, req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var resp httpapi.ListResponse[httpapi.CardCount]
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Len(t, resp.Items, 2)
	})

	t.Run("game not found returns 404 error", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		id := uuid.New()
		s.svcGames.EXPECT().ListCardCounts(id).Return(nil, service.ErrGameNotFound)

		req := newRequestWithVars(http.MethodGet, "/games/"+id.String()+"/cards/count", map[string]string{PathId1: id.String()})
		w := httptest.NewRecorder()

		assertApiErr(t, s.controller.ListCardCounts(w, req), http.StatusNotFound)
	})
}

func TestGames_AddDeck(t *testing.T) {
	t.Parallel()

	t.Run("adds deck and returns 204", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		id := uuid.New()
		s.svcGames.EXPECT().AddDeck(id).Return(nil)

		req := newRequestWithVars(http.MethodPost, "/games/"+id.String()+"/deck", map[string]string{PathId1: id.String()})
		w := httptest.NewRecorder()

		err := s.controller.AddDeck(w, req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("game not found returns 404 error", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		id := uuid.New()
		s.svcGames.EXPECT().AddDeck(id).Return(service.ErrGameNotFound)

		req := newRequestWithVars(http.MethodPost, "/games/"+id.String()+"/deck", map[string]string{PathId1: id.String()})
		w := httptest.NewRecorder()

		assertApiErr(t, s.controller.AddDeck(w, req), http.StatusNotFound)
	})
}

func TestGames_Shuffle(t *testing.T) {
	t.Parallel()

	t.Run("shuffles and returns 204", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		id := uuid.New()
		s.svcGames.EXPECT().Shuffle(id).Return(nil)

		req := newRequestWithVars(http.MethodPost, "/games/"+id.String()+"/shuffle", map[string]string{PathId1: id.String()})
		w := httptest.NewRecorder()

		err := s.controller.Shuffle(w, req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("game not found returns 404 error", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		id := uuid.New()
		s.svcGames.EXPECT().Shuffle(id).Return(service.ErrGameNotFound)

		req := newRequestWithVars(http.MethodPost, "/games/"+id.String()+"/shuffle", map[string]string{PathId1: id.String()})
		w := httptest.NewRecorder()

		assertApiErr(t, s.controller.Shuffle(w, req), http.StatusNotFound)
	})
}

func TestGames_AddPlayer(t *testing.T) {
	t.Parallel()

	t.Run("adds player and returns 200", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		gameId := uuid.New()
		player := &svcmodel.Player{Id: uuid.New(), Cards: svcmodel.NewCards()}
		s.svcGames.EXPECT().AddPlayer(gameId).Return(player, nil)

		req := newRequestWithVars(http.MethodPost, "/games/"+gameId.String()+"/players", map[string]string{PathId1: gameId.String()})
		w := httptest.NewRecorder()

		err := s.controller.AddPlayer(w, req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var resp httpapi.Player
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, player.Id, resp.Id.UUID)
	})

	t.Run("game not found returns 404 error", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		gameId := uuid.New()
		s.svcGames.EXPECT().AddPlayer(gameId).Return(nil, service.ErrGameNotFound)

		req := newRequestWithVars(http.MethodPost, "/games/"+gameId.String()+"/players", map[string]string{PathId1: gameId.String()})
		w := httptest.NewRecorder()

		assertApiErr(t, s.controller.AddPlayer(w, req), http.StatusNotFound)
	})
}

func TestGames_DeletePlayer(t *testing.T) {
	t.Parallel()

	t.Run("deletes player and returns 204", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		gameId := uuid.New()
		playerId := uuid.New()
		s.svcGames.EXPECT().DeletePlayer(gameId, playerId).Return(nil)

		req := newRequestWithVars(http.MethodDelete, "/games/"+gameId.String()+"/players/"+playerId.String(), map[string]string{
			PathId1: gameId.String(),
			PathId2: playerId.String(),
		})
		w := httptest.NewRecorder()

		err := s.controller.DeletePlayer(w, req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("game not found returns 404 error", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		gameId := uuid.New()
		playerId := uuid.New()
		s.svcGames.EXPECT().DeletePlayer(gameId, playerId).Return(service.ErrGameNotFound)

		req := newRequestWithVars(http.MethodDelete, "/games/"+gameId.String()+"/players/"+playerId.String(), map[string]string{
			PathId1: gameId.String(),
			PathId2: playerId.String(),
		})
		w := httptest.NewRecorder()

		assertApiErr(t, s.controller.DeletePlayer(w, req), http.StatusNotFound)
	})

	t.Run("missing player id returns 400 error", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		gameId := uuid.New()

		req := newRequestWithVars(http.MethodDelete, "/games/"+gameId.String()+"/players/", map[string]string{PathId1: gameId.String()})
		w := httptest.NewRecorder()

		assertApiErr(t, s.controller.DeletePlayer(w, req), http.StatusBadRequest)
	})
}

func TestGames_ListPlayers(t *testing.T) {
	t.Parallel()

	t.Run("returns players list with 200", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		gameId := uuid.New()
		players := []*svcmodel.Player{
			{Id: uuid.New(), Cards: svcmodel.NewCards()},
			{Id: uuid.New(), Cards: svcmodel.NewCards()},
		}
		s.svcGames.EXPECT().ListPlayers(gameId).Return(players, nil)

		req := newRequestWithVars(http.MethodGet, "/games/"+gameId.String()+"/players", map[string]string{PathId1: gameId.String()})
		w := httptest.NewRecorder()

		err := s.controller.ListPlayers(w, req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var resp httpapi.ListResponse[*httpapi.Player]
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Len(t, resp.Items, 2)
	})

	t.Run("game not found returns 404 error", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		gameId := uuid.New()
		s.svcGames.EXPECT().ListPlayers(gameId).Return(nil, service.ErrGameNotFound)

		req := newRequestWithVars(http.MethodGet, "/games/"+gameId.String()+"/players", map[string]string{PathId1: gameId.String()})
		w := httptest.NewRecorder()

		assertApiErr(t, s.controller.ListPlayers(w, req), http.StatusNotFound)
	})
}

func TestGames_ListPlayersCards(t *testing.T) {
	t.Parallel()

	t.Run("returns player cards with 200", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		gameId := uuid.New()
		playerId := uuid.New()
		cards := []svcmodel.Card{
			{Value: 7, Suit: svcmodel.CardSuit_Hearts},
			{Value: 3, Suit: svcmodel.CardSuit_Clubs},
		}
		s.svcGames.EXPECT().ListPlayersCards(gameId, playerId).Return(cards, nil)

		req := newRequestWithVars(http.MethodGet, "/games/"+gameId.String()+"/players/"+playerId.String()+"/cards", map[string]string{
			PathId1: gameId.String(),
			PathId2: playerId.String(),
		})
		w := httptest.NewRecorder()

		err := s.controller.ListPlayersCards(w, req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var resp httpapi.ListResponse[httpapi.Card]
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Len(t, resp.Items, 2)
	})

	t.Run("game not found returns 404 error", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		gameId := uuid.New()
		playerId := uuid.New()
		s.svcGames.EXPECT().ListPlayersCards(gameId, playerId).Return(nil, service.ErrGameNotFound)

		req := newRequestWithVars(http.MethodGet, "/games/"+gameId.String()+"/players/"+playerId.String()+"/cards", map[string]string{
			PathId1: gameId.String(),
			PathId2: playerId.String(),
		})
		w := httptest.NewRecorder()

		assertApiErr(t, s.controller.ListPlayersCards(w, req), http.StatusNotFound)
	})

	t.Run("player not found returns 404 error", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		gameId := uuid.New()
		playerId := uuid.New()
		s.svcGames.EXPECT().ListPlayersCards(gameId, playerId).Return(nil, service.ErrPlayerNotFound)

		req := newRequestWithVars(http.MethodGet, "/games/"+gameId.String()+"/players/"+playerId.String()+"/cards", map[string]string{
			PathId1: gameId.String(),
			PathId2: playerId.String(),
		})
		w := httptest.NewRecorder()

		assertApiErr(t, s.controller.ListPlayersCards(w, req), http.StatusNotFound)
	})
}

func TestGames_PostPlayerCard(t *testing.T) {
	t.Parallel()

	t.Run("deals card to player and returns 204", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		gameId := uuid.New()
		playerId := uuid.New()
		s.svcGames.EXPECT().AddPlayerCard(gameId, playerId).Return(nil)

		req := newRequestWithVars(http.MethodPost, "/games/"+gameId.String()+"/players/"+playerId.String()+"/cards", map[string]string{
			PathId1: gameId.String(),
			PathId2: playerId.String(),
		})
		w := httptest.NewRecorder()

		err := s.controller.PostPlayerCard(w, req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("game not found returns 404 error", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		gameId := uuid.New()
		playerId := uuid.New()
		s.svcGames.EXPECT().AddPlayerCard(gameId, playerId).Return(service.ErrGameNotFound)

		req := newRequestWithVars(http.MethodPost, "/games/"+gameId.String()+"/players/"+playerId.String()+"/cards", map[string]string{
			PathId1: gameId.String(),
			PathId2: playerId.String(),
		})
		w := httptest.NewRecorder()

		assertApiErr(t, s.controller.PostPlayerCard(w, req), http.StatusNotFound)
	})

	t.Run("player not found returns 404 error", func(t *testing.T) {
		t.Parallel()

		s := newGamesSuite(t)
		gameId := uuid.New()
		playerId := uuid.New()
		s.svcGames.EXPECT().AddPlayerCard(gameId, playerId).Return(service.ErrPlayerNotFound)

		req := newRequestWithVars(http.MethodPost, "/games/"+gameId.String()+"/players/"+playerId.String()+"/cards", map[string]string{
			PathId1: gameId.String(),
			PathId2: playerId.String(),
		})
		w := httptest.NewRecorder()

		assertApiErr(t, s.controller.PostPlayerCard(w, req), http.StatusNotFound)
	})
}
