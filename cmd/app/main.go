package main

import (
	"EasyRentGo/config"
	"EasyRentGo/internal/app"
	"EasyRentGo/pkg/logger"
	"log"
)

func main() {
	// Config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Logger
	l := logger.New(cfg.Log.Level)

	app.Migrations(l)
	app.Run(cfg, l)
}
