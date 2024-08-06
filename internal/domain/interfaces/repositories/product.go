package repositories

import (
	"Fridger/internal/domain/models"
	"context"
)

type ProductRepo interface {
	Save(ctx context.Context, product *models.Product) error
}
