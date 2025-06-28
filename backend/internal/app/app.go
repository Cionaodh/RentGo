package app

import (
	"EasyRentGo/config"
	"EasyRentGo/internal/controller/http"
	"EasyRentGo/internal/repo/nodb"
	"EasyRentGo/internal/usecase/producttemp"
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
	rentPoint := rentpoint.New(nodb.NewRentPoint())
	tmpProduct := producttemp.New(nodb.NewTemplate())

	// Инициализация серевера
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	http.NewRouter(
		httpServer.App,
		cfg,
		rentPoint,  // домен точки проката
		tmpProduct, // домен шаблона продукта
		l)

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
