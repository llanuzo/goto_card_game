package main

import (
	"log/slog"

	"github.com/llanuzo/card-game/internal/app"
	"github.com/llanuzo/card-game/internal/log"
)

func main() {
	logger := log.NewLogger("startup")
	slog.SetDefault(logger)

	err := app.Start(logger)
	if err != nil {
		logger.Error("startup failed", slog.Any("error", err))
	}
}
