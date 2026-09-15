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

type ProductUsecase struct {
	productRepo  repo.Product
	templateRepo repo.TemplateProduct
	l            logger.Interface
}

func NewProductUsecase(p repo.Product, t repo.TemplateProduct, l logger.Interface) *ProductUsecase {
	return &ProductUsecase{p, t, l}
}

func (p *ProductUsecase) Create(ctx context.Context, in CreateProductInput) (entity.Products, error) {

	if in.TemplateId == uuid.Nil {
		return entity.Products{}, ErrInvalidTemplateID
	}

	if in.Number <= 0 || in.Number > 100000 {
		return entity.Products{}, ErrInvalidProductQty
	}

	products, err := p.productRepo.Create(ctx, repotype.CreateProductInput{
		TemplateId: in.TemplateId,
		Status:     entity.StatusUnused,
		Number:     in.Number,
	})
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNotFound):
			return entity.Products{}, ErrTemplateNotFound
		default:
			p.l.Error("ProductUsecase - Create(): %v", err)
			return entity.Products{}, ErrCreateProduct
		}
	}

	return products, nil
}

func (p *ProductUsecase) GetAll(ctx context.Context) ([]entity.Product, error) {
	products, err := p.productRepo.GetAll(ctx)
	if err != nil {
		p.l.Error("ProductUsecase.GetAll: %v", err)
		return nil, ErrFetchProducts
	}
	return products, nil
}

func (p *ProductUsecase) GetByID(ctx context.Context, id uuid.UUID) (entity.Product, error) {
	if id == uuid.Nil {
		return entity.Product{}, ErrInvalidProductID
	}

	product, err := p.productRepo.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNotFound):
			return entity.Product{}, ErrProductNotFound
		default:
			p.l.Error("ProductUsecase - GetByID: %v", err)
			return entity.Product{}, ErrFetchProducts
		}
	}

	return product, nil
}

func (p *ProductUsecase) List(ctx context.Context, in ProductParamsInput) ([]entity.ProductsRP, error) {

	products, err := p.productRepo.List(ctx, repotype.ProductParams{
		Status:      in.Status,
		RentPointID: in.RentPointID,
		TemplateID:  in.TemplateID,
	})
	if err != nil {
		p.l.Error("ProductUsecase - List: %v", err)
		return nil, ErrFetchProducts
	}

	return products, nil
}

// ReturnToRentPoint возвращает арендованный продукт на точку проката по завершении аренды
func (p *ProductUsecase) ReturnToRentPoint(ctx context.Context, in AddProductsInput) error {
	if in.RentPointID == uuid.Nil {
		return ErrRentPointNotFound
	}

	if len(in.ProductIDs) != 1 {
		return ErrSingleProductRequired
	}

	ids, err := p.productRepo.AddToRentPoint(ctx, repotype.AddToRentPointInput{
		RentPointID: in.RentPointID,
		ProductIDs:  in.ProductIDs,
	}, entity.StatusRented)
	if err != nil {
		if errors.Is(err, repoerrors.ErrRentPointNotFound) {
			return ErrRentPointNotFound
		}
		p.l.Error("ProductUsecase - ReturnToRentPoint - AddToRentPoint: %v", err)
		return ErrAttachProducts
	}

	if len(ids) != len(in.ProductIDs) {
		return ErrProductNotFound
	}

	return nil
}

// MultipleAddToRentPoint прикрепление множества продуктов к точке проката
func (p *ProductUsecase) MultipleAddToRentPoint(ctx context.Context, in AddProductsInput) ([]entity.ProductsRP, error) {
	if in.RentPointID == uuid.Nil {
		return nil, ErrRentPointNotFound //
	}

	if len(in.ProductIDs) == 0 {
		return nil, ErrEmptyProductIDs
	}

	ids, err := p.productRepo.AddToRentPoint(ctx, repotype.AddToRentPointInput{
		RentPointID: in.RentPointID,
		ProductIDs:  in.ProductIDs,
	}, entity.StatusUnused)
	if err != nil {
		if errors.Is(err, repoerrors.ErrRentPointNotFound) {
			return nil, ErrRentPointNotFound
		}
		p.l.Error("ProductUsecase.MultipleAddToRentPoint: %v", err)
		return nil, ErrAttachProducts
	}

	// либо прикреплены все, либо ни один
	if len(ids) != len(in.ProductIDs) {
		return nil, ErrProductsNotAvailable
	}
	// TODO: добавить поле, в котором будут указаны продукты, которые не удалось прикрепить

	products, err := p.productRepo.List(ctx, repotype.ProductParams{
		IDs: &ids,
	})
	if err != nil {
		p.l.Error("ProductUsecase - MultipleAddToRentPoint: %v", err)
		return nil, ErrFetchProducts
	}

	return products, nil
}

func (p *ProductUsecase) Delete(context.Context, uuid.UUID) error {
	return nil
}
