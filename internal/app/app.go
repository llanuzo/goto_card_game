package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Start(logger *slog.Logger) error {
	logger.Info("starting card-game")
	defer func() {
		logger.Info("exiting card-game")
	}()

	// .. Initialize services

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Shutdown with timeout
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// .. Perform shutdown on initialized services

	return nil
}

func GracefulShutdown(logger *slog.Logger) {
	logger.Info("gracefully shutting down")
	pid := os.Getpid()
	syscall.Kill(pid, syscall.SIGTERM)
}
