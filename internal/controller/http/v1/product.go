package v1

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/usecase"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type productRoutes struct {
	productUsecase usecase.Product
	v              *validator.Validate
}

func newProductRoutes(p usecase.Product, v *validator.Validate) *productRoutes {
	return &productRoutes{p, v}
}

type ProductDTO struct {
	TemplateId uuid.UUID `json:"template_id"` // id шаблона (обязательно)
	// id точки проката (необязательный)
	// статус - назначается на бизнес слое ()
	Number int `json:"number"` // количество (обязательно)
}

func (p *productRoutes) create(ctx *fiber.Ctx) error {
	var body ProductDTO
	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, ErrInvalidRequestBody.Error())
	}

	if err := p.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, ErrInvalidParameters.Error())
	}

	result, err := p.productUsecase.Create(ctx.UserContext(), usecase.CreateProductInput{
		TemplateId: body.TemplateId,
		Number:     body.Number,
	})
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidTemplateID),
			errors.Is(err, usecase.ErrInvalidProductQty):
			return errorResponse(ctx, http.StatusBadRequest, err.Error())

		case errors.Is(err, usecase.ErrTemplateNotFound):
			return errorResponse(ctx, http.StatusNotFound, err.Error())

		default:
			return errorResponse(ctx, http.StatusInternalServerError, "failed to create product")
		}
	}

	return ctx.Status(http.StatusCreated).JSON(result)
}

func (p *productRoutes) getAll(ctx *fiber.Ctx) error {
	products, err := p.productUsecase.GetAll(ctx.UserContext())
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, "failed to fetch products")
	}

	return ctx.Status(http.StatusOK).JSON(products)
}

func (p *productRoutes) getByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid product id format")
	}

	product, err := p.productUsecase.GetByID(ctx.UserContext(), id)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidProductID):
			return errorResponse(ctx, http.StatusBadRequest, err.Error())

		case errors.Is(err, usecase.ErrProductNotFound):
			return errorResponse(ctx, http.StatusNotFound, err.Error())

		default:
			return errorResponse(ctx, http.StatusInternalServerError, "failed to fetch product")
		}
	}

	return ctx.Status(http.StatusOK).JSON(product)
}

type ProductParamsDTO struct {
	Status      *entity.ProductStatus `json:"status" validate:"omitempty"`
	TemplateID  *uuid.UUID            `json:"template_id" validate:"omitempty"`
	RentPointID *uuid.UUID            `json:"rentpoint_id" validate:"omitempty"`
}

func (p *productRoutes) list(ctx *fiber.Ctx) error {
	var input ProductParamsDTO

	if statusStr := ctx.Query("status"); statusStr != "" {
		st := entity.ProductStatus(statusStr)
		if !st.IsValid() {
			return errorResponse(ctx, http.StatusBadRequest, "invalid status")
		}
		input.Status = &st
	}

	if tmpIDStr := ctx.Query("template_id"); tmpIDStr != "" {
		id, err := uuid.Parse(tmpIDStr)
		if err != nil {
			return errorResponse(ctx, http.StatusBadRequest, "invalid template id format")
		}
		input.TemplateID = &id
	}

	if rpIDStr := ctx.Query("rentpoint_id"); rpIDStr != "" {
		id, err := uuid.Parse(rpIDStr)
		if err != nil {
			return errorResponse(ctx, http.StatusBadRequest, "invalid rentpoint id format")
		}
		input.RentPointID = &id
	}

	products, err := p.productUsecase.List(ctx.UserContext(), usecase.ProductParamsInput{
		Status:      input.Status,
		TemplateID:  input.TemplateID,
		RentPointID: input.RentPointID,
	})
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, "failed to fetch products")
	}

	return ctx.Status(http.StatusOK).JSON(products)
}
