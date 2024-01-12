package domain

type (
	ItemOption func(*Item)

	Item struct {
		productID   string
		productName string
		storeID     string
		storeName   string
		price       float64
		quantity    int
	}
)

func NewItem(productID string, storeID string, price float64, quantity int, opts ...ItemOption) *Item {
	item := &Item{
		productID: productID,
		storeID:   storeID,
		price:     price,
		quantity:  quantity,
	}

	for _, opt := range opts {
		opt(item)
	}

	return item
}

func (i Item) ProductID() string {
	return i.productID
}

func (i Item) ProductName() string {
	return i.productName
}

func (i Item) StoreID() string {
	return i.storeID
}

func (i Item) StoreName() string {
	return i.storeName
}

func (i Item) Price() float64 {
	return i.price
}

func (i Item) Quantity() int {
	return i.quantity
}

func WithProductName(productName string) ItemOption {
	return func(item *Item) {
		item.productName = productName
	}
}

func WithStoreName(storeName string) ItemOption {
	return func(item *Item) {
		item.storeName = storeName
	}
}
