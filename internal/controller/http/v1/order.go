package v1

import (
	"EasyRentGo/internal/usecase"
	"errors"
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

type CreateOrderDTO struct {
	ProductID uuid.UUID `json:"product_id" validate:"required" example:""`
}

type CompleteOrderDTO struct {
	FinishingPointID uuid.UUID `json:"finishing_point_id" validate:"required" example:""`
}

func (r *orderRoutes) create(ctx *fiber.Ctx) error {
	// Достаем ID пользователя из контекста (положенный туда в middleware)
	userID, ok := ctx.Locals("userID").(uuid.UUID)
	if !ok {
		return errorResponse(ctx, http.StatusUnauthorized, "unauthorized")
	}

	var body CreateOrderDTO
	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, ErrInvalidRequestBody.Error())
	}

	if err := r.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, ErrInvalidParameters.Error())
	}

	// Передаем UserID в юзкейс
	order, err := r.orderUsecase.Create(ctx.UserContext(), usecase.CreateOrderInput{
		ProductID: body.ProductID,
		UserID:    userID,
	})

	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidProductID):
			return errorResponse(ctx, http.StatusBadRequest, err.Error())
		case errors.Is(err, usecase.ErrProductNotFound):
			return errorResponse(ctx, http.StatusNotFound, err.Error())
		case errors.Is(err, usecase.ErrProductAlreadyAssigned),
			errors.Is(err, usecase.ErrInvalidOrderState):
			return errorResponse(ctx, http.StatusConflict, err.Error())
		case errors.Is(err, usecase.ErrCreateOrder):
			return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
		}
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusCreated).JSON(order)
}

func (r *orderRoutes) getAll(ctx *fiber.Ctx) error {
	orders, err := r.orderUsecase.GetAll(ctx.UserContext())
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(orders)
}

func (r *orderRoutes) getByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid order id format")
	}

	order, err := r.orderUsecase.GetByID(ctx.UserContext(), id)
	if err != nil {

		switch {
		case errors.Is(err, usecase.ErrInvalidOrderID):
			return errorResponse(ctx, http.StatusBadRequest, err.Error())

		case errors.Is(err, usecase.ErrOrderNotFound):
			return errorResponse(ctx, http.StatusNotFound, err.Error())

		case errors.Is(err, usecase.ErrFetchOrders):
			return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
		}

		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(order)
}

func (r *orderRoutes) complete(ctx *fiber.Ctx) error {
	// Достаем ID текущего пользователя
	userID, ok := ctx.Locals("userID").(uuid.UUID)
	if !ok {
		return errorResponse(ctx, http.StatusUnauthorized, "unauthorized")
	}

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid order id format")
	}

	var body CompleteOrderDTO

	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, ErrInvalidRequestBody.Error())
	}

	if err := r.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, ErrInvalidParameters.Error())
	}

	order, err := r.orderUsecase.Complete(ctx.UserContext(), usecase.CompleteOrderInput{
		ID:               id,
		FinishingPointID: body.FinishingPointID,
		UserID:           userID, // Добавлено поле защиты
	})
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidOrderID),
			errors.Is(err, usecase.ErrInvalidRentpointID):
			return errorResponse(ctx, http.StatusBadRequest, err.Error())

		case errors.Is(err, usecase.ErrOrderNotFound):
			return errorResponse(ctx, http.StatusNotFound, err.Error())

		case errors.Is(err, usecase.ErrOrderNotActive),
			errors.Is(err, usecase.ErrInvalidOrderState),
			errors.Is(err, usecase.ErrProductAlreadyAssigned):
			return errorResponse(ctx, http.StatusConflict, err.Error())

		case errors.Is(err, usecase.ErrCompleteOrder):
			return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
		// Я бы еще добавил обработку ошибки доступа (Forbidden), если заказ чужой:
		case errors.Is(err, usecase.ErrForbidden):
			return errorResponse(ctx, http.StatusForbidden, err.Error())
		}

		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(order)
}
