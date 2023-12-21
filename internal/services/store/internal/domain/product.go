package domain

import (
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/es"
	"github.com/stackus/errors"
)

const ProductAggregate = "stores.Product"

var (
	ErrNameIsBlank            = errors.Wrap(errors.ErrBadRequest, "the product name cannot be blank")
	ErrProductPriceIsNegative = errors.Wrap(errors.ErrBadRequest, "the product price cannot be negative")
	ErrNotAPriceIncrease      = errors.Wrap(errors.ErrBadRequest, "the price change would be a decrease")
	ErrNotAPriceDecrease      = errors.Wrap(errors.ErrBadRequest, "the price change would be a increase")
)

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Product)(nil)

type Product struct {
	es.Aggregate
	storeID     string
	name        string
	description string
	sku         string
	price       float64
}

var _ es.EventSourcedAggregate = (*Product)(nil)

func NewProduct(id string) *Product {
	return &Product{
		Aggregate: es.NewAggregate(id, ProductAggregate),
	}
}

func (p *Product) Init(
	id string,
	storeID string,
	name string,
	description string,
	sku string,
	price float64,
) (ddd.Event, error) {
	p.AggregateName()
	if name == "" {
		return nil, ErrNameIsBlank
	}

	if price < 0 {
		return nil, ErrProductPriceIsNegative
	}

	p.AddEvent(ProductAddEvent, &ProductAdded{
		StoreID:     storeID,
		Name:        name,
		Description: description,
		SKU:         sku,
		Price:       price,
	})

	return ddd.NewEvent(ProductAddEvent, p), nil
}

func (p Product) StoreID() string {
	return p.storeID
}

func (p Product) Name() string {
	return p.name
}

func (p Product) Description() string {
	return p.description
}

func (p Product) SKU() string {
	return p.sku
}

func (p Product) Price() float64 {
	return p.price
}

// Key implements registry.Registerable
func (Product) Key() string {
	return ProductAggregate
}

func (p *Product) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *ProductAdded:
		p.storeID = payload.StoreID
		p.name = payload.Name
		p.description = payload.Description
		p.sku = payload.SKU
		p.price = payload.Price

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", p, event.EventName(), payload)
	}

	return nil
}

func (p *Product) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *ProductV1:
		p.storeID = ss.StoreID
		p.name = ss.Name
		p.description = ss.Description
		p.sku = ss.SKU
		p.price = ss.Price
	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", p, snapshot)
	}

	return nil
}

func (p Product) ToSnapshot() es.Snapshot {
	return ProductV1{
		StoreID:     p.storeID,
		Name:        p.name,
		Description: p.description,
		SKU:         p.sku,
		Price:       p.price,
	}
}
