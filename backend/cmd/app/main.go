package main

import (
	"EasyRentGo/config"
	"EasyRentGo/internal/app"
	"log"
)

func main() {
	// Инициализация конфигуратора
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Запуск приложения
	app.Run(cfg)
}
