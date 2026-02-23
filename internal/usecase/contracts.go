package usecase

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo"
	"EasyRentGo/pkg/logger"
	"context"

	"github.com/google/uuid"
)

type CreateRentpointInput struct {
	Name string
	Addr string
}

type AddProductsInput struct {
	RentpointID uuid.UUID
	ProductsID  uuid.UUIDs
}

type RentPoint interface {
	CreateRentpoint(context.Context, CreateRentpointInput) (entity.RentPoint, error)
	GetAll(context.Context) ([]entity.RentPoint, error)
	GetByID(context.Context, uuid.UUID) (entity.ProductRentPoint, error)
	Delete(context.Context, uuid.UUID) error
	AddProduct(context.Context, AddProductsInput) (entity.ProductRentPoint, error)
}

type CreateTemplateInput struct {
	Name        string
	Description string
	Price       float64
}

type Template interface {
	Create(context.Context, CreateTemplateInput) (entity.ProductTemplate, error)
	GetAll(context.Context) ([]entity.ProductTemplate, error)
	GetByID(context.Context, uuid.UUID) (entity.ProductTemplate, error)
}

type CreateProductInput struct {
	TemplateId uuid.UUID
	Number     int
}

type ProductParamsInput struct {
	Status      *entity.ProductStatus
	RentPointID *uuid.UUID
	TemplateID  *uuid.UUID
}

type Product interface {
	Create(context.Context, CreateProductInput) (entity.Products, error)
	GetAll(context.Context) ([]entity.Product, error)
	GetByID(context.Context, uuid.UUID) (entity.Product, error)
	Delete(context.Context, uuid.UUID) error
	List(context.Context, ProductParamsInput) ([]entity.ProductsRP, error)

	// GetByStatus(context.Context, entity.ProductStatus) ([]entity.Product, error)
	// GetByRentPointID(context.Context, uuid.UUID) ([]entity.Product, error)
	// SetStatus(context.Context, uuid.UUID, entity.ProductStatus) (entity.Product, error)
	// SetRentPoint(context.Context, uuid.UUID, uuid.UUID) (entity.Product, error)
	// Reserve(context.Context, uuid.UUID) (entity.Product, error)
	// RentStart(context.Context, uuid.UUID) error
	// RentFinish(context.Context, uuid.UUID, uuid.UUID) error
}

type CreateOrderInput struct {
	ProductID uuid.UUID
}

type CompleteOrderInput struct {
	ID               uuid.UUID
	FinishingPointID uuid.UUID
}

type Order interface {
	Create(context.Context, CreateOrderInput) (entity.Order, error)
	GetAll(context.Context) ([]entity.Order, error)
	GetByID(context.Context, uuid.UUID) (entity.Order, error)
	Complete(context.Context, CompleteOrderInput) (entity.Order, error)
}

type Usecases struct {
	RentPoint
	Template
	Product
	Order
}

type UsecaseDependencies struct {
	Repos *repo.Repositories
}

func NewUsecase(d UsecaseDependencies, l logger.Interface) *Usecases {
	return &Usecases{
		RentPoint: NewRentPointUsecase(d.Repos.RentPoint, d.Repos.Product, l),
		Template:  NewTemplateUsecase(d.Repos.TemplateProduct, l),
		Product:   NewProductUsecase(d.Repos.Product, d.Repos.TemplateProduct, l),
		Order:     NewOrderUsecase(d.Repos.Order, l),
	}
}
