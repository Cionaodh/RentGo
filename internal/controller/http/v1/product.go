package v1

import (
	"EasyRentGo/internal/usecase"
	"EasyRentGo/pkg/logger"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type productRoutes struct {
	pointUsecase usecase.Product
	l            logger.Interface
	v            *validator.Validate
}

func newProductRoutes(p usecase.Product, l logger.Interface, v *validator.Validate) *productRoutes {
	return &productRoutes{p, l, v}
}

type ProductDTO struct {
	TemplateId uuid.UUID `json:"template_id"` // id шаблона (обязательно)
	// id точки проката (необязательный)
	// статус - назначается на бизнес слое ()
	Number int `json:"number"` // количество (обязательно)
}

func (p *productRoutes) create(ctx *fiber.Ctx) error {
	// Считываем тело запроса в DTO
	var body ProductDTO
	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	// Валидируем
	if err := p.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	// Вызываем Usecase
	products, err := p.pointUsecase.Create(
		ctx.UserContext(),
		usecase.CreateProductInput{
			TemplateId: body.TemplateId,
			Number:     body.Number,
		})
	if err != nil {
		return errorResponse(ctx, http.StatusInternalServerError, "product service problems")
	}

	return ctx.Status(http.StatusCreated).JSON(products)
}

func (p *productRoutes) getAll(ctx *fiber.Ctx) error {
	products, err := p.pointUsecase.GetAll(ctx.UserContext())
	if err != nil {
		p.l.Error(err, "http - v1 - getAll")
		return errorResponse(ctx, http.StatusInternalServerError, "product service problems")
	}

	return ctx.Status(http.StatusOK).JSON(products)
}

func (p *productRoutes) getByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		p.l.Error(err, "http - v1 - getByID - invalid UUID format")
		return errorResponse(ctx, http.StatusBadRequest, "invalid rentpoint ID format")
	}

	product, err := p.pointUsecase.GetByID(ctx.UserContext(), id)
	if err != nil {
		p.l.Error(err, "http - v1 - create")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to get product")
	}

	return ctx.Status(http.StatusOK).JSON(product)
}
