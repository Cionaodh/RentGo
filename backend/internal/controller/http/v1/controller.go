package v1

import (
	"EasyRentGo/internal/usecase"

	"EasyRentGo/pkg/logger"

	"github.com/go-playground/validator/v10"
)

// Controller - контроллер домена RentPoint API v1
type V1RentPointController struct {
	rp usecase.RentPointUseCase
	l  logger.Interface
	v  *validator.Validate
}

// Controller - контроллер домена ProductTemplate API v1
type V1ProductTmpController struct {
	tmp usecase.TemplateUseCase
	l   logger.Interface
	v   *validator.Validate
}

// // Controller - контроллер домена Product API v1
// type V1ProductController struct {
// 	p usecase.ProductUseCase
// 	l logger.Interface
// 	v *validator.Validate
// }
