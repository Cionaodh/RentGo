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

	type response struct {
		TemplateID  uuid.UUID  `json:"template_id"`
		RentPointID uuid.UUID  `json:"rentpoint_id"`
		Status      string     `json:"status"`
		Num         int        `json:"number"`
		Ids         uuid.UUIDs `json:"ids"`
	}

	prod := response{
		TemplateID:  products.TemplateID,
		RentPointID: products.RentPointID,
		Status:      string(products.Status),
		Ids:         products.IDs,
	}

	return ctx.Status(http.StatusCreated).JSON(prod)
}

func (p *productRoutes) getAll(ctx *fiber.Ctx) error {
	products, err := p.pointUsecase.GetAll(ctx.UserContext())
	if err != nil {
		p.l.Error(err, "http - v1 - getAll")
		return errorResponse(ctx, http.StatusInternalServerError, "product service problems")
	}

	type response struct {
		Name        string     `json:"name"`
		Price       float64    `json:"price"`
		RentPointID uuid.UUID  `json:"rentpoint_id"`
		Status      string     `json:"status"`
		Num         int        `json:"number"`
		Ids         uuid.UUIDs `json:"ids"`
	}

	allProd := make([]response, 0, len(products))
	for _, prod := range products {
		allProd = append(allProd, response{
			Name:        prod.Name,
			Price:       prod.Price,
			RentPointID: prod.RentPointID,
			Status:      string(prod.Status),
			Ids:         prod.IDs,
		})
	}

	return ctx.Status(http.StatusOK).JSON(allProd)
}

func (p *productRoutes) getByID(ctx *fiber.Ctx) error {

	return nil
}
