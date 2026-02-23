package pgdb

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo/repoerrors"
	repotype "EasyRentGo/internal/repo/repotypes"
	"EasyRentGo/pkg/postgres"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type OrderRepo struct {
	*postgres.Postgres
}

func NewOrderRepo(pg *postgres.Postgres) *OrderRepo {
	return &OrderRepo{pg}
}

func (o *OrderRepo) Create(ctx context.Context, in repotype.CreateOrderInput) (entity.Order, error) {

	tx, err := o.Pool.Begin(ctx)
	if err != nil {
		return entity.Order{}, fmt.Errorf("OrderRepo - Create - BeginTX(): %w", err)
	}
	defer tx.Rollback(ctx)

	var startPointID uuid.UUID

	// проверяем что продукт свободен и получаем его точку проката
	const selectProductQuery = `
		SELECT rentpoint_id 
		FROM public.products 
		WHERE id = $1 
			AND status = 'Free'
		FOR UPDATE;
		`

	err = tx.QueryRow(ctx, selectProductQuery, in.ProductID).Scan(&startPointID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Order{}, repoerrors.ErrProductNotAvailable // продукта не существует
		}
		return entity.Order{}, fmt.Errorf("OrderRepo - Create - QueryRow(): %w", err)
	}

	// Резервируем продукт если он
	const reserveProductQuery = `
			UPDATE public.products
         SET status = 'Rented',
            rentpoint_id = NULL
         WHERE id = $1; 
		`

	cmdTag, err := tx.Exec(ctx, reserveProductQuery, in.ProductID)
	if err != nil {
		return entity.Order{}, fmt.Errorf("OrderRepo - Create - Exec(): %w", err)
	}

	// Проверяем была ли изменена строка
	if cmdTag.RowsAffected() == 0 {
		return entity.Order{}, repoerrors.ErrProductStateInvalid // строка не была изменена
	}

	// Создаём заказ
	const insertOrderQuery = `
            INSERT INTO public.orders (
                product_id,
                starting_point_id
            )
            VALUES ($1, $2)
            RETURNING
                id,
                status,
                product_id,
                starting_point_id,
                finishing_point_id,
                rent_started_at,
                rent_finished_at;
      `

	var order entity.Order

	err = tx.QueryRow(ctx, insertOrderQuery, in.ProductID, startPointID).Scan(
		&order.ID,
		&order.Status,
		&order.ProductID,
		&order.StartPointID,
		&order.FigishPointID,
		&order.StartedAT,
		&order.FinishedAT,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return entity.Order{}, repoerrors.ErrForeignKeyViolation // если внешний ключ не существует
			}
		}
		return entity.Order{}, fmt.Errorf("OrderRepo - Create - QueryRow(): %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return entity.Order{}, fmt.Errorf("OrderRepo - Create - Commit(): %w", err)
	}

	return order, nil
}

func (o *OrderRepo) GetAll(ctx context.Context) ([]entity.Order, error) {

	const query = `
    SELECT 
        id, 
        status, 
        product_id, 
        starting_point_id,
        finishing_point_id,
        rent_started_at,
        rent_finished_at
    FROM public.orders
	 ORDER BY rent_started_at DESC;
	`

	rows, err := o.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - GetAll - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(
			&order.ID,
			&order.Status,
			&order.ProductID,
			&order.StartPointID,
			&order.FigishPointID,
			&order.StartedAT,
			&order.FinishedAT,
		); err != nil {
			return nil, fmt.Errorf("OrderRepo - GetAll - rows.Scan: %w", err)
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("OrderRepo - GetAll - rows.Err: %w", err)
	}

	return orders, nil
}

func (o *OrderRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.Order, error) {

	const selectOrderQuery = `
		SELECT 
			id, 
        status, 
        product_id, 
        starting_point_id,
        finishing_point_id,
        rent_started_at,
        rent_finished_at
		FROM public.orders
		WHERE id = $1;
	`

	var order entity.Order
	if err := o.Pool.QueryRow(ctx, selectOrderQuery, id).Scan(
		&order.ID,
		&order.Status,
		&order.ProductID,
		&order.StartPointID,
		&order.FigishPointID,
		&order.StartedAT,
		&order.FinishedAT,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Order{}, repoerrors.ErrNotFound // заказ не найден
		}
		return entity.Order{}, fmt.Errorf("OrderRepo - GetByID - QueryRow(): %w", err)
	}

	return order, nil
}

func (o *OrderRepo) Complete(ctx context.Context, in repotype.CompleteOrderInput) (entity.Order, error) {

	tx, err := o.Pool.Begin(ctx)
	if err != nil {
		return entity.Order{}, fmt.Errorf("OrderRepo - Complete - BeginTX(): %w", err)
	}
	defer tx.Rollback(ctx)

	var order entity.Order

	// TODO: можно разделить запрос на несколько чтобы проверить отдельно существование заказа и его статус
	// закрытие заказа
	const completeOrderQuery = `
        UPDATE public.orders
        SET
            status = 'Completed',
            finishing_point_id = $1,
            rent_finished_at = now()
        WHERE id = $2
          AND status = 'Active'
        RETURNING
            id,
            status,
            product_id,
            starting_point_id,
            finishing_point_id,
            rent_started_at,
            rent_finished_at;
    `

	if err := tx.QueryRow(
		ctx,
		completeOrderQuery,
		in.FinishingPointID,
		in.ID,
	).Scan(
		&order.ID,
		&order.Status,
		&order.ProductID,
		&order.StartPointID,
		&order.FigishPointID,
		&order.StartedAT,
		&order.FinishedAT,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Order{}, repoerrors.ErrOrderNotActive // либо заказа не существует, либо у него неверный статус
		}
		return entity.Order{}, fmt.Errorf("OrderRepo - Complete - QueryRow(): %w", err)
	}

	// освобождение продукта
	const freeProductQuery = `
        UPDATE public.products
        SET
            rentpoint_id = $1,
            status = 'Free'
        WHERE id = $2
          AND status = 'Rented';
    `

	cmdTag, err := tx.Exec(
		ctx,
		freeProductQuery,
		in.FinishingPointID,
		order.ProductID,
	)
	if err != nil {
		return entity.Order{}, fmt.Errorf("OrderRepo - Complete - Exec(): %w", err)
	}

	// Проверяем сколько строк было изменено
	if cmdTag.RowsAffected() == 0 {
		return entity.Order{}, repoerrors.ErrProductStateInvalid // строка не была изменена
	}

	if err := tx.Commit(ctx); err != nil {
		return entity.Order{}, fmt.Errorf("OrderRepo - Complete - Commit(): %w", err)
	}

	return order, nil
}

// TODO: Добавить ф-ю получение id заказа по id продукта
// func (o *OrderRepo) GetByProductID(ctx context.Context, id uuid.UUID) (entity.Order, error){}
