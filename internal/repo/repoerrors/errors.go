package repoerrors

import "errors"

var (
	ErrAlreadyExists       = errors.New("object already exists")
	ErrNotFound            = errors.New("not found")
	ErrForeignKeyViolation = errors.New("foreign key violation")

	ErrProductNotAvailable = errors.New("product not available")
	ErrProductStateInvalid = errors.New("product state invalid")
	ErrOrderNotActive      = errors.New("order not active")

	ErrRentPointNotFound = errors.New("rentpoint not found")
)
