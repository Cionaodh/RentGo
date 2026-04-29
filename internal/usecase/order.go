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

type OrderUsecase struct {
	repo repo.Order
	l    logger.Interface
}

func NewOrderUsecase(rp repo.Order, l logger.Interface) *OrderUsecase {
	return &OrderUsecase{rp, l}
}

func (o *OrderUsecase) Create(ctx context.Context, in CreateOrderInput) (entity.Order, error) {
	if in.ProductID == uuid.Nil {
		return entity.Order{}, ErrInvalidProductID
	}
	if in.UserID == uuid.Nil {
		return entity.Order{}, ErrInvalidUserID
	}

	order, err := o.repo.Create(ctx, repotype.CreateOrderInput{
		ProductID: in.ProductID,
		UserID:    in.UserID, // Передаем в БД
	})
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrProductNotAvailable):
			return entity.Order{}, ErrProductAlreadyAssigned // Продукт уже занят
		case errors.Is(err, repoerrors.ErrProductStateInvalid), errors.Is(err, repoerrors.ErrNotFound):
			return entity.Order{}, ErrProductNotFound // Продукта не существует
		}

		o.l.Error("OrderUsecase - Create: %v", err)
		return entity.Order{}, ErrCreateOrder
	}

	return order, nil
}

func (o *OrderUsecase) GetAll(ctx context.Context) ([]entity.Order, error) {

	orders, err := o.repo.GetAll(ctx)
	if err != nil {
		o.l.Error("OrderUsecase - GetAll: %v", err)
		return nil, ErrFetchOrders
	}

	return orders, nil
}

func (o *OrderUsecase) GetByID(ctx context.Context, id uuid.UUID) (entity.Order, error) {
	if id == uuid.Nil {
		return entity.Order{}, ErrInvalidOrderID
	}

	order, err := o.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return entity.Order{}, ErrOrderNotFound
		}
		o.l.Error("OrderUsecase - GetByID: %v", err)
		return entity.Order{}, ErrFetchOrders
	}

	return order, nil
}

func (o *OrderUsecase) Complete(ctx context.Context, in CompleteOrderInput) (entity.Order, error) {
	if in.ID == uuid.Nil {
		return entity.Order{}, ErrInvalidOrderID
	}
	if in.FinishingPointID == uuid.Nil {
		return entity.Order{}, ErrInvalidRentpointID
	}

	// _, err := o.GetByID(ctx, in.UserID, in.ID)
	// if err != nil {
	// 	return entity.Order{}, err // Если не найдено или чужой заказ (ErrForbidden)
	// }

	order, err := o.repo.Complete(ctx, repotype.CompleteOrderInput{
		ID:               in.ID,
		FinishingPointID: in.FinishingPointID,
		// Можно прокинуть UserID и сюда, если в БД нужна двойная проверка
	})
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrOrderNotActive):
			return entity.Order{}, ErrOrderNotActive
		case errors.Is(err, repoerrors.ErrProductStateInvalid):
			return entity.Order{}, ErrInvalidOrderState
		case errors.Is(err, repoerrors.ErrForeignKeyViolation):
			return entity.Order{}, ErrRentPointNotFound
		}

		o.l.Error("OrderUsecase - Complete: %v", err)
		return entity.Order{}, ErrCompleteOrder
	}

	return order, nil
}
