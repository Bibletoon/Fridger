package repositories

import (
	"Fridger/internal/domain/models"
	"context"
)

type ProductRepo interface {
	Add(ctx context.Context, product *models.Product) error
	GetBySerial(ctx context.Context, serial string) (*models.Product, error)
	DeleteBySerial(ctx context.Context, serial string) error
}
