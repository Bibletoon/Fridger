package errors

import "errors"

var (
	ErrProductNotFound = errors.New("product not found")
	ErrResponseRead    = errors.New("response read error")
	ErrProductExists   = errors.New("product already exists")
)
