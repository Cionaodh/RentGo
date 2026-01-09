package v1

import (
	"EasyRentGo/internal/usecase"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Controller - контроллер домена RentPoint API v1
type orderRoutes struct {
	orderUsecase usecase.Order
	v            *validator.Validate
}

func newOrderRoutes(rp usecase.Order, v *validator.Validate) *orderRoutes {
	return &orderRoutes{rp, v}
}

type OrderDTO struct {
	ID uuid.UUID `json:"product_id" validate:"required" example:""`
}

func (r *orderRoutes) create(ctx *fiber.Ctx) error {
	var body OrderDTO

	// Read body request
	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	// Validation
	if err := r.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "request body validation error")
	}

	order, err := r.orderUsecase.Create(ctx.UserContext(), usecase.CreateOrderInput{
		ProductID: body.ID,
	})
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.Status(http.StatusCreated).JSON(order)
}

func (r *orderRoutes) getAll(ctx *fiber.Ctx) error {
	orders, err := r.orderUsecase.GetAll(ctx.UserContext())
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, "failed to get orders")
	}

	return ctx.Status(http.StatusOK).JSON(orders)
}

func (r *orderRoutes) getByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid order ID format")
	}

	point, err := r.orderUsecase.GetByID(ctx.UserContext(), id)
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, "failed to get order")
	}

	return ctx.Status(http.StatusOK).JSON(point)
}

type CompleteOrderDTO struct {
	FinishingPointID uuid.UUID `json:"finishing_point_id" validate:"required" example:""`
}

func (r *orderRoutes) complete(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid order ID format")
	}

	var body CompleteOrderDTO

	// Read body request
	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	// Validation
	if err := r.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "request body validation error")
	}

	order, err := r.orderUsecase.Complete(ctx.UserContext(), usecase.CompleteOrderInput{
		ID:               id,
		FinishingPointID: body.FinishingPointID,
	})
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.Status(http.StatusCreated).JSON(order)
}
