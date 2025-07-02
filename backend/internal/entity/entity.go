package entity

import (
	"fmt"

	"github.com/google/uuid"
)

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
	StatusUnused   ProductStatus = "unused"   // "Неисползуется"
	StatusFree     ProductStatus = "free"     // "Свободен"
	StatusReserved ProductStatus = "reserved" // "Зарезервирован"
	StatusRented   ProductStatus = "rented"   // "Арендован"
)

type Product struct {
	ID          uuid.UUID
	TemplateID  uuid.UUID
	RentPointID uuid.UUID
	Status      ProductStatus
}

// ParseProductStatus конвертирует строку в ProductStatus (с проверкой валидности).
func ParseProductStatus(s string) (ProductStatus, error) {
	switch s {
	case string(StatusUnused), string(StatusFree),
		string(StatusReserved), string(StatusRented):
		return ProductStatus(s), nil
	default:
		return "", fmt.Errorf("invalid ProductStatus: %s", s)
	}
}

// // MustParseProductStatus паникует при невалидном статусе (использовать только в тестах/инициализации).
// func MustParseProductStatus(s string) ProductStatus {
// 	status, err := ParseProductStatus(s)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return status
// }
