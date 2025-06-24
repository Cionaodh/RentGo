// Package v1 implements routing paths. Each services in own file.
package http

// Инициализируем роутер, middleware

import (
	"EasyRentGo/internal/controller/http/middleware"
	v1 "EasyRentGo/internal/controller/http/v1"
	"EasyRentGo/internal/usecase"

	"EasyRentGo/config"
	"EasyRentGo/pkg/logger"

	// _ "EasyRentGo/docs" // Swagger docs.
	"github.com/gofiber/fiber/v2"
)

// NewRouter -.
// Swagger spec:
// @title       Титул
// @description Using a translation service as an example
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(app *fiber.App, cfg *config.Config, rp usecase.RentPoint, l logger.Interface) {
	// Options
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))

	// // Prometheus metrics
	// if cfg.Metrics.Enabled {
	// 	prometheus := fiberprometheus.New("my-service-name")
	// 	prometheus.RegisterAt(app, "/metrics")
	// 	app.Use(prometheus.Middleware)
	// }

	// // Swagger
	// if cfg.Swagger.Enabled {
	// 	app.Get("/swagger/*", swagger.HandlerDefault)
	// }

	// // K8s probe
	// app.Get("/healthz", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })

	// Routers
	apiV1Group := app.Group("/v1")
	{
		v1.NewRentRoutes(apiV1Group, rp, l)
	}
}
