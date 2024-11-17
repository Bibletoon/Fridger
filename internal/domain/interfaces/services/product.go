package services

import (
	"Fridger/internal/domain/models"
	"context"
)

type ProductService interface {
	GetAllActiveProducts(ctx context.Context) ([]*models.Product, error)
	GetExpiringProducts(ctx context.Context, daysBeforeExpiration int) (*models.ProductsCollection, error)
	GetProductByDatamatrix(ctx context.Context, datamatrix string) (*models.Product, error)
	AddProductByDatamatix(ctx context.Context, datamatrix string) (*models.Product, error)
	DeleteProductByDatamatrix(ctx context.Context, datamatrix string) error
}
