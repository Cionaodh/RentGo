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
	ID       uuid.UUID  `db:"id"`
	Name     string     `db:"name"`
	Addr     string     `db:"addr"`
	Products []Products `db:"product"`
	// TODO: Координаты (45.63545 74.54345)
}

type ProductTemp struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       float64   `db:"price"`
}

type ProductStatus string

const (
	StatusUnused   ProductStatus = "Unused"
	StatusFree     ProductStatus = "Free"
	StatusReserved ProductStatus = "Reserved"
	StatusRented   ProductStatus = "Rented"
)

type Products struct {
	TemplateID  uuid.UUID     `db:"template_id"`
	RentPointID uuid.UUID     `db:"rentpoint_id"`
	Status      ProductStatus `db:"status"`
	IDs         uuid.UUIDs    `db:"ids"`
}

type Product struct {
	ID          uuid.UUID     `db:"id"`
	TemplateID  uuid.UUID     `db:"template_id"`
	RentPointID uuid.UUID     `db:"rentpoint_id"`
	Status      ProductStatus `db:"status"`
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
