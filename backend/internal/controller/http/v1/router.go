package v1

// Содержит маршруты всех конечных точек

import (
	"EasyRentGo/internal/usecase"

	"EasyRentGo/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewRoutes(apiV1Group fiber.Router, rp usecase.RentPointUseCase, tmp usecase.TemplateUseCase, l logger.Interface) {
	r := &V1RentPointController{
		rp: rp,
		l:  l,
		v:  validator.New(validator.WithRequiredStructEnabled()),
	}

	// Group endpoints
	rentPointGroup := apiV1Group.Group("/rentpoint")

	{
		rentPointGroup.Post("/", r.create)    // POST /rentpoint/
		rentPointGroup.Get("/", r.getAll)     // GET /rentpoint/
		rentPointGroup.Get("/:id", r.getByID) // GET /rentpoint/{id}
	}

	tp := &V1ProductTmpController{
		tmp: tmp,
		l:   l,
		v:   validator.New(validator.WithRequiredStructEnabled()),
	}

	// Group endpoints
	tmpProductGroup := apiV1Group.Group("/tmpproduct")

	{
		tmpProductGroup.Post("/", tp.create)    // POST /tmpproduct/
		tmpProductGroup.Get("/", tp.getAll)     // GET /tmpproduct/
		tmpProductGroup.Get("/:id", tp.getByID) // GET /tmpproduct/{id}
	}

}
