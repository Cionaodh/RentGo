package pgdb

import (
	"EasyRentGo/internal/entity"
	rp "EasyRentGo/internal/repo/repotypes"
	"EasyRentGo/pkg/postgres"
	"context"
	"fmt"

	"github.com/google/uuid"
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
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Status,
			&product.RentPointID,
		)
		if err != nil {
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
	err := p.Pool.QueryRow(ctx, sql, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Status,
		&product.RentPointID,
	)
	if err != nil {
		return entity.Product{}, fmt.Errorf("ProductRepo - GetByID - p.Pool.QueryRow: %w", err)
	}

	return product, nil
}

func (p *ProductRepo) GetByRentpoint(ctx context.Context, id uuid.UUID) ([]entity.ProductsRP, error) {

	sqlGet := `
		SELECT p.id, t,name, p.status, t.price
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
