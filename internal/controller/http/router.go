// Package v1 implements routing paths. Each services in own file.
package http

import (
	"EasyRentGo/internal/controller/http/middleware"
	v1 "EasyRentGo/internal/controller/http/v1"
	"EasyRentGo/internal/usecase"

	"EasyRentGo/config"
	"EasyRentGo/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func NewRouter(app *fiber.App, cfg *config.Config, uc *usecase.Usecases, l logger.Interface) {
	// Options
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))

	// Routers
	apiV1Group := app.Group("/v1")
	{
		v1.NewRoutes(apiV1Group, uc)
	}
}
