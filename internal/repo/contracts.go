package repo

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo/pgdb"
	rp "EasyRentGo/internal/repo/repotypes"
	"EasyRentGo/pkg/postgres"
	"context"

	"github.com/google/uuid"
)

type RentPoint interface {
	Create(context.Context, rp.CreateRentpointInput) (entity.RentPoint, error)
	GetAll(context.Context) ([]entity.RentPoint, error)
	GetByID(context.Context, uuid.UUID) (entity.RentPoint, error)
	Delete(context.Context, uuid.UUID) error
	// TODO: AddPriduct(context.Context, uuid.UUID) error// Прикрепление продукта к точке проката
}

type TemplateProduct interface {
	Create(context.Context, rp.CreateTemplateInput) (entity.ProductTemp, error)
	GetAll(context.Context) ([]entity.ProductTemp, error)
	GetByID(context.Context, uuid.UUID) (entity.ProductTemp, error)
}

type Product interface {
	Create(context.Context, rp.CreateProductInput) (entity.Products, error)
	GetAll(context.Context) ([]entity.ProductsTemp, error)
	GetByID(context.Context, uuid.UUID) (entity.Products, error)
	GetProductsByRentpoint(context.Context, uuid.UUID) ([]entity.Products, error)

	// GetByStatus(context.Context, entity.ProductStatus) ([]entity.Product, error)
	// SetStatus(context.Context, uuid.UUID, entity.ProductStatus) error
	// SetRentPoint(context.Context, uuid.UUID, uuid.UUID) (entity.Product, error)
	// RentProduct(context.Context, uuid.UUID) (entity.Product, error)
	// ReserveProduct(context.Context, uuid.UUID) (entity.Product, error)
}

type Repositories struct {
	RentPoint
	TemplateProduct
	Product
}

func NewPostgresRepo(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		RentPoint:       pgdb.NewRentpointRepo(pg),
		TemplateProduct: pgdb.NewProductTemplateRepo(pg),
		Product:         pgdb.NewProductRepo(pg),
	}
}
