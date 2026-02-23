package v1

import (
	"EasyRentGo/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewRoutes(apiV1Group fiber.Router, uc *usecase.Usecases) {

	//
	rentPointGroup := apiV1Group.Group("/rentpoints")
	{
		r := newRentPointRoutes(uc.RentPoint, validator.New(validator.WithRequiredStructEnabled()))
		rentPointGroup.Post("/", r.create)         // POST /v1/rentpoints/
		rentPointGroup.Get("/", r.getAll)          // GET  /v1/rentpoints/
		rentPointGroup.Get("/:id", r.getByID)      // GET  /v1/rentpoints/{id}
		rentPointGroup.Post("/:id", r.addProducts) // POST  /v1/rentpoints/{id} -- TODO: Удалить конечную точку - данный метод должен вызываться в /v1/product
	}

	//
	templateGroup := apiV1Group.Group("/templates")
	{
		t := newTemplateRoutes(uc.Template, validator.New(validator.WithRequiredStructEnabled()))
		templateGroup.Post("/", t.create)    // POST /v1/templates/
		templateGroup.Get("/", t.getAll)     // GET  /v1/templates/
		templateGroup.Get("/:id", t.getByID) // GET  /v1/templates/{id}
	}

	//
	productGroup := apiV1Group.Group("/products")
	{
		p := newProductRoutes(uc.Product, validator.New(validator.WithRequiredStructEnabled()))
		productGroup.Post("/", p.create)     // POST /v1/products/
		productGroup.Get("/", p.getAll)      // GET /v1/products/
		productGroup.Get("/params/", p.list) // GET /v1/products/
		productGroup.Get("/:id", p.getByID)  // GET /v1/products/{id}
		// productGroup.Patch("/a")) // добавление продукта в точку проката
	}

	//
	orderGroup := apiV1Group.Group("/orders")
	{
		o := newOrderRoutes(uc.Order, validator.New(validator.WithRequiredStructEnabled()))
		orderGroup.Post("/", o.create)                // POST /v1/orders/
		orderGroup.Get("/", o.getAll)                 // GET /v1/orders/
		orderGroup.Get("/:id", o.getByID)             // GET /v1/orders/{id}
		orderGroup.Patch("/:id/complete", o.complete) // PATCH /v1/orders/{id}
	}
}
