package services

import (
	"Fridger/internal/domain/interfaces/clients"
	"Fridger/internal/domain/interfaces/repositories"
	"Fridger/internal/domain/interfaces/services"
	"Fridger/internal/domain/models"
	"Fridger/internal/helpers"
	"context"
)

type productService struct {
	productRepo repositories.ProductRepo
	crptClient  clients.CrptClient
}

func NewProductService(productRepo repositories.ProductRepo, crptClient clients.CrptClient) services.ProductService {
	return &productService{productRepo, crptClient}
}

func (s *productService) AddProductByDatamatix(ctx context.Context, datamatrix string) (*models.Product, error) {
	product, err := s.crptClient.GetByDatamatrix(ctx, datamatrix)
	if err != nil {
		return nil, err
	}

	err = s.productRepo.Add(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) FindProductByDatamatrix(ctx context.Context, datamatrix string) (*models.Product, error) {
	cis := helpers.ParseCis(datamatrix)
	product, err := s.productRepo.FindByCis(ctx, cis)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) DeleteProductByDatamatrix(ctx context.Context, datamatrix string) error {
	cis := helpers.ParseCis(datamatrix)
	err := s.productRepo.DeleteByCis(ctx, cis)
	if err != nil {
		return err
	}

	return nil
}
