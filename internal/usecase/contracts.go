package usecase

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo"
	"context"

	"github.com/google/uuid"
)

type CreateRentpointInput struct {
	Name string
	Addr string
}

type RentPoint interface {
	CreateRentpoint(context.Context, CreateRentpointInput) (entity.RentPoint, error)
	GetAll(context.Context) ([]entity.RentPoint, error)
	GetByID(context.Context, uuid.UUID) (entity.ProductRentPoint, error)
	Delete(context.Context, uuid.UUID) error
	// AddProduct(context.Context, uuid.UUID, uuid.UUID)
}

type CreateTemplateInput struct {
	Name        string
	Description string
	Price       float64
}

type Template interface {
	Create(context.Context, CreateTemplateInput) (entity.ProductTemp, error)
	GetAll(context.Context) ([]entity.ProductTemp, error)
	GetByID(context.Context, uuid.UUID) (entity.ProductTemp, error)
}

type CreateProductInput struct {
	TemplateId uuid.UUID
	Number     int
}

type Product interface {
	Create(context.Context, CreateProductInput) (entity.Products, error)
	GetAll(context.Context) ([]entity.Products, error)
	GetByID(context.Context, uuid.UUID) (entity.Products, error)
	Delete(context.Context, uuid.UUID) error

	// GetByStatus(context.Context, entity.ProductStatus) ([]entity.Product, error)
	// GetByRentPointID(context.Context, uuid.UUID) ([]entity.Product, error)
	// SetStatus(context.Context, uuid.UUID, entity.ProductStatus) (entity.Product, error)
	// SetRentPoint(context.Context, uuid.UUID, uuid.UUID) (entity.Product, error)
	// Reserve(context.Context, uuid.UUID) (entity.Product, error)
	// RentStart(context.Context, uuid.UUID) error
	// RentFinish(context.Context, uuid.UUID, uuid.UUID) error
}

type Usecases struct {
	RentPoint
	Template
	Product
}

type UsecaseDependencies struct {
	Repos *repo.Repositories
}

func NewUsecase(d UsecaseDependencies) *Usecases {
	return &Usecases{
		RentPoint: NewRentPointUsecase(d.Repos.RentPoint, d.Repos.Product),
		Template:  NewTeplateProductUsecase(d.Repos.TemplateProduct),
		Product:   NewProductUsecase(d.Repos.Product, d.Repos.TemplateProduct),
	}
}
