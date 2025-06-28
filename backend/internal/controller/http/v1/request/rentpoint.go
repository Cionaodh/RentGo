package request

type RentPoint struct {
	Name string `json:"name"       validate:"required"  example:"RentPoint 1"`
	Addr string `json:"addr"       validate:"required"  example:"г. Калининград, ул Баласа"`
}

type TmpProduct struct {
	Name        string `json:"name" validate:"required" example:"Велосипед 1"`
	Description string `json:"descriptiont" example:"Описание продукта"`
	Price       int    `json:"price" validate:"required" example:"800"`
}
