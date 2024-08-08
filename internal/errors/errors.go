package errors

import "errors"

var (
	ErrProductNotFoundInCrpt = errors.New("product not found")
	ErrCodeRead              = errors.New("code read error")
	ErrProductExists         = errors.New("product already exists")
)
