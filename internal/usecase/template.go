package usecase

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo"
	"EasyRentGo/internal/repo/repoerrors"
	repotype "EasyRentGo/internal/repo/repotypes"
	"EasyRentGo/pkg/logger"
	"context"
	"errors"

	"github.com/google/uuid"
)

type TemplateUsecase struct {
	repo repo.TemplateProduct
	l    logger.Interface
}

func NewTemplateUsecase(r repo.TemplateProduct, l logger.Interface) *TemplateUsecase {
	return &TemplateUsecase{
		repo: r,
		l:    l,
	}
}

func (uc *TemplateUsecase) Create(ctx context.Context, t CreateTemplateInput) (entity.ProductTemplate, error) {
	// Проверка входных данных
	if t.Name == "" {
		return entity.ProductTemplate{}, ErrInvalidTemplateName
	}
	if t.Price <= 0 {
		return entity.ProductTemplate{}, ErrInvalidTemplatePrice
	}

	temp, err := uc.repo.Create(ctx, repotype.CreateTemplateInput{
		Name:        t.Name,
		Description: t.Description,
		Price:       t.Price,
	})
	if err != nil {
		if errors.Is(err, repoerrors.ErrAlreadyExists) {
			return entity.ProductTemplate{}, ErrTemplateAlreadyExists
		}
		uc.l.Error("TemplateUsecase - Create: %v", err)
		return entity.ProductTemplate{}, ErrCreateTemplate
	}

	return temp, nil
}

func (uc *TemplateUsecase) GetAll(ctx context.Context) ([]entity.ProductTemplate, error) {
	templates, err := uc.repo.GetAll(ctx)
	if err != nil {
		uc.l.Error("TemplateUsecase - GetAll: %v", err)
		return nil, ErrGetAllTemplate
	}

	// if len(templates) == 0 {
	// 	return nil, errors.New("no templates found")
	// }

	return templates, nil
}

func (uc *TemplateUsecase) GetByID(ctx context.Context, id uuid.UUID) (entity.ProductTemplate, error) {
	if id == uuid.Nil {
		return entity.ProductTemplate{}, ErrInvalidTemplateID
	}

	template, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return entity.ProductTemplate{}, ErrTemplateNotFound
		}
		uc.l.Error("TemplateUsecase - GetByID: %v", err)
		return entity.ProductTemplate{}, ErrFetchTemplates
	}

	return template, nil
}
