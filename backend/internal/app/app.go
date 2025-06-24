package app

import (
	"EasyRentGo/config"
	"EasyRentGo/internal/controller/http"
	"EasyRentGo/internal/repo/persistent"
	"EasyRentGo/internal/usecase/rentpoint"
	"EasyRentGo/pkg/httpserver"
	"EasyRentGo/pkg/logger"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) { // Передача конфигуратора

	// Инициализация логгера
	l := logger.New(cfg.Log.Level)

	// Инициализация БД

	// Инициализация слоя бизнес-логики
	rentpoint := rentpoint.New(persistent.New())

	// Инициализация серевера
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	http.NewRouter(httpServer.App, cfg, rentpoint, l)

	// Start servers
	httpServer.Start()

	// Graceful shutdown (Ожиданеие сигнала, завершение)
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
	err := httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

}
