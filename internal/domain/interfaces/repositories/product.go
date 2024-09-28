package repositories

import (
	"Fridger/internal/domain/models"
	"context"
	"time"
)

type ProductRepo interface {
	Add(ctx context.Context, product *models.Product) error
	GetExpiringBefore(ctx context.Context, time time.Time) ([]*models.Product, error)
	GetAllActive(ctx context.Context) ([]*models.Product, error)
	GetBySerial(ctx context.Context, serial string) (*models.Product, error)
	DeleteBySerial(ctx context.Context, serial string) error
}
