package rentpoint

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo"
	"EasyRentGo/internal/usecase"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type RentPoint struct {
	repo      repo.RentPointRepo
	productUC usecase.ProductUseCase
}

func New(rp repo.RentPointRepo, pUC usecase.ProductUseCase) *RentPoint {
	return &RentPoint{
		repo:      rp,
		productUC: pUC,
	}
}

// SetProductUC: установка связи с ProductUseCase после создания объекта
func (rp *RentPoint) SetProductUC(pUC usecase.ProductUseCase) {
	rp.productUC = pUC
}

// func (rp *RentPoint) Create(ctx context.Context, point entity.RentPoint) (entity.RentPoint, error) {
// 	point.ID = uuid.New() // генерируем уникальный id

// 	if err := rp.repo.Create(ctx, point); err != nil {
// 		return entity.RentPoint{}, fmt.Errorf("RentPointUseCase - Create - rp.repo.Create: %w", err)
// 	}

// 	return point, nil
// }

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

func (rp *RentPoint) AddProduct(ctx context.Context, idRP string, idP string) error {

}

// func (rp *RentPoint) AddProducts(ctx context.Context, pointID uuid.UUID, productIDs []uuid.UUID) error {
// 	// 1. Получаем точку проката
// 	point, err := rp.repo.GetByID(ctx, pointID)
// 	if err != nil { return err }

// 	// 2. Проверяем, что продукты доступны (не привязаны к другой точке)
// 	availableProducts, err := rp.productUC.GetAvailableProducts(ctx) // Зависимость от ProductUseCase
// 	if err != nil { return err }

// 	// 3. Фильтруем только разрешенные продукты
// 	var validProducts []uuid.UUID
// 	for _, prodID := range productIDs {
// 		 if contains(availableProducts, prodID) {
// 			  validProducts = append(validProducts, prodID)
// 		 }
// 	}

// 	// 4. Добавляем продукты в RentPoint
// 	point.Products = append(point.Products, validProducts...)
// 	if err := rp.repo.Update(ctx, point); err != nil { return err }

// 	// 5. Обновляем RentPointID у продуктов (вызов ProductUseCase)
// 	for _, prodID := range validProducts {
// 		 if err := rp.productUC.SetRentPoint(ctx, prodID, pointID); err != nil {
// 			  // Можно добавить откат или логирование
// 			  return fmt.Errorf("failed to update product %v: %w", prodID, err)
// 		 }
// 	}

// 	return nil
// }
