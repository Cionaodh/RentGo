package nodb

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo"
	"context"
	"fmt"

	"github.com/google/uuid"
)

var _ repo.RentPointRepo = (*DBRentPoint)(nil)

type point struct {
	id       uuid.UUID
	name     string
	addr     string
	products []*product
	// Координаты (45.63545 74.54345 - чтобы указать точную точку)
}

var rentPoints []point

type DBRentPoint struct {
	DB *[]point
}

func NewRentPoint() *DBRentPoint {
	DefaultPoint := point{
		name:     "Аренда",
		products: make([]*product, 0),
	}
	rentPoints = append(rentPoints, DefaultPoint)

	products = append(products, product{})
	return &DBRentPoint{
		DB: &rentPoints,
	}
}

// Создаем нулевую точку проката для хранения арендованных продуктов

// Создание новой точки проката
func (db *DBRentPoint) Create(ctx context.Context, rp entity.RentPoint) error {
	// Проверяем, есть ли уже точка с таким ID
	for _, existing := range rentPoints {
		if existing.id == rp.ID {
			return fmt.Errorf("точка проката с ID %s уже существует", rp.ID)
		}
	}

	p := point{
		id:       rp.ID,
		name:     rp.Name,
		addr:     rp.Addr,
		products: make([]*product, 0),
	}

	rentPoints = append(rentPoints, p)
	return nil
}

// Получение всех точек проката (не передаем прикрепленные продукты)
func (db *DBRentPoint) GetAll(ctx context.Context) ([]entity.RentPoint, error) {
	rentPointArr := make([]entity.RentPoint, 0)
	for i := 0; i < len(rentPoints); i++ {
		rentPointArr = append(rentPointArr, entity.RentPoint{
			ID:   rentPoints[i].id,
			Name: rentPoints[i].name,
			Addr: rentPoints[i].addr,
		})
	}
	return rentPointArr, nil
}

// Получение точки проката по id (со списком продуктов)
func (db *DBRentPoint) GetByID(ctx context.Context, id uuid.UUID) (entity.RentPoint, error) {

	rp, err := db.getObject(id)
	if err != nil {
		return entity.RentPoint{}, err
	}

	return entity.RentPoint{
		ID:   rp.id,
		Name: rp.name,
		Addr: rp.addr,
		// Products: rp.Products, // TODO: добавить
	}, nil

	// return entity.RentPoint{}, fmt.Errorf("точка проката с ID %s не найдена", id)
}

// Получение объекта по id
func (db *DBRentPoint) getObject(id uuid.UUID) (point, error) {
	for _, rp := range rentPoints {
		if rp.id == id {
			return rp, nil
		}
	}
	return point{}, fmt.Errorf("object RentPoint: %s - not found", id)
}

//
