package repotype

import (
	"EasyRentGo/internal/entity"

	"github.com/google/uuid"
)

type CreateProductInput struct {
	TemplateId uuid.UUID
	Status     entity.ProductStatus
	Number     int
}
