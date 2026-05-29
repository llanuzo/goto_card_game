package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/llanuzo/card-game/internal/config"
	"github.com/llanuzo/card-game/internal/log"

	"github.com/llanuzo/card-game/internal/http"
	"github.com/llanuzo/card-game/internal/service"
)

type shutdowner interface {
	Shutdown(ctx context.Context) error
}

func Start(conf config.App, logger *log.Logger) error {
	logger.Infof("starting card-game")
	defer func() {
		logger.Infof("exiting card-game")
	}()

	svcs := service.NewServices()

	httpApi := http.NewApi(conf.HttpPort, svcs)
	go func() {
		logger.Infof("http api on port %d", conf.HttpPort)
		err := httpApi.Start()
		if err != nil {
			logger.Errorf("http server returned an err: %v", err)
			GracefulShutdown(logger)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shutdowners := []shutdowner{
		httpApi,
	}

	for _, val := range shutdowners {
		logger.Infof("shutting down %T...", val)
		if err := val.Shutdown(shutdownCtx); err != nil {
			logger.Errorf("error shutting down %T: %v", val, err)
			continue
		}
	}

	return nil
}

func GracefulShutdown(logger *log.Logger) {
	logger.Infof("grafully shutdown initiated")
	pid := os.Getpid()
	syscall.Kill(pid, syscall.SIGTERM)
}
