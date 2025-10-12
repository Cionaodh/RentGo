package v1

// Содержит маршруты всех конечных точек

import (
	"EasyRentGo/internal/usecase"

	"EasyRentGo/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewRoutes(apiV1Group fiber.Router, uc *usecase.Usecases, l logger.Interface) {

	//
	rentPointGroup := apiV1Group.Group("/rentpoint")
	{
		r := newRentPointRoutes(uc.RentPoint, l, validator.New(validator.WithRequiredStructEnabled()))
		rentPointGroup.Post("/", r.create)    // POST /rentpoint/
		rentPointGroup.Get("/", r.getAll)     // GET /rentpoint/
		rentPointGroup.Get("/:id", r.getByID) // GET /rentpoint/{id}
	}

	//
	tmpProductGroup := apiV1Group.Group("/template")
	{
		t := newTemplateRoutes(uc.Template, l, validator.New(validator.WithRequiredStructEnabled()))
		tmpProductGroup.Post("/", t.create)    // POST /template/
		tmpProductGroup.Get("/", t.getAll)     // GET /template/
		tmpProductGroup.Get("/:id", t.getByID) // GET /template/{id}
	}

}
