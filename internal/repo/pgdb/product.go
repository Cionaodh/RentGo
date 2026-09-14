package pgdb

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo/repoerrors"
	rp "EasyRentGo/internal/repo/repotypes"
	"EasyRentGo/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ProductRepo struct {
	*postgres.Postgres
}

func NewProductRepo(pg *postgres.Postgres) *ProductRepo {
	return &ProductRepo{pg}
}

func (p *ProductRepo) Create(ctx context.Context, in rp.CreateProductInput) (entity.Products, error) {

	sql := `
			WITH inserted AS (
    			INSERT INTO products (template_id, status)
    			SELECT $1, $2
    			FROM generate_series(1, $3)
    			RETURNING id, status, template_id, rentpoint_id
			)
			SELECT 
    			status,
   			template_id,
    			rentpoint_id,
    			ARRAY_AGG(id) AS ids
			FROM inserted
			GROUP BY status, template_id, rentpoint_id;
			`

	var product entity.Products
	err := p.Pool.QueryRow(ctx, sql, in.TemplateId, in.Status, in.Number).Scan(
		&product.Status,
		&product.TemplateID,
		&product.RentPointID,
		&product.IDs,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503": // если шаблона не существует
				return entity.Products{}, repoerrors.ErrNotFound
			}
		}
		return entity.Products{}, fmt.Errorf("ProductRepo - Create - QueryRow()Scan(): %w", err)
	}

	return product, nil
}

func (p *ProductRepo) GetAll(ctx context.Context) ([]entity.Product, error) {
	sql := `
	SELECT 
		 p.id,
		 t.name AS name,
	    t.price AS price,
	    p.status,
	    p.rentpoint_id
	FROM products p
	JOIN templates t ON p.template_id = t.id
	ORDER BY p.rentpoint_id, p.template_id, p.status;
    `

	rows, err := p.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("ProductRepo - GetAll - Query(): %w", err)
	}
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		var product entity.Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Status,
			&product.RentPointID,
		); err != nil {
			return nil, fmt.Errorf("ProductRepo - GetAll - Scan(): %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ProductRepo - GetAll - rows.Err(): %w", err)
	}

	return products, nil
}

func (p *ProductRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.Product, error) {

	sql := `
	SELECT 
	 p.id,
	 t.name AS name,
    t.price AS price,
    p.status,
    p.rentpoint_id
	FROM products p
	JOIN templates t ON p.template_id = t.id
	WHERE p.id = $1;
	`

	var product entity.Product
	if err := p.Pool.QueryRow(ctx, sql, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Status,
		&product.RentPointID,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Product{}, repoerrors.ErrNotFound
		}
		return entity.Product{}, fmt.Errorf("ProductRepo - GetByID - p.Pool.QueryRow: %w", err)
	}

	return product, nil
}

func (p *ProductRepo) GetByRentpoint(ctx context.Context, id uuid.UUID) ([]entity.ProductsRP, error) {

	sqlGet := `
		SELECT p.id, t.name, p.status, t.price
		FROM products p
		LEFT JOIN templates t ON p.template_id = t.id
		WHERE p.id = $1
	`

	rows, err := p.Pool.Query(ctx, sqlGet, id)
	if err != nil {
		return nil, fmt.Errorf("ProductRepo - GetByRentpoint - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var products []entity.ProductsRP
	for rows.Next() {
		var prod entity.ProductsRP
		err := rows.Scan(&prod.Id, &prod.Name, &prod.Status, &prod.Price)
		if err != nil {
			return nil, fmt.Errorf("ProductRepo - GetByRentpoint - rows.Scan: %w", err)
		}
		products = append(products, prod)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ProductRepo - GetByRentpoint - rows.Err: %w", err)
	}

	return products, nil
}

func (p *ProductRepo) List(ctx context.Context, in rp.ProductParams) ([]entity.ProductsRP, error) {

	baseQuery := `
        SELECT 
            p.id,
            t.name,
            p.status,
            t.price
        FROM public.products p
        JOIN public.templates t ON p.template_id = t.id
    `

	var (
		conditions []string
		args       []any
		argID      = 1
	)

	if in.Status != nil {
		conditions = append(conditions, fmt.Sprintf("p.status = $%d", argID))
		args = append(args, *in.Status)
		argID++
	}

	if in.TemplateID != nil {
		conditions = append(conditions, fmt.Sprintf("p.template_id = $%d", argID))
		args = append(args, *in.TemplateID)
		argID++
	}

	if in.RentPointID != nil {
		conditions = append(conditions, fmt.Sprintf("p.rentpoint_id = $%d", argID))
		args = append(args, *in.RentPointID)
		argID++
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := p.Pool.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("ProductRepo - List - Query(): %w", err)
	}
	defer rows.Close()

	var products []entity.ProductsRP
	for rows.Next() {
		var product entity.ProductsRP
		if err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Status,
			&product.Price,
		); err != nil {
			return nil, fmt.Errorf("ProductRepo - List - Scan(): %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ProductRepo - List - rows.Err(): %w", err)
	}

	return products, nil
}

// DetachByRentPointID - открепляем все продукты от точки проката
func (p *ProductRepo) DetachByRentPointID(ctx context.Context, rentPointID uuid.UUID) error {
	const query = `
		UPDATE products
		SET status = $2, rentpoint_id = NULL
		WHERE rentpoint_id = $1`

	if _, err := p.Pool.Exec(ctx, query, rentPointID, entity.StatusUnused); err != nil {
		return fmt.Errorf("ProductRepo - DetachByRentPointID - p.Pool.Exec: %w", err)
	}

	return nil
}

// TODO: Возвращаем обновлен ли продукт {при завершении аренды меняется и статус и точка проката}
// func (p *ProductRepo) UpdateStatusById(ctx context.Context, id uuid.UUID) error {
// 	return false, nil
// }

// func (p *ProductRepo) GetByStatus(context.Context, entity.ProductStatus) ([]entity.Product, error) {
// 	return []entity.Product{}, nil
// }

// func (p *ProductRepo) SetStatus(context.Context, uuid.UUID, entity.ProductStatus) error {
// 	return nil
// }

// func (p *ProductRepo) SetRentPoint(context.Context, uuid.UUID, uuid.UUID) (entity.Product, error) {
// 	return entity.Product{}, nil
// }

// func (p *ProductRepo) RentProduct(context.Context, uuid.UUID) (entity.Product, error) {
// 	return entity.Product{}, nil
// }

// func (p *ProductRepo) ReserveProduct(context.Context, uuid.UUID) (entity.Product, error) {
// 	return entity.Product{}, nil
// }
