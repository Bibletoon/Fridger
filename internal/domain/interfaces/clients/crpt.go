package clients

import (
	"Fridger/internal/infrastructure/clients/dto"
	"context"
)

type CrptClient interface {
	GetByDatamatrix(ctx context.Context, datamatrix string) (*dto.ProductInfoDto, error)
}
