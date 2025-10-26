package pgdb

import (
	"EasyRentGo/internal/entity"
	rp "EasyRentGo/internal/repo/repotypes"
	"EasyRentGo/pkg/postgres"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ProductTemplateRepo struct {
	*postgres.Postgres
}

func NewProductTemplateRepo(pg *postgres.Postgres) *ProductTemplateRepo {
	return &ProductTemplateRepo{pg}
}

func (pt *ProductTemplateRepo) Create(ctx context.Context, temp rp.CreateTemplateInput) (entity.ProductTemp, error) {
	sql := `
		INSERT INTO templates (name, description, price)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, price;
	`

	var template entity.ProductTemp
	if err := pt.Pool.QueryRow(ctx, sql,
		temp.Name,
		temp.Description,
		temp.Price,
	).Scan(
		&template.ID,
		&template.Name,
		&template.Description,
		&template.Price,
	); err != nil {
		return entity.ProductTemp{}, fmt.Errorf("ProductTemplateRepo - Create - pt.Pool.QueryRow: %w", err)
	}

	return template, nil
}

func (pt *ProductTemplateRepo) GetAll(ctx context.Context) ([]entity.ProductTemp, error) {
	sql := `
		SELECT * 
		FROM templates;`

	rows, err := pt.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("ProductTemplateRepo - GetAll - pt.Pool.Query: %w", err)
	}
	defer rows.Close()

	var templates []entity.ProductTemp
	for rows.Next() {
		var temp entity.ProductTemp
		err := rows.Scan(&temp.ID, &temp.Name, &temp.Description, &temp.Price)
		if err != nil {
			return nil, fmt.Errorf("ProductTemplateRepo - GetAll - rows.Scan: %w", err)
		}
		templates = append(templates, temp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ProductTemplateRepo - GetAll - rows.Err: %w", err)
	}

	return templates, nil
}

func (pt *ProductTemplateRepo) GetByID(context.Context, uuid.UUID) (entity.ProductTemp, error) {
	return entity.ProductTemp{}, nil
}
