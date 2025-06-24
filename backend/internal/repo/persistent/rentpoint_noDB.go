package persistent

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo"
	"context"
	"fmt"

	"github.com/google/uuid"
)

var _ repo.RentPointRepo = (*DBRentPoint)(nil)

type Point struct {
	ID       uuid.UUID
	Name     string
	Addr     string
	Products []*entity.Product
	// Координаты (45.63545 74.54345 - чтобы указать точную точку)
}

var RentPoints []Point

type DBRentPoint struct {
	DB *[]Point
}

func New() *DBRentPoint {
	DefaultPoint := Point{
		Name:     "Аренда",
		Products: make([]*entity.Product, 0),
	}
	RentPoints = append(RentPoints, DefaultPoint)

	return &DBRentPoint{
		DB: &RentPoints,
	}
}

// Создаем нулевую точку проката для хранения арендованных продуктов

// Создание новой точки проката
func (db *DBRentPoint) Create(ctx context.Context, rp entity.RentPoint) error {
	// Проверяем, есть ли уже точка с таким ID
	for _, existing := range RentPoints {
		if existing.ID == rp.ID {
			return fmt.Errorf("точка проката с ID %s уже существует", rp.ID)
		}
	}

	p := Point{
		ID:       rp.ID,
		Name:     rp.Name,
		Addr:     rp.Addr,
		Products: make([]*entity.Product, 0),
	}

	RentPoints = append(RentPoints, p)
	return nil
}

// Получение всех точек проката (не передаем прикрепленные продукты)
func (db *DBRentPoint) GetAll(ctx context.Context) ([]entity.RentPoint, error) {
	rentPointArr := make([]entity.RentPoint, 0)
	for i := 0; i < len(RentPoints); i++ {
		rentPointArr = append(rentPointArr, entity.RentPoint{
			ID:   RentPoints[i].ID,
			Name: RentPoints[i].Name,
			Addr: RentPoints[i].Addr,
		})
	}
	return rentPointArr, nil
}

// Получение точки проката по id (со списком продуктов)
func (db *DBRentPoint) GetByID(ctx context.Context, id uuid.UUID) (entity.RentPoint, error) {
	for _, rp := range RentPoints {
		if rp.ID == id {
			// Возвращаем полную копию объекта
			return entity.RentPoint{
				ID:   rp.ID,
				Name: rp.Name,
				Addr: rp.Addr,
			}, nil
		}
	}
	return entity.RentPoint{}, fmt.Errorf("точка проката с ID %s не найдена", id)
}
