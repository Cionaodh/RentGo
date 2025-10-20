package usecase

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo"
	repotype "EasyRentGo/internal/repo/repotypes"
	"context"
	"errors"

	"github.com/google/uuid"
)

type TeplateProductUsecase struct {
	repo repo.TemplateProduct
}

func NewTeplateProductUsecase(r repo.TemplateProduct) *TeplateProductUsecase {
	return &TeplateProductUsecase{
		repo: r,
	}
}

func (uc *TeplateProductUsecase) Create(ctx context.Context, t CreateTemplateInput) (entity.ProductTemp, error) {
	// Валидация входных данных
	if t.Name == "" {
		return entity.ProductTemp{}, errors.New("invalid template name")
	}
	if t.Price <= 0 {
		return entity.ProductTemp{}, errors.New("invalid template price")
	}

	temp, err := uc.repo.Create(ctx, repotype.CreateTemplateInput{
		Name:        t.Name,
		Description: t.Description,
		Price:       t.Price,
	})
	if err != nil {
		return entity.ProductTemp{}, err
	}

	return temp, nil
}

func (uc *TeplateProductUsecase) GetAll(ctx context.Context) ([]entity.ProductTemp, error) {
	templates, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(templates) == 0 {
		return nil, errors.New("no templates found")
	}

	return templates, nil
}

func (uc *TeplateProductUsecase) GetByID(ctx context.Context, id uuid.UUID) (entity.ProductTemp, error) {
	if id == uuid.Nil {
		return entity.ProductTemp{}, errors.New("invalid template ID")
	}

	template, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return entity.ProductTemp{}, err
	}

	return template, nil
}
