package product

type ProductUseCase struct {
	// указание на репу (например, repo.ProductRepo)
	// указание на внешние интерфейсы (например, ProductWebAPI)
}

func New() *ProductUseCase { // передаем 
	return &ProductUseCase{}
}

// Реализуем интерфейс 
