package domain

import "context"

type ProductRepository interface {
	AddProduct(
		ctx context.Context,
		id string,
		storeID string,
		name string,
		description string,
		sku string,
		price float64,
	) error
}
