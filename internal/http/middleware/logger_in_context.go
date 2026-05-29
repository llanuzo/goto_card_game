package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/llanuzo/card-game/internal/log"
)

const ctxKeyLogger ctxKey = "logger"

func NewLoggerInContext() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := log.NewLogger("http_request")
			logger = logger.With(slog.String("http_req_method", r.Method))
			logger = logger.With(slog.String("http_req_path", r.URL.Path))

			r = r.WithContext(context.WithValue(r.Context(), ctxKeyLogger, logger))

			next.ServeHTTP(w, r)
		})
	}
}

func GetLoggerFromContext(ctx context.Context) *log.Logger {
	logger, ok := ctx.Value(ctxKeyLogger).(*log.Logger)
	if !ok {
		logger = log.FallbackLogger
		logger.Errorf("failed to get middleware injected %s from context, defaulting to fallback logger", ctxKeyLogger)
	}

	return logger
}
