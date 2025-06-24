package v1

// Содержит маршруты всех конечных точек

import (
	"EasyRentGo/internal/usecase"

	"EasyRentGo/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewRentRoutes(apiV1Group fiber.Router, rp usecase.RentPoint, l logger.Interface) {
	r := &V1{rp: rp, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	// Group endpoints
	rentGroup := apiV1Group.Group("/rentpoint")

	{
		rentGroup.Post("/", r.create)    // POST /rentpoint/
		rentGroup.Get("/", r.getAll)     // GET /rentpoint/
		rentGroup.Get("/:id", r.getByID) // GET /rentpoint/{id}
	}
}
