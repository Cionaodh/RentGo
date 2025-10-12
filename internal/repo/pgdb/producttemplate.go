package pgdb

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/pkg/postgres"
	"context"

	"github.com/google/uuid"
)

type ProductTemplateRepo struct {
	*postgres.Postgres
}

func NewProductTemplateRepo(pg *postgres.Postgres) *ProductTemplateRepo {
	return &ProductTemplateRepo{pg}
}

func (pt *ProductTemplateRepo) Create(context.Context, entity.ProductTemp) error {
	return nil
}

func (pt *ProductTemplateRepo) GetAll(context.Context) ([]entity.ProductTemp, error) {
	return []entity.ProductTemp{}, nil
}

func (pt *ProductTemplateRepo) GetByID(context.Context, uuid.UUID) (entity.ProductTemp, error) {
	return entity.ProductTemp{}, nil
}
