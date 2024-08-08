package services

import (
	"Fridger/internal/domain/models"
	"context"
)

type ProductService interface {
	GetProductByDatamatrix(ctx context.Context, datamatrix string) (*models.Product, error)
	AddProductByDatamatix(ctx context.Context, datamatrix string) (*models.Product, error)
	DeleteProductByDatamatrix(ctx context.Context, datamatrix string) error
}
