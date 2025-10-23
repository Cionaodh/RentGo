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
		return entity.Products{}, fmt.Errorf("failed to create products: %w", err)
	}

	return product, nil
}

// TODO: написать JOIN с таблицей шаблонов
func (p *ProductRepo) GetAll(ctx context.Context) ([]entity.Products, error) {
	sql := `
        SELECT 
            template_id,
            rentpoint_id,
            status,
            ARRAY_AGG(id) AS ids
        FROM products
        GROUP BY template_id, rentpoint_id, status
        ORDER BY template_id, rentpoint_id, status
    `

	rows, err := p.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}
	defer rows.Close()

	var products []entity.Products
	for rows.Next() {
		var product entity.Products
		err := rows.Scan(
			&product.TemplateID,
			&product.RentPointID,
			&product.Status,
			&product.IDs,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return products, nil
}

func (p *ProductRepo) GetByID(context.Context, uuid.UUID) (entity.Products, error) {
	return entity.Products{}, nil
}

func (p *ProductRepo) GetProductsByRentpoint(context.Context, uuid.UUID) ([]entity.Products, error) {
	// TODO: Находим все продуткы принадлежащие пункту проката по id
	// TODO: JOIN с таблицей шаблонов (чтобы выводилась вся информация по продукту)
	return []entity.Products{}, nil
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
