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

type RentPointRepo struct {
	*postgres.Postgres
}

func NewRentpointRepo(pg *postgres.Postgres) *RentPointRepo {
	return &RentPointRepo{pg}
}

func (r *RentPointRepo) Create(ctx context.Context, in rp.CreateRentpointInput) (entity.RentPoint, error) {
	sql := `
        INSERT INTO rentpoints (name, addr)
        VALUES ($1, $2)
        RETURNING id, name, addr;
    `

	var point entity.RentPoint
	err := r.Pool.QueryRow(ctx, sql, in.Name, in.Addr).
		Scan(&point.ID, &point.Name, &point.Addr)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return entity.RentPoint{}, repoerrors.ErrAlreadyExists
			}
		}
		return entity.RentPoint{}, fmt.Errorf("RentpointRepo - Create - QueryRow().Scan(): %w", err)
	}

	return point, nil
}

func (r *RentPointRepo) GetAll(ctx context.Context) ([]entity.RentPoint, error) {
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
func (r *RentPointRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.RentPoint, error) {
	// TODO: Проверяем удалена ли точка или существует

	sql := `
	SELECT id, name, addr
	FROM rentpoints
	WHERE id = $1;`

	var rp entity.RentPoint
	err := r.Pool.QueryRow(ctx, sql, id).Scan(&rp.ID, &rp.Name, &rp.Addr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.RentPoint{}, repoerrors.ErrNotFound
		}
		return entity.RentPoint{}, fmt.Errorf("RentpointRepo - GetByID - r.Pool.QueryRow: %w", err)
	}
	return rp, nil
}

func (r *RentPointRepo) Delete(ctx context.Context, rentPointID uuid.UUID) error {

	const query = `
		UPDATE rentpoints
		SET is_deleted = TRUE
		WHERE id = $1 AND is_deleted = FALSE`

	tag, err := r.Pool.Exec(ctx, query, rentPointID)
	if err != nil {
		return fmt.Errorf("RentpointRepo - Delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return repoerrors.ErrRentPointNotFound // возвращаем 404
	}

	return nil
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
