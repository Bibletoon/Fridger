package models

import "time"

type ProductsCollection struct {
	products []*Product
}

func NewProductsCollection(products []*Product) *ProductsCollection {
	return &ProductsCollection{
		products: products,
	}
}

func (pc *ProductsCollection) All() []*Product {
	return pc.products
}

func (pc *ProductsCollection) Len() int {
	return len(pc.products)
}

func (pc *ProductsCollection) GetExpired() []*Product {
	t := time.Now()
	expiredProducts := make([]*Product, 0)

	for _, p := range pc.products {
		if p.ExpirationDate.Before(t) {
			expiredProducts = append(expiredProducts, p)
		}
	}

	return expiredProducts
}

func (pc *ProductsCollection) GetNotExpired() []*Product {
	t := time.Now()
	expiredProducts := make([]*Product, 0)

	for _, p := range pc.products {
		if p.ExpirationDate.After(t) {
			expiredProducts = append(expiredProducts, p)
		}
	}

	return expiredProducts
}
