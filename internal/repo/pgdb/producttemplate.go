package pgdb

import (
	"EasyRentGo/internal/entity"
	rp "EasyRentGo/internal/repo/repotypes"
	"EasyRentGo/pkg/postgres"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ProductTemplateRepo struct {
	*postgres.Postgres
}

func NewProductTemplateRepo(pg *postgres.Postgres) *ProductTemplateRepo {
	return &ProductTemplateRepo{pg}
}

func (pt *ProductTemplateRepo) Create(ctx context.Context, temp rp.CreateTemplateInput) (entity.ProductTemp, error) {
	// Добавить строку, которая содержит
	// (имя, описание, Цену)

	sql := `
		INSERT INTO template (name, description, price)
		VALUES ($1, $2, $3)
		RETURNING *;
	`

	rows, err := pt.Pool.Query(ctx, sql,
		temp.Name,
		temp.Description,
		temp.Price,
	)
	if err != nil {
		return entity.ProductTemp{}, fmt.Errorf("ProductTemplateRepo - Create - pt.Pool.Query: %w", err)
	}

	template, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[entity.ProductTemp])
	if err != nil {
		return entity.ProductTemp{}, fmt.Errorf("ProductTemplateRepo - Create - pgx.CollectExactlyOneRow: %w", err)
	}

	return template, nil
}

func (pt *ProductTemplateRepo) GetAll(context.Context) ([]entity.ProductTemp, error) {
	return []entity.ProductTemp{}, nil
}

func (pt *ProductTemplateRepo) GetByID(context.Context, uuid.UUID) (entity.ProductTemp, error) {
	return entity.ProductTemp{}, nil
}
