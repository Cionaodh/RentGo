package usecase

import (
	"EasyRentGo/internal/entity"
	"context"

	"github.com/google/uuid"
)

// Интерфейс для взаимодействия с бизнес слоем
type Rent interface {
}

type Product interface {
}

type RentPointUseCase interface {
	Create(context.Context, entity.RentPoint) (entity.RentPoint, error)
	GetAll(context.Context) ([]entity.RentPoint, error)
	GetByID(context.Context, uuid.UUID) (entity.RentPoint, error)
}

type TemplateUseCase interface {
	Create(context.Context, entity.ProductTemp) (entity.ProductTemp, error)
	GetAll(context.Context) ([]entity.ProductTemp, error)
	GetByID(context.Context, uuid.UUID) (entity.ProductTemp, error)
}

type ProductUseCase interface {
}
