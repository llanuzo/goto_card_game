package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/llanuzo/card-game/internal/http/controller"
	"github.com/llanuzo/card-game/internal/http/middleware"
	"github.com/llanuzo/card-game/pkg/httpapi"
)

type ErrWithOrigin interface {
	GetOrigin() string
	Error() string
}

type HandlerWithError func(w http.ResponseWriter, r *http.Request) error

type errorHandler func(err error) *httpapi.ErrorResponse

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	apiErr := &httpapi.ErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Message:    "internal server error",
	}

	handlers := []errorHandler{
		handleControllerErrors,
	}

	for _, handler := range handlers {
		if handledApiErr := handler(err); handledApiErr != nil {
			apiErr = handledApiErr
			break
		}
	}

	if apiErr.StatusCode == http.StatusInternalServerError {
		logger := middleware.GetLoggerFromContext(ctx)

		if errWithOrigin, ok := errors.AsType[ErrWithOrigin](err); ok {
			logger.Errorf("unexpected error: %v %v", err, errWithOrigin.GetOrigin())
		} else {
			logger.Errorf("unexpected error: %v", err)
		}
	} else if apiErr.Message == "" {
		apiErr.Message = err.Error()
	}

	w.WriteHeader(apiErr.StatusCode)
	json.NewEncoder(w).Encode(apiErr)
}

func handleControllerErrors(err error) *httpapi.ErrorResponse {
	var apiErr *httpapi.ErrorResponse

	if err, ok := errors.AsType[controller.ErrApiResponse](err); ok {
		apiErr = &httpapi.ErrorResponse{
			StatusCode: err.HttpStatusCode,
		}
	}

	return apiErr
}
