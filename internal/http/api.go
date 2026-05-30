package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/llanuzo/card-game/internal/http/controller"
	"github.com/llanuzo/card-game/internal/http/middleware"
	"github.com/llanuzo/card-game/internal/service"
)

type Api struct {
	server *http.Server
}

func NewApi(port int, game service.Game) Api {
	r := mux.NewRouter()

	api := Api{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: r,
		},
	}

	mwLoggerInContext := middleware.NewLoggerInContext()

	r = r.PathPrefix("/api/v1").Subrouter().StrictSlash(false)
	r.Use(mwLoggerInContext)

	games := controller.NewGames(game)

	api.addRoute(r, http.MethodGet, "/games", games.List)
	api.addRoute(r, http.MethodPost, "/games", games.Post)
	api.addRoute(r, http.MethodDelete, "/games/{id1}", games.Delete)
	api.addRoute(r, http.MethodGet, "/games/{id1}/cards-by-suit", games.GetCardsBySuit)
	api.addRoute(r, http.MethodGet, "/games/{id1}/card-counts", games.ListCardCounts)

	return api
}

func (a Api) Start() error {
	err := a.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("unexpected error returned by http api server: %w", err)
	}

	return nil
}

func (a Api) Shutdown(ctx context.Context) error {
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shufown http server: %w", err)
	}

	return nil
}

func (a Api) addRoute(r *mux.Router, method string, path string, handler HandlerWithError, middlewares ...mux.MiddlewareFunc) {
	var httpHandlerFunc http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			HandleError(w, r, err)
			return
		}
	}

	var httpHandler http.Handler = httpHandlerFunc

	for i := len(middlewares) - 1; i >= 0; i-- {
		httpHandler = middlewares[i](httpHandler)
	}

	r.Handle(path, httpHandler).Methods(method)
}
