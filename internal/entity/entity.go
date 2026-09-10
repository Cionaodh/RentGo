package entity

import (
	"time"

	"github.com/google/uuid"
)

// User — базовая сущность пользователя
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// UserProfile - включает данные юзера и его заказы
type UserProfile struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Orders   []Order   `json:"orders"`
}

type RentPoint struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"`
	Addr string    `json:"addr" db:"addr"`
	// TODO: Координаты (45.63545 74.54345)
}

type ProductRentPoint struct {
	ID       uuid.UUID    `json:"id"`
	Name     string       `json:"name"`
	Addr     string       `json:"addr"`
	Products []ProductsRP `json:"products"`
}

type ProductTemplate struct {
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

func (s ProductStatus) IsValid() bool {
	switch s {
	case StatusUnused,
		StatusFree,
		StatusReserved,
		StatusRented:
		return true
	default:
		return false
	}
}

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

type ProductsRP struct {
	Id     uuid.UUID     `json:"id"`
	Name   string        `json:"name"`
	Status ProductStatus `json:"status"`
	Price  float64       `json:"price"`
}

type Product struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Price       float64       `json:"price"`
	Status      ProductStatus `json:"status"`
	RentPointID uuid.UUID     `json:"rentpoint_id"`
}

type OrderStatus string

const (
	StatusActive    OrderStatus = "Active"    // Активный
	StatusCompleted OrderStatus = "Completed" // Завершенный
	StatusExpired   OrderStatus = "Expired"   // Просрочен - если пользователь не завершил поездку за купленное время
)

type Order struct {
	ID            uuid.UUID   `json:"id"`
	UserID        uuid.UUID   `json:"user_id"`
	Status        OrderStatus `json:"status"`
	ProductID     uuid.UUID   `json:"product_id"`
	StartPointID  uuid.UUID   `json:"start_point_id"`
	FigishPointID *uuid.UUID  `json:"finish_point_id,omitempty"`
	StartedAT     time.Time   `json:"start_at"`
	FinishedAT    *time.Time  `json:"finish_at,omitempty"`
}

// // ParseProductStatus конвертирует строку в ProductStatus
// func ParseProductStatus(s string) (ProductStatus, error) {
// 	switch s {
// 	case string(StatusUnused), string(StatusFree),
// 		string(StatusReserved), string(StatusRented):
// 		return ProductStatus(s), nil
// 	default:
// 		return "", fmt.Errorf("invalid ProductStatus: %s", s)
// 	}
// }

// // MustParseProductStatus паникует при невалидном статусе
// func MustParseProductStatus(s string) ProductStatus {
// 	status, err := ParseProductStatus(s)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return status
// }
