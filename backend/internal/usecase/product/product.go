package product

import (
	"EasyRentGo/internal/repo"
)

type Product struct {
	repo repo.ProductRepo // указание на репу (например, repo.ProductRepo)
	// указание на внешние интерфейсы (например, ProductWebAPI)
}

func New(p repo.ProductRepo) *Product { // передаем
	return &Product{
		repo: p,
	}
}
