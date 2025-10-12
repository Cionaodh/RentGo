package v1

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/usecase"
	"EasyRentGo/pkg/logger"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Controller - контроллер домена ProductTemplate API v1
type templateRoutes struct {
	tmp usecase.Template
	l   logger.Interface
	v   *validator.Validate
}

func newTemplateRoutes(tmp usecase.Template, l logger.Interface, v *validator.Validate) *templateRoutes {
	return &templateRoutes{tmp, l, v}
}

// Обработчики маршрутов

type TemplateDTO struct {
	Name        string `json:"name" validate:"required" example:"Велосипед 1"`
	Description string `json:"descriptiont" example:"Описание продукта"`
	Price       int    `json:"price" validate:"required" example:"800"`
}

func (r *templateRoutes) create(ctx *fiber.Ctx) error {
	var body TemplateDTO

	// Read body request
	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - create")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
		// return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	// Validation
	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "http - v1 - create")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	res, err := r.tmp.Create(
		ctx.UserContext(),
		entity.ProductTemp{
			Name:        body.Name,
			Description: body.Description,
			Price:       body.Price,
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - create")
		return errorResponse(ctx, http.StatusInternalServerError, "rentpoint service problems")
	}

	return ctx.Status(http.StatusCreated).JSON(res)
}

func (r *templateRoutes) getAll(ctx *fiber.Ctx) error {
	// points, err := r.rp.GetAll(ctx.UserContext())
	res, err := r.tmp.GetAll(ctx.UserContext())
	if err != nil {
		r.l.Error(err, "http - v1 - getAll")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to get rent points")
	}

	return ctx.Status(http.StatusOK).JSON(res)
}

func (r *templateRoutes) getByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - getByID - invalid UUID format")
		return errorResponse(ctx, http.StatusBadRequest, "invalid rentpoint ID format")
	}

	point, err := r.tmp.GetByID(ctx.UserContext(), id)
	if err != nil {
		// if err.Error() == "Product template by ID "+id.String()+" not found" {
		// 	return errorResponse(ctx, http.StatusNotFound, "rentpoint not found")
		// }
		r.l.Error(err, "http - v1 - getByID")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to get rent point")
	}

	return ctx.Status(http.StatusOK).JSON(point)
}
