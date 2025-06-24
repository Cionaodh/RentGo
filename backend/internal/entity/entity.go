package entity

import "github.com/google/uuid"

// Точка проката
type RentPoint struct {
	ID       uuid.UUID
	Name     string
	Addr     string
	Products []uuid.UUID
	// Координаты (45.63545 74.54345 - чтобы указать точную точку)
}

// Шаблон товара
type ProductTemp struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       int
}

// Статус продукта
type ProductStatus string

const (
	StatusFree     ProductStatus = "free"
	StatusReserved ProductStatus = "reserved"
	StatusRented   ProductStatus = "rented"
)

// Продукт
type Product struct {
	ID          uuid.UUID
	TemplateID  uuid.UUID
	RentPointID uuid.UUID
	Status      ProductStatus
}
