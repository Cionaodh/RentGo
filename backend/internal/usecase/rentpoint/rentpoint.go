package rentpoint

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type RentPoint struct {
	repo repo.RentPointRepo // указание на репу (например, repo.RentPointRepo)
	// указание на внешние интерфейсы (например, RentPointWebAPI)
}

func New(rp repo.RentPointRepo) *RentPoint {
	return &RentPoint{
		repo: rp,
	}
}

func (rp *RentPoint) Create(ctx context.Context, point entity.RentPoint) (entity.RentPoint, error) {
	point.ID = uuid.New() // генерируем уникальный id

	if err := rp.repo.Create(ctx, point); err != nil {
		return entity.RentPoint{}, fmt.Errorf("RentPointUseCase - Create - rp.repo.Create: %w", err)
	}

	return point, nil
}

func (rp *RentPoint) GetAll(ctx context.Context) ([]entity.RentPoint, error) {
	rentPoints, err := rp.repo.GetAll(ctx)
	if err != nil {
		return []entity.RentPoint{}, fmt.Errorf("RentPointUseCase - GetAll - rp.repo.GetAll: %w", err)
	}

	return rentPoints, nil
}

func (rp *RentPoint) GetByID(ctx context.Context, id uuid.UUID) (entity.RentPoint, error) {

	point, err := rp.repo.GetByID(ctx, id)
	if err != nil {
		return entity.RentPoint{}, fmt.Errorf("RentPointUseCase - GetByID - rp.repo.GetByDI: %w", err)
	}

	return point, nil
}
