package usecase

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo"
	repotype "EasyRentGo/internal/repo/repotypes"
	"EasyRentGo/pkg/logger"
	"context"

	"github.com/google/uuid"
)

type OrderUsecase struct {
	repo repo.Order
	l    logger.Interface
}

func NewOrderUsecase(rp repo.Order, l logger.Interface) *OrderUsecase {
	return &OrderUsecase{rp, l}
}

func (o *OrderUsecase) Create(ctx context.Context, in CreateOrderInput) (entity.Order, error) {

	order, err := o.repo.Create(ctx, repotype.CreateOrderInput{
		ProductID: in.ProductID,
	})
	if err != nil {
		return entity.Order{}, err
	}

	return order, nil
}

func (o *OrderUsecase) GetAll(ctx context.Context) ([]entity.Order, error) {

	orders, err := o.repo.GetAll(ctx)
	if err != nil {
		return []entity.Order{}, nil
	}

	return orders, nil
}

func (o *OrderUsecase) GetByID(ctx context.Context, id uuid.UUID) (entity.Order, error) {

	order, err := o.repo.GetByID(ctx, id)
	if err != nil {
		return entity.Order{}, err
	}

	return order, nil
}

func (o *OrderUsecase) Complete(ctx context.Context, in CompleteOrderInput) (entity.Order, error) {

	return o.repo.Complete(ctx, repotype.CompleteOrderInput{
		ID:               in.ID,
		FinishingPointID: in.FinishingPointID,
	})
}
