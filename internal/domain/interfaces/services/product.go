package services

import (
	"Fridger/internal/domain/models"
	"context"
)

type ProductService interface {
	AddProductByDatamatix(ctx context.Context, datamatrix string) (*models.Product, error)
}
