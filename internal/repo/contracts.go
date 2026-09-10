package repo

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo/pgdb"
	repotype "EasyRentGo/internal/repo/repotypes"
	rp "EasyRentGo/internal/repo/repotypes"
	"EasyRentGo/pkg/postgres"
	"context"

	"github.com/google/uuid"
)

type User interface {
	Create(ctx context.Context, input repotype.CreateUserInput) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	GetProfile(ctx context.Context, id uuid.UUID) (entity.UserProfile, error)
}

type RentPoint interface {
	Create(context.Context, rp.CreateRentpointInput) (entity.RentPoint, error)
	GetAll(context.Context) ([]entity.RentPoint, error)
	// Delete(context.Context, uuid.UUID) error
	AddProducts(context.Context, rp.AddProductsInput) (entity.ProductRentPoint, error) // Только добавляет продукт
	GetByID(context.Context, uuid.UUID) (entity.RentPoint, error)
}

type TemplateProduct interface {
	Create(context.Context, rp.CreateTemplateInput) (entity.ProductTemplate, error)
	GetAll(context.Context) ([]entity.ProductTemplate, error)
	GetByID(context.Context, uuid.UUID) (entity.ProductTemplate, error)
}

type Product interface {
	Create(context.Context, rp.CreateProductInput) (entity.Products, error)
	GetAll(context.Context) ([]entity.Product, error)
	GetByID(context.Context, uuid.UUID) (entity.Product, error)
	GetByRentpoint(context.Context, uuid.UUID) ([]entity.ProductsRP, error)
	List(context.Context, rp.ProductParams) ([]entity.ProductsRP, error)

	// GetByParameters() ([]entity.Products, error)

	// GetByStatus(context.Context, entity.ProductStatus) ([]entity.Product, error)
	// SetStatus(context.Context, uuid.UUID, entity.ProductStatus) error
	// SetRentPoint(context.Context, uuid.UUID, uuid.UUID) (entity.Product, error)
	// RentProduct(context.Context, uuid.UUID) (entity.Product, error)
	// ReserveProduct(context.Context, uuid.UUID) (entity.Product, error)
}

type Order interface {
	Create(context.Context, rp.CreateOrderInput) (entity.Order, error)
	GetAll(context.Context) ([]entity.Order, error)
	GetByID(context.Context, uuid.UUID) (entity.Order, error)
	Complete(context.Context, rp.CompleteOrderInput) (entity.Order, error)
}

type Repositories struct {
	RentPoint
	TemplateProduct
	Product
	Order
	User
}

func NewPostgresRepo(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		RentPoint:       pgdb.NewRentpointRepo(pg),
		TemplateProduct: pgdb.NewProductTemplateRepo(pg),
		Product:         pgdb.NewProductRepo(pg),
		Order:           pgdb.NewOrderRepo(pg),
		User:            pgdb.NewUserRepo(pg),
	}
}
