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

	// AddPriduct(context.Context, uuid.UUID) error// Прикрепление продукта к точке проката
}

type ProductTempRepo interface {
	Create(context.Context, entity.ProductTemp) error
	GetAll(context.Context) ([]entity.ProductTemp, error)
	GetByID(context.Context, uuid.UUID) (entity.ProductTemp, error)
}

// (Шаблон, кол-во (стандартно = 1), (необязательный параметр = прикрепление к точке проката, стандартно nil), стандартный статус unused)
type ProductRepo interface {
	// id шаблона,
	Create(context.Context, uuid.UUID, int) ([]entity.Product, error)                   // Создание продукта на основе шаблона (id шаблона, кол-во продуктов)
	GetAll(context.Context) ([]entity.Product, error)                                   // получение всех продуктов
	GetByID(context.Context, uuid.UUID) (entity.Product, error)                         // Получение продукта по id
	GetByStatus(context.Context, entity.ProductStatus) ([]entity.Product, error)        // Получение продуктов с определенным статусом
	SetStatus(context.Context, uuid.UUID, entity.ProductStatus) (entity.Product, error) // Изменение статуса продукта
	// -- Прикрепление продукта к точке проката (id точки проката, id продукта)
	// Мы должны:
	// 1)поменять указатель на точку,
	// 2) поменять статус продукта
	// 3) Добавить объекту RentPoint указатель на продукт (вызвать функцию добавления_продукта(id продукта) )
	// -- Прикрепление продуктов к точке проката

	// -- Изменить точку проката для продутка
	// --
}

// Пользователь сдает арендованный продкут - сотрудник проката сканирует QR код на продукте -> переходит в карточку продукта -> заканчивает аренду + прикрепляет продукт к своей точке проката
// Функция закончить_аренду(продукт, точка_проката) -> ф-я меняет статус продукта и прикрепляет его к точки проката
