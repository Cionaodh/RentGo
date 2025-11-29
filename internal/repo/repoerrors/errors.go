package repoerrors

import "errors"

var (
	ErrRentPointAlreadyExists = errors.New("rentpoint already exists")
	ErrNotFound               = errors.New("not found")
)
