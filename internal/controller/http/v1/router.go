package v1

import (
	"EasyRentGo/internal/controller/http/middleware"
	"EasyRentGo/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewRoutes(apiV1Group fiber.Router, uc *usecase.Usecases) {

	//
	rentPointGroup := apiV1Group.Group("/rentpoints")
	{
		r := newRentPointRoutes(uc.RentPoint, validator.New(validator.WithRequiredStructEnabled()))
		rentPointGroup.Post("/", r.create)    // POST /v1/rentpoints/
		rentPointGroup.Get("/", r.getAll)     // GET  /v1/rentpoints/
		rentPointGroup.Get("/:id", r.getByID) // GET  /v1/rentpoints/{id}
		rentPointGroup.Delete("/:id", r.delete)
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
		productGroup.Post("/", p.create)          // POST /v1/products/
		productGroup.Get("/", p.getAll)           // GET /v1/products/
		productGroup.Get("/params/", p.list)      // GET /v1/products/
		productGroup.Get("/:id", p.getByID)       // GET /v1/products/{id}
		productGroup.Patch("/", p.addToRentPoint) // PATCH /v1/products/
	}

	//
	orderGroup := apiV1Group.Group("/orders", middleware.Auth())
	{
		o := newOrderRoutes(uc.Order, validator.New(validator.WithRequiredStructEnabled()))
		orderGroup.Post("/", o.create)                // POST /v1/orders/
		orderGroup.Get("/", o.getAll)                 // GET /v1/orders/
		orderGroup.Get("/:id", o.getByID)             // GET /v1/orders/{id}
		orderGroup.Patch("/:id/complete", o.complete) // PATCH /v1/orders/{id}
	}

	// Пользователи
	userGroup := apiV1Group.Group("/users")
	{
		u := newUserRoutes(uc.User, validator.New(validator.WithRequiredStructEnabled()))
		userGroup.Post("/register", u.register)
		userGroup.Post("/login", u.login)
		userGroup.Get("/:id", middleware.Auth(), u.getProfile)
		// userGroup.Get("/:id", u.getProfile)
	}
}
