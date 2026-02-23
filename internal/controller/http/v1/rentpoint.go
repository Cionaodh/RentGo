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
type rentPointRoutes struct {
	pointUsecase usecase.RentPoint
	v            *validator.Validate
}

func newRentPointRoutes(rp usecase.RentPoint, v *validator.Validate) *rentPointRoutes {
	return &rentPointRoutes{rp, v}
}

type RentpointDTO struct {
	Name string `json:"name"       validate:"required"  example:"RentPoint 1"`
	Addr string `json:"addr"       validate:"required"  example:"г. Калининград, ул Баласа"`
}

func (r *rentPointRoutes) create(ctx *fiber.Ctx) error {
	var body RentpointDTO

	// Read body request
	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, ErrInvalidRequestBody.Error())
	}

	// Validation
	if err := r.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, ErrInvalidParameters.Error())
	}

	point, err := r.pointUsecase.CreateRentpoint(
		ctx.UserContext(),
		usecase.CreateRentpointInput{
			Name: body.Name,
			Addr: body.Addr,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrRentPointAlreadyExists):
			return errorResponse(ctx, http.StatusConflict, err.Error())

		case errors.Is(err, usecase.ErrFieldIsEmpty),
			errors.Is(err, usecase.ErrFieldIsTooLong):
			return errorResponse(ctx, http.StatusBadRequest, err.Error())

		default:
			return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
		}
	}

	return ctx.Status(http.StatusCreated).JSON(point)
}

func (r *rentPointRoutes) getAll(ctx *fiber.Ctx) error {
	points, err := r.pointUsecase.GetAll(ctx.UserContext())
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(points)
}

func (r *rentPointRoutes) getByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid rentpoint ID format")
	}

	point, err := r.pointUsecase.GetByID(ctx.UserContext(), id)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrRentPointNotFound):
			return errorResponse(ctx, http.StatusNotFound, err.Error())
		default:
			return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
		}
	}

	return ctx.Status(http.StatusOK).JSON(point)
}

type AddProductsDTO struct {
	IdProducts uuid.UUIDs `json:"products"`
}

func (r *rentPointRoutes) addProducts(ctx *fiber.Ctx) error {

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid rentpoint ID format")
	}

	var body AddProductsDTO
	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	if err := r.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "request body validation error")
	}

	rentpoint, err := r.pointUsecase.AddProduct(ctx.UserContext(), usecase.AddProductsInput{
		RentpointID: id,
		ProductsID:  body.IdProducts,
	})
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrRentPointNotFound),
			errors.Is(err, usecase.ErrProductNotFound):
			return errorResponse(ctx, http.StatusNotFound, err.Error())

		case errors.Is(err, usecase.ErrProductAlreadyAssigned):
			return errorResponse(ctx, http.StatusConflict, err.Error())

		default:
			return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
		}
	}

	return ctx.Status(http.StatusOK).JSON(rentpoint)
}
