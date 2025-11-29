package usecase

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo"
	repotype "EasyRentGo/internal/repo/repotypes"
	"EasyRentGo/pkg/logger"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type RentPointUsecase struct {
	pointRepo   repo.RentPoint
	productRepo repo.Product
	l           logger.Interface
}

func NewRentPointUsecase(rp repo.RentPoint, p repo.Product, l logger.Interface) *RentPointUsecase {
	return &RentPointUsecase{
		pointRepo:   rp,
		productRepo: p,
		l:           l,
	}
}

func (rp *RentPointUsecase) CreateRentpoint(ctx context.Context, in CreateRentpointInput) (entity.RentPoint, error) {
	// point.ID = uuid.New()

	point, err := rp.pointRepo.Create(ctx, repotype.CreateRentpointInput{
		Name: in.Name,
		Addr: in.Addr,
	})
	if err != nil {
		return entity.RentPoint{}, fmt.Errorf("RentPointUseCase - Create - rp.repo.Create: %w", err)
	}

	return point, nil
}

func (rp *RentPointUsecase) GetAll(ctx context.Context) ([]entity.RentPoint, error) {
	rentPoints, err := rp.pointRepo.GetAll(ctx)
	if err != nil {
		return []entity.RentPoint{}, fmt.Errorf("RentPointUseCase - GetAll - rp.repo.GetAll: %w", err)
	}

	return rentPoints, nil
}

func (rp *RentPointUsecase) GetByID(ctx context.Context, id uuid.UUID) (entity.ProductRentPoint, error) {

	point, err := rp.pointRepo.GetByID(ctx, id)
	if err != nil {
		return entity.ProductRentPoint{}, fmt.Errorf("RentPointUseCase - GetByID - rp.repo.GetByDI: %w", err)
	}

	return point, nil
}

func (rp *RentPointUsecase) Delete(context.Context, uuid.UUID) error {
	return nil
}

func (p *RentPointUsecase) AddProduct(ctx context.Context, in AddProductsInput) (entity.ProductRentPoint, error) {
	rentpoint, err := p.pointRepo.AddProducts(ctx, repotype.AddProductsInput{
		ID_rentpoint: in.ID_rentpoint,
		IDs_products: in.IDs_products,
		Status:       entity.StatusFree,
	})
	if err != nil {
		return entity.ProductRentPoint{}, fmt.Errorf("RentPointUsecase - AddProduct - p.pointRepo.AddProducts: %w", err)
	}
	return rentpoint, nil
	// TODO: Добавить проверку, что продукты доступны (не привязаны к другой точке)
}
