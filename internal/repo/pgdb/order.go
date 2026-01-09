package pgdb

import (
	"EasyRentGo/internal/entity"
	repotype "EasyRentGo/internal/repo/repotypes"
	"EasyRentGo/pkg/postgres"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type OrderRepo struct {
	*postgres.Postgres
}

func NewOrderRepo(pg *postgres.Postgres) *OrderRepo {
	return &OrderRepo{pg}
}

func (o *OrderRepo) Create(ctx context.Context, in repotype.CreateOrderInput) (entity.Order, error) {

	var order entity.Order

	err := pgx.BeginTxFunc(ctx, o.Pool, pgx.TxOptions{}, func(tx pgx.Tx) error {

		const queryRentPointID = `
		SELECT rentpoint_id 
		FROM public.products 
		WHERE 
			id = $1 AND 
			status = 'Free';
		`

		var startPointID uuid.UUID
		err := tx.QueryRow(ctx, queryRentPointID, in.ProductID).Scan(&startPointID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return err // entity.ErrProductNotAvailable // TODO:
			}
			return fmt.Errorf("rent product: %w", err)
		}

		const reserveProductQuery = `
			UPDATE public.products
         SET status = 'Rented',
            rentpoint_id = NULL
         WHERE id = $1
            AND status = 'Free'
            AND rentpoint_id IS NOT NULL
		`

		cmdTag, err := tx.Exec(ctx, reserveProductQuery, in.ProductID)
		if err != nil {
			return fmt.Errorf("err: %w", err)
		}

		// Проверяем сколько строк было изменено
		if cmdTag.RowsAffected() == 0 {
			return fmt.Errorf("ErrProductStateInvalid") // entity.ErrProductStateInvalid // TODO:
		}

		// 2. Создаём заказ
		const sqlInsertOrder = `
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

		return tx.QueryRow(
			ctx,
			sqlInsertOrder,
			in.ProductID,
			startPointID,
		).Scan(
			&order.ID,
			&order.Status,
			&order.ProductID,
			&order.StartPointID,
			&order.FigishPointID,
			&order.StartedAT,
			&order.FinishedAT,
		)
	})

	if err != nil {
		return entity.Order{}, err
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

	orders := make([]entity.Order, 0, 32)
	for rows.Next() {
		var order entity.Order
		err := rows.Scan(
			&order.ID,
			&order.Status,
			&order.ProductID,
			&order.StartPointID,
			&order.FigishPointID,
			&order.StartedAT,
			&order.FinishedAT,
		)
		if err != nil {
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
		WHERE id = $1;
	`

	var order entity.Order
	if err := o.Pool.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.Status,
		&order.ProductID,
		&order.StartPointID,
		&order.FigishPointID,
		&order.StartedAT,
		&order.FinishedAT,
	); err != nil {
		return entity.Order{}, err
	}

	return order, nil
}

func (o *OrderRepo) Complete(ctx context.Context, in repotype.CompleteOrderInput) (entity.Order, error) {

	var order entity.Order

	err := pgx.BeginTxFunc(ctx, o.Pool, pgx.TxOptions{}, func(tx pgx.Tx) error {

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
				return err // entity.ErrOrderNotActive // TODO:
			}
			return fmt.Errorf("complete order: %w", err)
		}

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
			return fmt.Errorf("free product: %w", err)
		}

		// Проверяем сколько строк было изменено
		if cmdTag.RowsAffected() == 0 {
			return fmt.Errorf("ErrProductStateInvalid") // entity.ErrProductStateInvalid // TODO:
		}

		return nil
	})

	if err != nil {
		return entity.Order{}, err
	}

	return order, nil
}

// TODO: Добавить ф-ю получение id заказа по id продукта
// func (o *OrderRepo) GetByProductID(ctx context.Context, id uuid.UUID) (entity.Order, error){}
