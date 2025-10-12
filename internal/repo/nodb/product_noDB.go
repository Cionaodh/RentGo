package nodb

import (
	"EasyRentGo/internal/entity"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type product struct {
	id          uuid.UUID
	templateID  *template
	rentPointID *point
	status      entity.ProductStatus
}

var products []product

type DBProduct struct {
	DB *[]product
}

func NewProduct() *DBProduct { // передаем
	return &DBProduct{
		DB: &products,
	}
}

func (p *DBProduct) Create(ctx context.Context, id uuid.UUID, num int) error {

	return nil
} // Создание продукта на основе шаблона (id шаблона, кол-во продуктов)

func (p *DBProduct) GetAll(ctx context.Context) ([]entity.Product, error) {
	return []entity.Product{}, nil
} // получение всех продуктов

func (p *DBProduct) GetByID(ctx context.Context, id uuid.UUID) (entity.Product, error) {
	return entity.Product{}, nil
} // Получение продукта по id

func (p *DBProduct) GetByStatus(ctx context.Context, status entity.ProductStatus) ([]entity.Product, error) {
	return []entity.Product{}, nil
} // Получение продуктов с определенным статусом

func (p *DBProduct) SetStatus(ctx context.Context, id uuid.UUID, status entity.ProductStatus) (entity.Product, error) {
	return entity.Product{}, nil
}

// Получение объекта по id
func (db *DBProduct) getObject(id uuid.UUID) (product, error) {
	for _, p := range products {
		if p.id == id {
			return p, nil
		}
	}

	return product{}, fmt.Errorf("object Product: %s - not found", id)
}
