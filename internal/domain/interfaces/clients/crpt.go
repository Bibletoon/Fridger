package clients

import (
	"Fridger/internal/domain/models"
	"context"
)

type CrptClient interface {
	GetByDatamatrix(ctx context.Context, datamatrix string) (*models.Product, error)
}
