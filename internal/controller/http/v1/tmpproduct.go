package v1

import (
	"EasyRentGo/internal/usecase"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Controller - контроллер домена ProductTemplate API v1
type templateRoutes struct {
	temp usecase.Template
	v    *validator.Validate
}

func newTemplateRoutes(temp usecase.Template, v *validator.Validate) *templateRoutes {
	return &templateRoutes{temp, v}
}

// Обработчики маршрутов

type TemplateDTO struct {
	Name        string  `json:"name" validate:"required" example:"Велосипед 1"`
	Description string  `json:"description" example:"Описание продукта"`
	Price       float64 `json:"price" validate:"required" example:"800"`
}

func (r *templateRoutes) create(ctx *fiber.Ctx) error {
	var body TemplateDTO

	// Read body request
	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
		// return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	// Validation
	if err := r.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	temp, err := r.temp.Create(
		ctx.UserContext(),
		usecase.CreateTemplateInput{
			Name:        body.Name,
			Description: body.Description,
			Price:       body.Price,
		},
	)
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, "template service problems")
	}

	return ctx.Status(http.StatusCreated).JSON(temp)
}

func (r *templateRoutes) getAll(ctx *fiber.Ctx) error {
	templates, err := r.temp.GetAll(ctx.UserContext())
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, "failed to get all templates")
	}

	return ctx.Status(http.StatusOK).JSON(templates)
}

func (r *templateRoutes) getByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid rentpoint ID format")
	}

	point, err := r.temp.GetByID(ctx.UserContext(), id)
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, "failed to get rent point")
	}

	return ctx.Status(http.StatusOK).JSON(point)
}
