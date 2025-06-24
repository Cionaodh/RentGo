package repo

import (
	"EasyRentGo/internal/entity"
	"context"

	"github.com/google/uuid"
)

type RentPointRepo interface {
	Create(context.Context, entity.RentPoint) error               // Создание новой точки проката
	GetAll(context.Context) ([]entity.RentPoint, error)           // Получение всех точек проката (не передаем прикрепленные продукты)
	GetByID(context.Context, uuid.UUID) (entity.RentPoint, error) // Получение точки проката по id (со списком продуктов)
	// Удаление точки проката

	// Прикрепление продукта к точке проката
}
