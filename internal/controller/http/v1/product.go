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
	TemplateId string `json:"template_id" validate:"required,uuid4"` // id шаблона (обязательно)
	// id точки проката (необязательный)
	// статус - назначается на бизнес слое ()
	Number int `json:"number" validate:"required,gt=0,lte=1000"` // количество от 0 до 1000 (обязательно)
}

type ProductParamsDTO struct {
	Status      *entity.ProductStatus `json:"status" validate:"omitempty"` // поля могут быть пустыми
	TemplateID  *uuid.UUID            `json:"template_id" validate:"omitempty,uuid4"`
	RentPointID *uuid.UUID            `json:"rentpoint_id" validate:"omitempty,uuid4"`
}

func (p *productRoutes) create(ctx *fiber.Ctx) error {
	var body ProductDTO
	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, ErrInvalidRequestBody.Error())
	}

	if err := p.v.Struct(body); err != nil {
		// return errorResponse(ctx, http.StatusBadRequest, ErrInvalidParameters.Error())\
		return validationErrorResponse(ctx, err)
	}

	templateUUID, _ := uuid.Parse(body.TemplateId)

	result, err := p.productUsecase.Create(ctx.UserContext(), usecase.CreateProductInput{
		TemplateId: templateUUID,
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

func (p *productRoutes) list(ctx *fiber.Ctx) error {
	var params ProductParamsDTO

	if statusStr := ctx.Query("status"); statusStr != "" {
		st := entity.ProductStatus(statusStr)
		if !st.IsValid() {
			return errorResponse(ctx, http.StatusBadRequest, "invalid status")
		}
		params.Status = &st
	}

	if tmpIDStr := ctx.Query("template_id"); tmpIDStr != "" {
		id, err := uuid.Parse(tmpIDStr)
		if err != nil {
			return errorResponse(ctx, http.StatusBadRequest, "invalid template id format")
		}
		params.TemplateID = &id
	}

	if rpIDStr := ctx.Query("rentpoint_id"); rpIDStr != "" {
		id, err := uuid.Parse(rpIDStr)
		if err != nil {
			return errorResponse(ctx, http.StatusBadRequest, "invalid rentpoint id format")
		}
		params.RentPointID = &id
	}

	products, err := p.productUsecase.List(ctx.UserContext(), usecase.ProductParamsInput{
		Status:      params.Status,
		TemplateID:  params.TemplateID,
		RentPointID: params.RentPointID,
	})
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, "failed to fetch products")
	}

	return ctx.Status(http.StatusOK).JSON(products)
}

type AddToRentPointDTO struct {
	ProductIDs  uuid.UUIDs `json:"product_ids" validate:"required,min=1,dive,required"`
	RentPointID uuid.UUID  `json:"rentpoint_id" validate:"required"`
}

func (p *productRoutes) addToRentPoint(ctx *fiber.Ctx) error {
	var body AddToRentPointDTO
	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}
	if err := p.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "request body validation error")
	}

	products, err := p.productUsecase.MultipleAddToRentPoint(ctx.UserContext(), usecase.AddProductsInput{
		RentPointID: body.RentPointID,
		ProductIDs:  body.ProductIDs,
	})
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrRentPointNotFound):
			return errorResponse(ctx, http.StatusNotFound, err.Error())
		case errors.Is(err, usecase.ErrEmptyProductIDs):
			return errorResponse(ctx, http.StatusBadRequest, err.Error())
		case errors.Is(err, usecase.ErrProductsNotAvailable):
			return errorResponse(ctx, http.StatusConflict, err.Error())
		default:
			return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
		}
	}

	return ctx.Status(http.StatusOK).JSON(products)
}

// ======
// ======

func validationErrorResponse(ctx *fiber.Ctx, err error) error {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return errorResponse(ctx, http.StatusBadRequest, "invalid parameters")
	}

	errorsMap := make(map[string]string)
	for _, fe := range ve {
		errorsMap[fe.Field()] = validationMessage(fe)
	}

	return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
		"errors": errorsMap,
	})
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "uuid4":
		return "must be valid UUIDv4"
	case "gt":
		return "must be greater than " + fe.Param()
	case "lte":
		return "must be less than or equal to " + fe.Param()
	case "product_status":
		return "invalid product status"
	default:
		return "invalid value"
	}
}
