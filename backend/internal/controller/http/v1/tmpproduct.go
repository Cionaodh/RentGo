package v1

import (
	"EasyRentGo/internal/controller/http/v1/request"
	"EasyRentGo/internal/entity"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Обработчики маршрутов

// @Summary     Create RentPoint
// @Description Create RentPoint
// @ID          create-rentpoint
// @Tags  	    rentpoint
// @Accept      json
// @Produce     json
// @Param       request body request.RentPoint true "Create RentPoint"
// @Success     200 {object} entity.RentPoint
// @Failure     400 {object} error
// @Failure     500 {object} error
// @Router      /rentpoint/ [post]
func (r *V1ProductTmpController) create(ctx *fiber.Ctx) error {
	var body request.TmpProduct

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

// @Summary     Get all RentPoints
// @Description Get all RentPoints
// @ID          get-all-rentpoints
// @Tags        rentpoint
// @Accept      json
// @Produce     json
// @Success     200 {array} entity.RentPoint
// @Failure     400 {object} error
// @Failure     500 {object} error
// @Router      /rentpoint/ [get]
func (r *V1ProductTmpController) getAll(ctx *fiber.Ctx) error {
	// points, err := r.rp.GetAll(ctx.UserContext())
	res, err := r.tmp.GetAll(ctx.UserContext())
	if err != nil {
		r.l.Error(err, "http - v1 - getAll")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to get rent points")
	}

	return ctx.Status(http.StatusOK).JSON(res)
}

// @Summary     Get RentPoint by ID
// @Description Get RentPoint by ID
// @ID          get-rentpoint-by-id
// @Tags        rentpoint
// @Accept      json
// @Produce     json
// @Param       id path string true "RentPoint ID"
// @Success     200 {object} entity.RentPoint
// @Failure     400 {object} error
// @Failure     404 {object} error
// @Failure     500 {object} error
// @Router      /rentpoint/{id} [get]
func (r *V1ProductTmpController) getByID(ctx *fiber.Ctx) error {
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
