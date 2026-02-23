package v1

import (
	"EasyRentGo/internal/controller/http/v1/response"
	"errors"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrInvalidParameters  = errors.New("invalid request parameters")
	ErrInvalidRequestBody = errors.New("invalid request body")
)

func errorResponse(ctx *fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(response.ErrorResp{Error: msg})
}
