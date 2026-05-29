package controller

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	PathId1 = "id1"
)

func loadUuidFromPath(r *http.Request, name string) (uuid.UUID, error) {
	idStr := mux.Vars(r)[name]
	if idStr == "" {
		return uuid.Nil, newErrApiResponse(http.StatusBadRequest, "missing required %s in path", name)
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, newErrApiResponse(http.StatusBadRequest, "invalid uuid path param %s for %s", idStr, name)
	}

	return id, nil
}
