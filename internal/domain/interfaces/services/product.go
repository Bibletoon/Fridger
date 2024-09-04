package services

import (
	"Fridger/internal/domain/models"
	"context"
)

type ProductService interface {
	GetAllActiveProducts(ctx context.Context) ([]*models.Product, error)
	GetProductByDatamatrix(ctx context.Context, datamatrix string) (*models.Product, error)
	AddProductByDatamatix(ctx context.Context, datamatrix string) (*models.Product, error)
	DeleteProductByDatamatrix(ctx context.Context, datamatrix string) error
}
