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
	"github.com/jackc/pgx/v5/pgconn"
)

// RentpointRepo -.
type RentpointRepo struct {
	*postgres.Postgres
}

// New -.
func NewRentpointRepo(pg *postgres.Postgres) *RentpointRepo {
	return &RentpointRepo{pg}
}

// Create - Создание новой точки проката.
func (r *RentpointRepo) Create(ctx context.Context, rp rp.CreateRentpointInput) (entity.RentPoint, error) {
	sql := `
        INSERT INTO rentpoints (name, addr)
        VALUES ($1, $2)
        RETURNING id, name, addr;
    `

	var point entity.RentPoint
	err := r.Pool.QueryRow(ctx, sql, rp.Name, rp.Addr).Scan(&point.ID, &point.Name, &point.Addr)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entity.RentPoint{}, repoerrors.ErrRentPointAlreadyExists
		}
		return entity.RentPoint{}, fmt.Errorf("RentpointRepo - Create: %w", err)
	}

	return point, nil
}

func (r *RentpointRepo) GetAll(ctx context.Context) ([]entity.RentPoint, error) {
	sql := `
		SELECT id, name, addr
		FROM rentpoints;`

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("RentpointRepo - GetAll - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var rentPoints []entity.RentPoint
	for rows.Next() {
		var rp entity.RentPoint
		err := rows.Scan(&rp.ID, &rp.Name, &rp.Addr)
		if err != nil {
			return nil, fmt.Errorf("RentpointRepo - GetAll - rows.Scan: %w", err)
		}
		rentPoints = append(rentPoints, rp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("RentpointRepo - GetAll - rows.Err: %w", err)
	}

	return rentPoints, nil
}

// GetByID - Получение точки проката по id (со списком продуктов).
func (r *RentpointRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.ProductRentPoint, error) {

	sql := `
	SELECT 
	    r.id,
	    r.name,
	    r.addr,
	    COALESCE(
	        JSON_AGG(
	            JSON_BUILD_OBJECT(
	                'id', p.id,
	                'name', t.name,
	                'price', t.price,
	                'status', p.status
	            )
	        ) FILTER (WHERE p.id IS NOT NULL),
	        '[]'
	    ) AS products
	FROM rentpoints r
	LEFT JOIN products p ON r.id = p.rentpoint_id
	LEFT JOIN templates t ON p.template_id = t.id
	WHERE r.id = $1
	GROUP BY r.id, r.name, r.addr;
	`

	var rp entity.ProductRentPoint
	if err := r.Pool.QueryRow(ctx, sql, id).Scan(&rp.ID, &rp.Name, &rp.Addr, &rp.Products); err != nil {
		return entity.ProductRentPoint{}, fmt.Errorf("RentpointRepo - GetByID - r.Pool.QueryRow: %w", err)
	}

	return rp, nil

	// TODO: разбить sql запрос в рамках одной транзакции
	// TODO: добавить обработку ошибки NotFound
}

func (r *RentpointRepo) Delete(ctx context.Context, id uuid.UUID) error {
	// Сначала удаляем связи с продуктами
	deleteProductsSQL := `
		DELETE FROM rent_point_products
		WHERE rent_point_id = $1
	`
	_, err := r.Pool.Exec(ctx, deleteProductsSQL, id)
	if err != nil {
		return fmt.Errorf("RentpointRepo - Delete - delete products: %w", err)
	}

	// Затем удаляем саму точку проката
	deleteSQL := `
		DELETE FROM rent_points
		WHERE id = $1
	`

	result, err := r.Pool.Exec(ctx, deleteSQL, id)
	if err != nil {
		return fmt.Errorf("RentpointRepo - Delete - r.Pool.Exec: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("RentpointRepo - Delete - point not found")
	}

	return nil
}

func (r *RentpointRepo) AddProducts(ctx context.Context, in rp.AddProductsInput) (entity.ProductRentPoint, error) {
	sqlUpdate := `
	UPDATE products
	SET 
	    rentpoint_id = $1,
	    status = $2
	WHERE id = ANY($3);
	`

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return entity.ProductRentPoint{}, fmt.Errorf("RentpointRepo - AddProducts - r.Pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, sqlUpdate, in.ID_rentpoint, in.Status, in.IDs_products)
	if err != nil {
		return entity.ProductRentPoint{}, fmt.Errorf("RentpointRepo - AddProducts - tx.Exec: %w", err)
	}

	sqlGet := `
	SELECT 
	    r.id,
	    r.name,
	    r.addr,
	    COALESCE(
	        JSON_AGG(
	            JSON_BUILD_OBJECT(
	                'id', p.id,
	                'name', t.name,
	                'price', t.price,
	                'status', p.status
	            )
	        ) FILTER (WHERE p.id IS NOT NULL),
	        '[]'
	    ) AS products
	FROM rentpoints r
	LEFT JOIN products p ON r.id = p.rentpoint_id
	LEFT JOIN templates t ON p.template_id = t.id
	WHERE r.id = $1
	GROUP BY r.id, r.name, r.addr;
	`

	var rp entity.ProductRentPoint
	if err = tx.QueryRow(ctx, sqlGet, in.ID_rentpoint).Scan(&rp.ID, &rp.Name, &rp.Addr, &rp.Products); err != nil {
		return entity.ProductRentPoint{}, fmt.Errorf("RentpointRepo - AddProducts - r.Pool.QueryRow: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return entity.ProductRentPoint{}, fmt.Errorf("RentpointRepo - AddProducts - tx.Commit: %w", err)
	}

	return rp, nil
}

// getPointProducts - получение списка продуктов точки проката
// func (r *RentpointRepo) getPointProducts(ctx context.Context, pointID uuid.UUID) ([]uuid.UUID, error) {
// 	sql := `
// 		SELECT product_id
// 		FROM rent_point_products
// 		WHERE rent_point_id = $1
// 		ORDER BY product_id
// 	`
//
// 	rows, err := r.Pool.Query(ctx, sql, pointID)
// 	if err != nil {
// 		return nil, fmt.Errorf("getPointProducts - r.Pool.Query: %w", err)
// 	}
// 	defer rows.Close()
//
// 	var products []uuid.UUID
// 	for rows.Next() {
// 		var productID uuid.UUID
// 		err := rows.Scan(&productID)
// 		if err != nil {
// 			return nil, fmt.Errorf("getPointProducts - rows.Scan: %w", err)
// 		}
// 		products = append(products, productID)
// 	}
//
// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("getPointProducts - rows.Err: %w", err)
// 	}
//
// 	return products, nil
// }

// updatePointProducts - обновление списка продуктов точки проката
// func (r *RentpointRepo) updatePointProducts(ctx context.Context, pointID uuid.UUID, products []uuid.UUID) error {
// 	// Удаляем старые связи
// 	deleteSQL := `
// 		DELETE FROM rent_point_products
// 		WHERE rent_point_id = $1
// 	`
// 	_, err := r.Pool.Exec(ctx, deleteSQL, pointID)
// 	if err != nil {
// 		return fmt.Errorf("updatePointProducts - delete: %w", err)
// 	}
//
// 	// Добавляем новые связи
// 	if len(products) == 0 {
// 		return nil
// 	}
//
// 	// Строим запрос для множественной вставки
// 	insertSQL := `
// 		INSERT INTO rent_point_products (rent_point_id, product_id)
// 		VALUES
// 	`
// 	args := make([]interface{}, 0, len(products)*2)
// 	argCounter := 1
//
// 	for i, productID := range products {
// 		if i > 0 {
// 			insertSQL += ","
// 		}
// 		insertSQL += fmt.Sprintf("($%d, $%d)", argCounter, argCounter+1)
// 		args = append(args, pointID, productID)
// 		argCounter += 2
// 	}
//
// 	_, err = r.Pool.Exec(ctx, insertSQL, args...)
// 	if err != nil {
// 		return fmt.Errorf("updatePointProducts - insert: %w", err)
// 	}
//
// 	return nil
// }
