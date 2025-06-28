package nodb

import (
	"EasyRentGo/internal/entity"
	"EasyRentGo/internal/repo"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type template struct {
	id          uuid.UUID
	name        string
	description string
	price       int
}

var _ repo.ProductTempRepo = (*DBTemplate)(nil)

var templates []template

type DBTemplate struct {
	DB *[]template
}

func NewTemplate() *DBTemplate {
	// templates = make([]template, 0)
	return &DBTemplate{
		DB: &templates,
	}
}

func (db *DBTemplate) Create(ctx context.Context, t entity.ProductTemp) error {
	template := template{
		id:          t.ID,
		name:        t.Name,
		description: t.Description,
		price:       t.Price,
	}

	*db.DB = append(*db.DB, template)
	return nil
}

func (db *DBTemplate) GetAll(ctx context.Context) ([]entity.ProductTemp, error) {
	templates := make([]entity.ProductTemp, 0, len(*db.DB))
	for _, t := range *db.DB {
		templates = append(templates, entity.ProductTemp{
			ID:          t.id,
			Name:        t.name,
			Description: t.description,
			Price:       t.price,
		})
	}
	return templates, nil
}

func (db *DBTemplate) GetByID(ctx context.Context, id uuid.UUID) (entity.ProductTemp, error) {
	tmp, err := db.getObject(id)
	if err != nil {
		return entity.ProductTemp{}, fmt.Errorf("error: %s", err)
	}

	return entity.ProductTemp{
		ID:          tmp.id,
		Name:        tmp.name,
		Description: tmp.description,
		Price:       tmp.price,
	}, nil
}

// Получение объекта по id
func (db *DBTemplate) getObject(id uuid.UUID) (template, error) {
	for _, tmp := range *db.DB {
		if tmp.id == id {
			return tmp, nil
		}
	}

	return template{}, fmt.Errorf("object Product: %s - not found", id)
}
