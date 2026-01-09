package repotype

import (
	"github.com/google/uuid"
)

type CreateOrderInput struct {
	ProductID uuid.UUID
}

type CompleteOrderInput struct {
	ID               uuid.UUID
	FinishingPointID uuid.UUID
}
