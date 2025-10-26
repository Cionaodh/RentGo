package entity

import (
	"github.com/google/uuid"
)

type RentPoint struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Addr string    `json:"addr"`
	// TODO: Координаты (45.63545 74.54345)
}

type ProductRentPoint struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Addr     string    `json:"addr"`
	Products []struct {
		ID     uuid.UUID     `json:"id"`
		Name   string        `json:"name"`
		Price  float64       `json:"price"`
		Status ProductStatus `json:"status"`
	} `json:"products"`
}

type ProductTemp struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"desc"`
	Price       float64   `json:"price"`
}

type ProductStatus string

const (
	StatusUnused   ProductStatus = "Unused"
	StatusFree     ProductStatus = "Free"
	StatusReserved ProductStatus = "Reserved"
	StatusRented   ProductStatus = "Rented"
)

// type ProductsTemp struct {
// 	RentPointID uuid.UUID
// 	Name        string
// 	Status      ProductStatus
// 	Price       float64
// 	IDs         uuid.UUIDs
// }

type Products struct {
	RentPointID uuid.UUID     `json:"rentpoint_id"`
	Status      ProductStatus `json:"status"`
	IDs         uuid.UUIDs    `json:"ids"`
	TemplateID  uuid.UUID     `json:"template_id"`
}

type Product struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Price       float64       `json:"price"`
	Status      ProductStatus `json:"status"`
	RentPointID uuid.UUID     `json:"rentpoint_id"`
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
