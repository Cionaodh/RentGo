package usecase

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo"
	repotype "EasyRentGo/internal/repo/repotypes"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type RentPointUsecase struct {
	pointRepo   repo.RentPoint
	productRepo repo.Product
}

func NewRentPointUsecase(rp repo.RentPoint, p repo.Product) *RentPointUsecase {
	return &RentPointUsecase{
		pointRepo:   rp,
		productRepo: p,
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

// Возвращает все данные о точке проката. И все подукты, которые принадлежат точке проката (вызывает getProducts)
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

// func (rp *RentPointUsecase) AddProduct(ctx context.Context, idRP string, idP string) error {
// 	return nil
// }

// func (rp *RentPoint) AddProducts(ctx context.Context, pointID uuid.UUID, productIDs []uuid.UUID) error {
// 	// 1. Получаем точку проката
// 	point, err := rp.repo.GetByID(ctx, pointID)
// 	if err != nil { return err }
//
// 	// 2. Проверяем, что продукты доступны (не привязаны к другой точке)
// 	availableProducts, err := rp.productUC.GetAvailableProducts(ctx) // Зависимость от ProductUseCase
// 	if err != nil { return err }
//
// 	// 3. Фильтруем только разрешенные продукты
// 	var validProducts []uuid.UUID
// 	for _, prodID := range productIDs {
// 		 if contains(availableProducts, prodID) {
// 			  validProducts = append(validProducts, prodID)
// 		 }
// 	}
//
// 	// 4. Добавляем продукты в RentPoint
// 	point.Products = append(point.Products, validProducts...)
// 	if err := rp.repo.Update(ctx, point); err != nil { return err }
//
// 	// 5. Обновляем RentPointID у продуктов (вызов ProductUseCase)
// 	for _, prodID := range validProducts {
// 		 if err := rp.productUC.SetRentPoint(ctx, prodID, pointID); err != nil {
// 			  // Можно добавить откат или логирование
// 			  return fmt.Errorf("failed to update product %v: %w", prodID, err)
// 		 }
// 	}
//
// 	return nil
// }

// func getProducts(context.Context, uuid.UUID) ([]entity.Product, error){} // получение всех продуктов точки проката
// 1. Залезаем в таблицу продуктов
// 2. Выводим все продукты, которые принадлежат выбранному
