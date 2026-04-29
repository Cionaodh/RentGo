package pgdb

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo/repoerrors"
	repotype "EasyRentGo/internal/repo/repotypes"
	"EasyRentGo/pkg/postgres"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(db *postgres.Postgres) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) Create(ctx context.Context, input repotype.CreateUserInput) (entity.User, error) {
	query := `
		INSERT INTO users (email, username, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, username, password_hash, created_at`

	var user entity.User
	err := r.Pool.QueryRow(ctx, query, input.Email, input.Username, input.PasswordHash).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		// Проверяем ошибку ограничения уникальности (например, если email уже есть в базе)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entity.User{}, repoerrors.ErrAlreadyExists
		}
		return entity.User{}, err
	}

	return user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	query := `SELECT id, email, username, password_hash, created_at FROM users WHERE email = $1`

	var user entity.User
	err := r.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repoerrors.ErrNotFound
		}
		return entity.User{}, err
	}

	return user, nil
}

func (r *UserRepo) GetProfile(ctx context.Context, id uuid.UUID) (entity.UserProfile, error) {
	// 1. Достаем основные данные юзера
	queryUser := `SELECT id, email, username FROM users WHERE id = $1`
	var profile entity.UserProfile

	err := r.Pool.QueryRow(ctx, queryUser, id).Scan(&profile.ID, &profile.Email, &profile.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.UserProfile{}, repoerrors.ErrNotFound
		}
		return entity.UserProfile{}, err
	}

	// 2. Вытягиваем список заказов пользователя
	queryOrders := `
    SELECT 
        id, 
        user_id, 
        status, 
        product_id, 
        starting_point_id, 
        finishing_point_id, 
        rent_started_at, 
        rent_finished_at 
    FROM orders 
    WHERE user_id = $1`

	rows, err := r.Pool.Query(ctx, queryOrders, id)
	if err != nil {
		return entity.UserProfile{}, err
	}
	defer rows.Close()

	// Инициализируем пустой слайс, чтобы в JSON не было null, если заказов нет
	profile.Orders = make([]entity.Order, 0)

	for rows.Next() {
		var o entity.Order
		err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.Status,
			&o.ProductID,
			&o.StartPointID,
			&o.FigishPointID,
			&o.StartedAT,
			&o.FinishedAT,
		)
		if err != nil {
			return entity.UserProfile{}, err
		}
		profile.Orders = append(profile.Orders, o)
	}

	if err = rows.Err(); err != nil {
		return entity.UserProfile{}, err
	}

	return profile, nil
}
