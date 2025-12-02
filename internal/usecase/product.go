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

type ProductUsecase struct {
	productRepo  repo.Product
	templateRepo repo.TemplateProduct
	l            logger.Interface
}

func NewProductUsecase(p repo.Product, t repo.TemplateProduct, l logger.Interface) *ProductUsecase {
	return &ProductUsecase{p, t, l}
}

func (p *ProductUsecase) Create(ctx context.Context, in CreateProductInput) (entity.Products, error) {
	// проверяем количество создаваемых объектов
	if in.Number <= 0 || in.Number > 100000 {
		return entity.Products{}, ErrNumberProduct
	}

	product := repotype.CreateProductInput{
		TemplateId: in.TemplateId,
		Status:     entity.StatusUnused, // устанавливаем default status
		Number:     in.Number,
	}

	products, err := p.productRepo.Create(ctx, product)
	if err != nil {
		// TODO: обработка ожидаемых ошибок: 1) несуществующий шаблон
		p.l.Error("ProductUsecase - Create - p.productRepo.Create: %v", err)
		return entity.Products{}, ErrCreateProduct
	}

	return products, nil
}

func (p *ProductUsecase) GetAll(ctx context.Context) ([]entity.Product, error) {
	products, err := p.productRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("ProductUsecase - GetAll - p.productRepo.GetAll: %w", err)
	}

	// if len(products) == 0 {
	// 	return nil, errors.New("no products found")
	// }

	return products, nil
}

func (p *ProductUsecase) GetByID(ctx context.Context, id uuid.UUID) (entity.Product, error) {
	product, err := p.productRepo.GetByID(ctx, id)
	if err != nil {
		return entity.Product{}, fmt.Errorf("ProductUsecase - GetByID - p.productRepo.GetByID: %w", err)
	}

	// TODO: проверить является ли пустым полученное значение

	return product, nil
}

func (p *ProductUsecase) Delete(context.Context, uuid.UUID) error {
	return nil
}

// func (p *ProductUsecase) GetByStatus(ctx context.Context, status string) ([]entity.Product, error) {
// 	// Проверка корректности статуса
// 	s, err := entity.ParseProductStatus(status)
// 	if err != nil {
// 		return []entity.Product{}, fmt.Errorf("%s", err)
// 	}
//
// 	return p.repo.GetByStatus(ctx, s)
// }

// func (p *ProductUsecase) SetStatus(ctx context.Context, id string, status string) (entity.Product, error) {
// 	// Проверка корректности статуса
// 	s, err := entity.ParseProductStatus(status)
// 	if err != nil {
// 		return entity.Product{}, fmt.Errorf("%s", err)
// 	}
//
// 	// проверка uuid на корректность и преобразование
// 	if err := uuid.Validate(id); err != nil {
// 		return entity.Product{}, fmt.Errorf("%s", err)
// 	}
//
// 	if err := p.repo.SetStatus(ctx, uuid.MustParse(id), s); err != nil {
// 		return entity.Product{}, fmt.Errorf("%s", err)
// 	}
//
// 	prod, err := p.repo.GetByID(ctx, uuid.MustParse(id))
// 	if err != nil {
// 		return entity.Product{}, fmt.Errorf("%s", err)
// 	}
//
// 	return prod, nil
// }

// func (p *ProductUsecase) SetRentPoint(ctx context.Context, idProduct string, idPoint string) (entity.Product, error) {
// 	// if err := uuid.Validate(idProduct); err != nil {
// 	// 	return entity.Product{}, fmt.Errorf("%s", err)
// 	// }
//
// 	// if err := uuid.Validate(idPoint); err != nil {
// 	// 	return entity.Product{}, fmt.Errorf("%s", err)
// 	// }
//
// 	// point, err := p.repo.SetRentPoint(ctx, uuid.MustParse(idProduct), uuid.MustParse(idPoint)) // Прикрепляем к продукту - пункт проката
// 	// if err != nil {
// 	// 	return entity.Product{}, fmt.Errorf("%s", err)
// 	// }
//
// 	// p.rentPointUC.
//
// 	// // p.rentPointUC.
//
// 	// // Получаем объект точки проката
// 	// // Вызываем фукнцию добавления продукта в хранилще продуктов Пункта проката
//
//
// 	// return entity.Product, nil
// 	return entity.Product{}, nil
// }

// func (p *Product) SetRentPoint(ctx context.Context, productID, pointID string) (entity.Product, error) {
// 	// 1. Валидация ID
// 	prodUUID, err := uuid.Parse(productID)
// 	if err != nil { /* ... */ }
//
// 	pointUUID, err := uuid.Parse(pointID)
// 	if err != nil { /* ... */ }
//
// 	// 2. Обновляем RentPointID у продукта
// 	if err := p.repo.SetRentPoint(ctx, prodUUID, pointUUID); err != nil {
// 		 return entity.Product{}, fmt.Errorf("failed to update product: %w", err)
// 	}
//
// 	// 3. Добавляем продукт в RentPoint (вызов другого домена)
// 	if err := p.rentPointUC.AddProduct(ctx, pointUUID, prodUUID); err != nil {
// 		 // Откатываем, если не удалось добавить в RentPoint
// 		 if rollbackErr := p.repo.SetRentPoint(ctx, prodUUID, uuid.Nil); rollbackErr != nil {
// 			  return entity.Product{}, fmt.Errorf("failed to add to rent point: %v, rollback failed: %w", err, rollbackErr)
// 		 }
// 		 return entity.Product{}, fmt.Errorf("failed to add to rent point: %w", err)
// 	}
//
// 	// 4. Возвращаем обновленный продукт
// 	return p.repo.GetByID(ctx, prodUUID)
// }
