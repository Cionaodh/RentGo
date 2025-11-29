package v1

import (
	"EasyRentGo/internal/controller/http/v1/response"
	"errors"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrInvalidParameters = errors.New("invalid request parameters")
)

func errorResponse(ctx *fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(response.ErrorResp{Error: msg})
}
