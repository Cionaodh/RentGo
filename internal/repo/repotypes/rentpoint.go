package repotype

import (
	"EasyRentGo/internal/entity"

	"github.com/google/uuid"
)

type CreateRentpointInput struct {
	Name string
	Addr string
}

type AddProductsInput struct {
	ID_rentpoint uuid.UUID
	IDs_products uuid.UUIDs
	Status       entity.ProductStatus
}
