package repositories

import (
	"Fridger/internal/domain/models"
	"context"
)

type ProductRepo interface {
	Add(ctx context.Context, product *models.Product) error
	FindByCis(ctx context.Context, cis string) (*models.Product, error)
	DeleteByCis(ctx context.Context, cis string) error
}
