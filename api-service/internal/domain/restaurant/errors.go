package restaurant

import "errors"

var (
	ErrPizzaSizeAlreadyExists = errors.New("pizza-size already exists")
	ErrEmailAlreadyExists     = errors.New("email already exists")
)
