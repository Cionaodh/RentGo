package v1

import (
	"EasyRentGo/internal/usecase"

	"EasyRentGo/pkg/logger"

	"github.com/go-playground/validator/v10"
)

// V1 -.
type V1 struct {
	rp usecase.RentPoint
	l  logger.Interface
	v  *validator.Validate
}
