package pgdb

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/pkg/postgres"
	"context"

	"github.com/google/uuid"
)

type ProductRepo struct {
	*postgres.Postgres
}

func NewProductRepo(pg *postgres.Postgres) *ProductRepo {
	return &ProductRepo{pg}
}

func (p *ProductRepo) Create(context.Context, entity.Product) error {
	return nil
}

func (p *ProductRepo) GetAll(context.Context) ([]entity.Product, error) {
	return []entity.Product{}, nil
}

func (p *ProductRepo) GetByID(context.Context, uuid.UUID) (entity.Product, error) {
	return entity.Product{}, nil
}

func (p *ProductRepo) GetByStatus(context.Context, entity.ProductStatus) ([]entity.Product, error) {
	return []entity.Product{}, nil
}

func (p *ProductRepo) SetStatus(context.Context, uuid.UUID, entity.ProductStatus) error {
	return nil
}

func (p *ProductRepo) SetRentPoint(context.Context, uuid.UUID, uuid.UUID) (entity.Product, error) {
	return entity.Product{}, nil
}

func (p *ProductRepo) RentProduct(context.Context, uuid.UUID) (entity.Product, error) {
	return entity.Product{}, nil
}

func (p *ProductRepo) ReserveProduct(context.Context, uuid.UUID) (entity.Product, error) {
	return entity.Product{}, nil
}
