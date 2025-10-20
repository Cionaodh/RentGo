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
	pointUsecase usecase.RentPoint
	l            logger.Interface
	v            *validator.Validate
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

	point, err := r.pointUsecase.CreateRentpoint(
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

	type response struct {
		Id   uuid.UUID `json:"id"`
		Name string    `json:"name"`
		Addr string    `json:"addr"`
	}

	p := response{
		Id:   point.ID,
		Name: point.Name,
		Addr: point.ID.String(),
	}

	// return ctx.Status(http.StatusCreated).JSON(point)
	return ctx.Status(http.StatusCreated).JSON(p)
}

func (r *rentPointRoutes) getAll(ctx *fiber.Ctx) error {
	points, err := r.pointUsecase.GetAll(ctx.UserContext())
	if err != nil {
		r.l.Error(err, "http - v1 - getAll")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to get rent points")
	}

	type response struct {
		Id   uuid.UUID `json:"id"`
		Name string    `json:"name"`
		Addr string    `json:"addr"`
	}

	ps := make([]response, 0, len(points))
	for _, point := range points {
		ps = append(ps, response{
			Id:   point.ID,
			Name: point.Name,
			Addr: point.Addr,
		})
	}

	return ctx.Status(http.StatusOK).JSON(ps)
}

// Get - с id --> структура с данными о точке проката и всех её продуктах
func (r *rentPointRoutes) getByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - getByID - invalid UUID format")
		return errorResponse(ctx, http.StatusBadRequest, "invalid rentpoint ID format")
	}

	point, err := r.pointUsecase.GetByID(ctx.UserContext(), id)
	if err != nil {
		// if err.Error() == "точка проката с ID "+id.String()+" не найдена" {
		// 	return errorResponse(ctx, http.StatusNotFound, "rentpoint not found")
		// }
		r.l.Error(err, "http - v1 - getByID")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to get rent point")
	}

	type response struct {
		Id   uuid.UUID `json:"id"`
		Name string    `json:"name"`
		Addr string    `json:"addr"`
		// products []struct {
		// 	ID         uuid.UUID `json:"id"`
		// 	TemplateID uuid.UUID `json:"template"` // TODO: подтянуть данные из шаблона
		// }
	}

	p := response{
		Id:   point.ID,
		Name: point.Name,
		Addr: point.Addr,
	}

	return ctx.Status(http.StatusOK).JSON(p)
}
