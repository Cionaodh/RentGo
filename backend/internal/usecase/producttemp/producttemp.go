package producttemp

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo"
	"context"
	"errors"

	"github.com/google/uuid"
)

type TeplateProduct struct {
	repo repo.ProductTempRepo
}

func New(r repo.ProductTempRepo) *TeplateProduct {
	return &TeplateProduct{
		repo: r,
	}
}

func (uc *TeplateProduct) Create(ctx context.Context, t entity.ProductTemp) (entity.ProductTemp, error) {
	// Валидация входных данных
	if t.Name == "" {
		return entity.ProductTemp{}, errors.New("invalid template name")
	}
	if t.Price <= 0 {
		return entity.ProductTemp{}, errors.New("invalid template price")
	}

	// Create UUID
	t.ID = uuid.New()

	err := uc.repo.Create(ctx, t)
	if err != nil {
		return entity.ProductTemp{}, err
	}

	return t, nil
}

func (uc *TeplateProduct) GetAll(ctx context.Context) ([]entity.ProductTemp, error) {
	templates, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(templates) == 0 {
		return nil, errors.New("no templates found")
	}

	return templates, nil
}

func (uc *TeplateProduct) GetByID(ctx context.Context, id uuid.UUID) (entity.ProductTemp, error) {
	if id == uuid.Nil {
		return entity.ProductTemp{}, errors.New("invalid template ID")
	}

	template, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return entity.ProductTemp{}, err
	}

	return template, nil
}
