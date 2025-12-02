package usecase

import "errors"

var (
	ErrRentPointAlreadyExists = errors.New("название или адрес точки проката уже существует в системе")
	ErrFieldIsTooLong         = errors.New("значение полей имени и/или адреса превышают допустимую длину")
	ErrCreateRentpoint        = errors.New("возникла ошибка при создании точки проката")
	ErrNumberProduct          = errors.New("превышено число создаваемых объектов")

	ErrCreateProduct = errors.New("возникла ошибка при создании продукта")
)
