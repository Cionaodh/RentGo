package usecase

import "errors"

var (
	ErrRentPointAlreadyExists = errors.New("Название или адрес точки проката уже существует в системе")
	ErrFieldIsTooLong         = errors.New("Значение полей имени и/или адреса превышают допустимую длину")
	ErrCreateRentpoint        = errors.New("Возникла ошибка при создании точки проката")
)
