package domain

import "context"

type CatalogProduct struct {
	ID          string
	StoreID     string
	Name        string
	Description string
	SKU         string
	Price       float64
}

type CatalogRepository interface {
	AddProduct(ctx context.Context, id string, storeID string, name string, description string, price float64) error
}
