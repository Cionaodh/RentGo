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

type ProductParams struct {
	Status      *entity.ProductStatus
	RentPointID *uuid.UUID
	TemplateID  *uuid.UUID
	IDs         *uuid.UUIDs
}

type AddToRentPointInput struct {
	RentPointID uuid.UUID
	ProductIDs  uuid.UUIDs
}
