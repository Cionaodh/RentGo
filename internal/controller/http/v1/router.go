package v1

// Содержит маршруты всех конечных точек

import (
	"EasyRentGo/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewRoutes(apiV1Group fiber.Router, uc *usecase.Usecases) {

	//
	rentPointGroup := apiV1Group.Group("/rentpoint")
	{
		r := newRentPointRoutes(uc.RentPoint, validator.New(validator.WithRequiredStructEnabled()))
		rentPointGroup.Post("/", r.create)         // POST /v1/rentpoint/
		rentPointGroup.Get("/", r.getAll)          // GET  /v1/rentpoint/
		rentPointGroup.Get("/:id", r.getByID)      // GET  /v1/rentpoint/{id}
		rentPointGroup.Post("/:id", r.addProducts) // POST  /v1/rentpoint/{id}
	}

	//
	templateGroup := apiV1Group.Group("/template")
	{
		t := newTemplateRoutes(uc.Template, validator.New(validator.WithRequiredStructEnabled()))
		templateGroup.Post("/", t.create)    // POST /v1/template/
		templateGroup.Get("/", t.getAll)     // GET  /v1/template/
		templateGroup.Get("/:id", t.getByID) // GET  /v1/template/{id}
	}

	//
	productGroup := apiV1Group.Group("/product")
	{
		p := newProductRoutes(uc.Product, validator.New(validator.WithRequiredStructEnabled()))
		productGroup.Post("/", p.create)    // POST /v1/product
		productGroup.Get("/", p.getAll)     // POST /v1/product
		productGroup.Get("/:id", p.getByID) // POST /v1/product/{id}
	}
}
