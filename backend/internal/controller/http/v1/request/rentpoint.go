package request

type RentPoint struct {
	Name string `json:"name"       validate:"required"  example:"RentPoint 1"`
	Addr string `json:"addr"       validate:"required"  example:"г. Калининград, ул Баласа"`
}
