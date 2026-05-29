package main

import (
	"log/slog"

	"github.com/llanuzo/card-game/internal/app"
	"github.com/llanuzo/card-game/internal/config"
	"github.com/llanuzo/card-game/internal/log"
)

func main() {
	logger := log.NewLogger("startup")
	slog.SetDefault(logger.Logger)

	err := app.Start(config.NewApp(), logger)
	if err != nil {
		logger.Errorf("startup failed: %v", err)
	}
}
