package product

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo"
	"EasyRentGo/internal/usecase"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Product struct {
	repo        repo.ProductRepo // указание на репу (например, repo.ProductRepo)
	// rentPointUC usecase.RentPointUseCase
}

func New(p repo.ProductRepo, rpUC usecase.RentPointUseCase) *Product {
	return &Product{
		repo:        p,
		rentPointUC: rpUC,
	}
}

// // SetRentPointUC: установка связи с RentPointUseCase после создания объекта
// func (p *Product) SetRentPointUC(rpUC usecase.RentPointUseCase) {
// 	p.rentPointUC = rpUC
// }

func (p *Product) Create(ctx context.Context, idTmp uuid.UUID) (entity.Product, error) {
	product := entity.Product{
		ID:          uuid.New(),
		TemplateID:  idTmp,
		RentPointID: uuid.Nil, // дефолтное значение - на складе
		Status:      entity.StatusUnused,
	}
	if err := p.repo.Create(ctx, product); err != nil {
		return entity.Product{}, fmt.Errorf("%s", err)
	}
	return product, nil
}

func (p *Product) GetAll(ctx context.Context) ([]entity.Product, error) {
	return p.repo.GetAll(ctx)
}

func (p *Product) GetByDI(ctx context.Context, id string) (entity.Product, error) {
	// проверка uuid на корректность и преобразование
	if err := uuid.Validate(id); err != nil {
		return entity.Product{}, fmt.Errorf("%s", err)
	}

	return p.repo.GetByID(ctx, uuid.MustParse(id))
}

func (p *Product) GetByStatus(ctx context.Context, status string) ([]entity.Product, error) {
	// Проверка корректности статуса
	s, err := entity.ParseProductStatus(status)
	if err != nil {
		return []entity.Product{}, fmt.Errorf("%s", err)
	}

	return p.repo.GetByStatus(ctx, s)
}

func (p *Product) SetStatus(ctx context.Context, id string, status string) (entity.Product, error) {
	// Проверка корректности статуса
	s, err := entity.ParseProductStatus(status)
	if err != nil {
		return entity.Product{}, fmt.Errorf("%s", err)
	}

	// проверка uuid на корректность и преобразование
	if err := uuid.Validate(id); err != nil {
		return entity.Product{}, fmt.Errorf("%s", err)
	}

	if err := p.repo.SetStatus(ctx, uuid.MustParse(id), s); err != nil {
		return entity.Product{}, fmt.Errorf("%s", err)
	}

	prod, err := p.repo.GetByID(ctx, uuid.MustParse(id))
	if err != nil {
		return entity.Product{}, fmt.Errorf("%s", err)
	}

	return prod, nil
}

func (p *Product) SetRentPoint(ctx context.Context, idProduct string, idPoint string) (entity.Product, error) {
	if err := uuid.Validate(idProduct); err != nil {
		return entity.Product{}, fmt.Errorf("%s", err)
	}

	if err := uuid.Validate(idPoint); err != nil {
		return entity.Product{}, fmt.Errorf("%s", err)
	}

	point, err := p.repo.SetRentPoint(ctx, uuid.MustParse(idProduct), uuid.MustParse(idPoint)) // Прикрепляем к продукту - пункт проката
	if err != nil {
		return entity.Product{}, fmt.Errorf("%s", err)
	}

	p.rentPointUC.

	// p.rentPointUC.

	// Получаем объект точки проката!!!!
	// Вызываем фукнцию добавления продукта в хранилще продуктов Пункта проката

	// TODO: как полчить доступ к домену Пункта проката в бизнес слое?
	return entity.Product{}, nil
}

// // Product UseCase
// func (p *Product) SetRentPoint(ctx context.Context, productID, pointID string) (entity.Product, error) {
// 	// 1. Валидация ID
// 	prodUUID, err := uuid.Parse(productID)
// 	if err != nil { /* ... */ }

// 	pointUUID, err := uuid.Parse(pointID)
// 	if err != nil { /* ... */ }

// 	// 2. Обновляем RentPointID у продукта
// 	if err := p.repo.SetRentPoint(ctx, prodUUID, pointUUID); err != nil {
// 		 return entity.Product{}, fmt.Errorf("failed to update product: %w", err)
// 	}

// 	// 3. Добавляем продукт в RentPoint (вызов другого домена)
// 	if err := p.rentPointUC.AddProduct(ctx, pointUUID, prodUUID); err != nil {
// 		 // Откатываем, если не удалось добавить в RentPoint
// 		 if rollbackErr := p.repo.SetRentPoint(ctx, prodUUID, uuid.Nil); rollbackErr != nil {
// 			  return entity.Product{}, fmt.Errorf("failed to add to rent point: %v, rollback failed: %w", err, rollbackErr)
// 		 }
// 		 return entity.Product{}, fmt.Errorf("failed to add to rent point: %w", err)
// 	}

// 	// 4. Возвращаем обновленный продукт
// 	return p.repo.GetByID(ctx, prodUUID)
// }

// // RentPoint UseCase
// func (rp *RentPoint) AddProduct(ctx context.Context, pointID, productID uuid.UUID) error {
// 	point, err := rp.repo.GetByID(ctx, pointID)
// 	if err != nil { /* ... */ }

// 	point.Products = append(point.Products, productID)

// 	return rp.repo.Update(ctx, point)
// }
