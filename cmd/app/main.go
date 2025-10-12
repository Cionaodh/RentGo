package main

import (
	"EasyRentGo/config"
	"EasyRentGo/internal/app"
	"log"
)

func main() {
	// init config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
