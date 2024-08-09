package errors

import "errors"

var (
	ErrNotFound              = errors.New("not found")
	ErrProductNotFoundInCrpt = errors.New("product not found")
	ErrCodeRead              = errors.New("code read error")
	ErrProductExists         = errors.New("product already exists")
)
