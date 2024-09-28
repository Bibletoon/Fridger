package services

import (
	"Fridger/internal/domain/interfaces/clients"
	"Fridger/internal/domain/interfaces/repositories"
	"Fridger/internal/domain/interfaces/services"
	"Fridger/internal/domain/models"
	"Fridger/internal/helpers"
	"context"
	"time"
)

type productService struct {
	productRepo repositories.ProductRepo
	crptClient  clients.CrptClient
}

func NewProductService(productRepo repositories.ProductRepo, crptClient clients.CrptClient) services.ProductService {
	return &productService{productRepo, crptClient}
}

func (s *productService) GetAllActiveProducts(ctx context.Context) ([]*models.Product, error) {
	products, err := s.productRepo.GetAllActive(ctx)

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *productService) GetExpiringProducts(ctx context.Context, daysBeforeExpiration int) (expiringProducts, expiredProducts []*models.Product, err error) {
	t := time.Now().AddDate(0, 0, -daysBeforeExpiration)
	products, err := s.productRepo.GetExpiringBefore(ctx, t)

	if err != nil {
		return nil, nil, err
	}

	expiringProducts, expiredProducts = splitProducts(products, t)

	return expiringProducts, expiredProducts, nil
}

func splitProducts(products []*models.Product, t time.Time) (expiringProducts, expiredProducts []*models.Product) {
	expiringProducts = make([]*models.Product, 0)
	expiredProducts = make([]*models.Product, 0)

	for _, product := range products {
		if product.ExpirationDate.Before(t) {
			expiredProducts = append(expiredProducts, product)
		} else {
			expiringProducts = append(expiringProducts, product)
		}
	}

	return expiringProducts, expiredProducts
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

func (s *productService) GetProductByDatamatrix(ctx context.Context, datamatrix string) (*models.Product, error) {
	cis := helpers.ParseCis(datamatrix)
	product, err := s.productRepo.GetBySerial(ctx, cis)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) DeleteProductByDatamatrix(ctx context.Context, datamatrix string) error {
	cis := helpers.ParseCis(datamatrix)
	err := s.productRepo.DeleteBySerial(ctx, cis)
	if err != nil {
		return err
	}

	return nil
}
