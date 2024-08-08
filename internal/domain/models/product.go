package models

import "time"

type Product struct {
	Name           string
	Gtin           int64
	Serial         string
	Category       string
	ExpirationDate time.Time
	IsActive       bool
	CreatedAt      time.Time
	DeletedAt      time.Time
}
