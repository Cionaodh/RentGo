package pgdb

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo/repoerrors"
	rp "EasyRentGo/internal/repo/repotypes"
	"EasyRentGo/pkg/postgres"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ProductTemplateRepo struct {
	*postgres.Postgres
}

func NewProductTemplateRepo(pg *postgres.Postgres) *ProductTemplateRepo {
	return &ProductTemplateRepo{pg}
}

func (pt *ProductTemplateRepo) Create(ctx context.Context, temp rp.CreateTemplateInput) (entity.ProductTemplate, error) {
	sql := `
		INSERT INTO templates (name, description, price)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, price;
	`

	var template entity.ProductTemplate

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

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return entity.ProductTemplate{}, repoerrors.ErrAlreadyExists
			}
		}

		return entity.ProductTemplate{}, fmt.Errorf("ProductTemplateRepo - Create - pt.Pool.QueryRow: %w", err)
	}

	return template, nil
}

func (pt *ProductTemplateRepo) GetAll(ctx context.Context) ([]entity.ProductTemplate, error) {
	sql := `
		SELECT id, name, description, price
		FROM templates;
		`

	rows, err := pt.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("ProductTemplateRepo - GetAll - pt.Pool.Query: %w", err)
	}
	defer rows.Close()

	var templates []entity.ProductTemplate
	for rows.Next() {
		var temp entity.ProductTemplate
		if err := rows.Scan(
			&temp.ID,
			&temp.Name,
			&temp.Description,
			&temp.Price,
		); err != nil {
			return nil, fmt.Errorf("ProductTemplateRepo - GetAll - rows.Scan: %w", err)
		}

		templates = append(templates, temp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ProductTemplateRepo - GetAll - rows.Err: %w", err)
	}

	return templates, nil
}

func (pt *ProductTemplateRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.ProductTemplate, error) {
	sql := `
		SELECT id, name, description, price
		FROM templates
		WHERE id = $1
	`

	var tmp entity.ProductTemplate
	err := pt.Pool.QueryRow(ctx, sql, id).Scan(&tmp.ID, &tmp.Name, &tmp.Description, &tmp.Price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ProductTemplate{}, repoerrors.ErrNotFound
		}

		return entity.ProductTemplate{}, fmt.Errorf("ProductTemplateRepo - GetByID - p.Pool.QueryRow.Scan: %w", err)
	}

	return tmp, nil
}
