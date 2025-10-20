package entity

import (
	"github.com/google/uuid"
)

type RentPoint struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
	Addr string    `db:"addr"`
	// TODO: Координаты (45.63545 74.54345)
}

type ProductRentPoint struct {
	ID       uuid.UUID `db:"id"`
	Name     string    `db:"name"`
	Addr     string    `db:"addr"`
	Products []Product `db:"product"`
	// TODO: Координаты (45.63545 74.54345)
}

type ProductTemp struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       int       `db:"price"`
}

type ProductStatus string

const (
	StatusUnused   ProductStatus = "unused"
	StatusFree     ProductStatus = "free"
	StatusReserved ProductStatus = "reserved"
	StatusRented   ProductStatus = "rented"
)

type Product struct {
	ID          uuid.UUID
	TemplateID  uuid.UUID
	RentPointID uuid.UUID
	Status      ProductStatus
}

// // ParseProductStatus конвертирует строку в ProductStatus (с проверкой валидности).
// func ParseProductStatus(s string) (ProductStatus, error) {
// 	switch s {
// 	case string(StatusUnused), string(StatusFree),
// 		string(StatusReserved), string(StatusRented):
// 		return ProductStatus(s), nil
// 	default:
// 		return "", fmt.Errorf("invalid ProductStatus: %s", s)
// 	}
// }

// // MustParseProductStatus паникует при невалидном статусе (использовать только в тестах/инициализации).
// func MustParseProductStatus(s string) ProductStatus {
// 	status, err := ParseProductStatus(s)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return status
// }
