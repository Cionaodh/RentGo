package v1

import (
	"EasyRentGo/internal/usecase"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Controller - контроллер для API пользователей (v1)
type userRoutes struct {
	user usecase.User // интерфейс
	v    *validator.Validate
}

func newUserRoutes(user usecase.User, v *validator.Validate) *userRoutes {
	return &userRoutes{user, v}
}

type AuthDTO struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"secret123"`
}

func (r *userRoutes) register(ctx *fiber.Ctx) error {
	var body AuthDTO

	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, ErrInvalidRequestBody.Error())
	}

	if err := r.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid email or password format")
	}

	// Генерация username (user_0001) происходит внутри usecase
	user, err := r.user.Register(ctx.UserContext(), usecase.RegisterUserInput{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrUserAlreadyExists):
			return errorResponse(ctx, http.StatusConflict, err.Error())
		default:
			return errorResponse(ctx, http.StatusInternalServerError, "failed to register user")
		}
	}

	return ctx.Status(http.StatusCreated).JSON(user)
}

func (r *userRoutes) login(ctx *fiber.Ctx) error {
	var body AuthDTO

	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, ErrInvalidRequestBody.Error())
	}

	if err := r.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid inputs")
	}

	user, err := r.user.Login(ctx.UserContext(), usecase.LoginUserInput{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidCredentials),
			errors.Is(err, usecase.ErrUserNotFound):
			return errorResponse(ctx, http.StatusUnauthorized, "invalid email or password")
		default:
			return errorResponse(ctx, http.StatusInternalServerError, "failed to login")
		}
	}

	return ctx.Status(http.StatusOK).JSON(user)
}

// getProfile - Получение профиля и заказов
func (r *userRoutes) getProfile(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "user id is required")
	}

	profile, err := r.user.GetProfile(ctx.UserContext(), userID)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrUserNotFound):
			return errorResponse(ctx, http.StatusNotFound, err.Error())
		default:
			return errorResponse(ctx, http.StatusInternalServerError, "failed to fetch user profile")
		}
	}

	return ctx.Status(http.StatusOK).JSON(profile)
}
