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

	if in.Addr == "" || in.Name == "" {
		return entity.RentPoint{}, ErrFieldIsEmpty
	}
	if len(in.Addr) > 50 || len(in.Name) > 50 {
		return entity.RentPoint{}, ErrFieldIsTooLong
	}

	point, err := rp.pointRepo.Create(ctx, repotype.CreateRentpointInput{
		Name: in.Name,
		Addr: in.Addr,
	})
	if err != nil {
		if errors.Is(err, repoerrors.ErrAlreadyExists) {
			return entity.RentPoint{}, ErrRentPointAlreadyExists
		}
		rp.l.Error("RentPointUseCase - Create: %v", err)
		return entity.RentPoint{}, ErrCreateRentpoint
	}

	return point, nil
}

func (rp *RentPointUsecase) GetAll(ctx context.Context) ([]entity.RentPoint, error) {
	points, err := rp.pointRepo.GetAll(ctx)
	if err != nil {
		rp.l.Error("RentPointUseCase - GetAll: %v", err)
		return nil, ErrFetchRentPoints
	}

	return points, nil
}

func (rp *RentPointUsecase) GetByID(ctx context.Context, id uuid.UUID) (entity.ProductRentPoint, error) {
	if id == uuid.Nil {
		return entity.ProductRentPoint{}, ErrRentPointNotFound
	}

	point, err := rp.pointRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return entity.ProductRentPoint{}, ErrRentPointNotFound
		}
		rp.l.Error("RentPointUseCase - GetByID: %v", err)
		return entity.ProductRentPoint{}, err
	}

	return point, nil
}

func (rp *RentPointUsecase) Delete(context.Context, uuid.UUID) error {
	return nil
}

func (rp *RentPointUsecase) AddProduct(ctx context.Context, in AddProductsInput) (entity.ProductRentPoint, error) {
	if in.RentpointID == uuid.Nil {
		return entity.ProductRentPoint{}, ErrRentPointNotFound
	}
	if len(in.ProductsID) == 0 {
		return entity.ProductRentPoint{}, ErrProductNotFound
	}

	point, err := rp.pointRepo.AddProducts(ctx, repotype.AddProductsInput{
		ID_rentpoint: in.RentpointID,
		IDs_products: in.ProductsID,
		Status:       entity.StatusFree,
	})
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNotFound):
			return entity.ProductRentPoint{}, ErrRentPointNotFound
		case errors.Is(err, repoerrors.ErrForeignKeyViolation):
			return entity.ProductRentPoint{}, ErrProductNotFound
		default:
			rp.l.Error("RentPointUsecase - AddProduct: %v", err)
			return entity.ProductRentPoint{}, err
		}
	}

	return point, nil
	// TODO: Добавить проверку, что продукты доступны (не привязаны к другой точке)
}
