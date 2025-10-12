package v1

import (
	"EasyRentGo/internal/usecase"
	"EasyRentGo/pkg/logger"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Controller - контроллер домена RentPoint API v1
type rentPointRoutes struct {
	rp usecase.RentPoint
	l  logger.Interface
	v  *validator.Validate
}

func newRentPointRoutes(rp usecase.RentPoint, l logger.Interface, v *validator.Validate) *rentPointRoutes {
	return &rentPointRoutes{rp, l, v}
}

type RentpointDTO struct {
	Name string `json:"name"       validate:"required"  example:"RentPoint 1"`
	Addr string `json:"addr"       validate:"required"  example:"г. Калининград, ул Баласа"`
}

// Обработчики маршрутов

func (r *rentPointRoutes) create(ctx *fiber.Ctx) error {
	var body RentpointDTO

	// Read body request
	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - create")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	// Validation
	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "http - v1 - create")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	point, err := r.rp.CreateRentpoint(
		ctx.UserContext(),
		usecase.CreateRentpointInput{
			Name: body.Name,
			Addr: body.Addr,
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - create")
		return errorResponse(ctx, http.StatusInternalServerError, "rentpoint service problems")
	}

	return ctx.Status(http.StatusCreated).JSON(point)
}

func (r *rentPointRoutes) getAll(ctx *fiber.Ctx) error {
	points, err := r.rp.GetAll(ctx.UserContext())
	if err != nil {
		r.l.Error(err, "http - v1 - getAll")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to get rent points")
	}

	return ctx.Status(http.StatusOK).JSON(points)
}

func (r *rentPointRoutes) getByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - getByID - invalid UUID format")
		return errorResponse(ctx, http.StatusBadRequest, "invalid rentpoint ID format")
	}

	point, err := r.rp.GetByID(ctx.UserContext(), id)
	if err != nil {
		// if err.Error() == "точка проката с ID "+id.String()+" не найдена" {
		// 	return errorResponse(ctx, http.StatusNotFound, "rentpoint not found")
		// }
		r.l.Error(err, "http - v1 - getByID")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to get rent point")
	}

	return ctx.Status(http.StatusOK).JSON(point)
}
