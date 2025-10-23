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

	rows, err := r.Pool.Query(ctx, sql,
		rp.Name,
		rp.Addr,
	)
	// _, err := r.Pool.Exec(ctx, sql, rp.Name, rp.Addr)
	if err != nil {
		return entity.RentPoint{}, fmt.Errorf("RentpointRepo - Create - r.Pool.Query: %w", err)
	}

	point, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[entity.RentPoint])
	if err != nil {
		return entity.RentPoint{}, fmt.Errorf("pgdb.RentPoint - Create - pgx.CollectExactlyOneRow: %w", err)
	}

	return point, nil
}

// GetAll - Получение всех точек проката (без прикрепленных продуктов).
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
func (r *RentpointRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.RentPoint, error) {
	// // Получаем основную информацию о точке проката
	sql := `
		SELECT id, name, addr
		FROM rentpoints
		WHERE id = $1
	`

	var rp entity.RentPoint
	err := r.Pool.QueryRow(ctx, sql, id).Scan(&rp.ID, &rp.Name, &rp.Addr)
	if err != nil {
		return entity.RentPoint{}, fmt.Errorf("RentpointRepo - GetByID - r.Pool.QueryRow: %w", err)
	}

	// // Получаем список продуктов, привязанных к точке проката
	// products, err := r.getPointProducts(ctx, id)
	// if err != nil {
	// 	return entity.RentPoint{}, fmt.Errorf("RentpointRepo - GetByID - getPointProducts: %w", err)
	// }

	// // rp.Products = products
	// return rp, nil
	return rp, nil
}

// Delete - Удаление точки проката.
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

// getPointProducts - получение списка продуктов точки проката
func (r *RentpointRepo) getPointProducts(ctx context.Context, pointID uuid.UUID) ([]uuid.UUID, error) {
	sql := `
		SELECT product_id
		FROM rent_point_products
		WHERE rent_point_id = $1
		ORDER BY product_id
	`

	rows, err := r.Pool.Query(ctx, sql, pointID)
	if err != nil {
		return nil, fmt.Errorf("getPointProducts - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var products []uuid.UUID
	for rows.Next() {
		var productID uuid.UUID
		err := rows.Scan(&productID)
		if err != nil {
			return nil, fmt.Errorf("getPointProducts - rows.Scan: %w", err)
		}
		products = append(products, productID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getPointProducts - rows.Err: %w", err)
	}

	return products, nil
}

// updatePointProducts - обновление списка продуктов точки проката
func (r *RentpointRepo) updatePointProducts(ctx context.Context, pointID uuid.UUID, products []uuid.UUID) error {
	// Удаляем старые связи
	deleteSQL := `
		DELETE FROM rent_point_products
		WHERE rent_point_id = $1
	`
	_, err := r.Pool.Exec(ctx, deleteSQL, pointID)
	if err != nil {
		return fmt.Errorf("updatePointProducts - delete: %w", err)
	}

	// Добавляем новые связи
	if len(products) == 0 {
		return nil
	}

	// Строим запрос для множественной вставки
	insertSQL := `
		INSERT INTO rent_point_products (rent_point_id, product_id)
		VALUES 
	`
	args := make([]interface{}, 0, len(products)*2)
	argCounter := 1

	for i, productID := range products {
		if i > 0 {
			insertSQL += ","
		}
		insertSQL += fmt.Sprintf("($%d, $%d)", argCounter, argCounter+1)
		args = append(args, pointID, productID)
		argCounter += 2
	}

	_, err = r.Pool.Exec(ctx, insertSQL, args...)
	if err != nil {
		return fmt.Errorf("updatePointProducts - insert: %w", err)
	}

	return nil
}
