package clients

import (
	"Fridger/internal/models"
	"context"
)

type CrptClient interface {
	GetByDatamatrix(ctx *context.Context, datamatrix string) (*models.ProductInfo, error)
}