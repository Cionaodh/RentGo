package repoerrors

import "errors"

var (
	ErrAlreadyExists = errors.New("object already exists")
	ErrNotFound      = errors.New("not found")
)
