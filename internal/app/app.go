package app

import (
	"EasyRentGo/config"
	"EasyRentGo/internal/controller/http"
	"EasyRentGo/internal/repo"
	"EasyRentGo/internal/usecase"
	"EasyRentGo/pkg/httpserver"
	"EasyRentGo/pkg/logger"
	"EasyRentGo/pkg/postgres"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config, l *logger.Logger) {

	// Postgres
	l.Info("Initializing postgres...")
	pg, err := postgres.New(cfg.PG.Conn, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Services and repos
	l.Info("Initializing services and repos...")
	uc := usecase.NewUsecase(usecase.UsecaseDependencies{
		Repos: repo.NewPostgresRepo(pg),
	})

	// init server
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	http.NewRouter(httpServer.App, cfg, uc, l)
	httpServer.Start()

	// Graceful shutdown
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

}
