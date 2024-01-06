package domain

type Items map[string]*Item

type Item struct {
	storeID     string
	storeName   string
	productID   string
	productName string
	price       float64
	quantity    int
}

func (i Item) StoreID() string {
	return i.storeID
}

func (i Item) StoreName() string {
	return i.storeName
}

func (i Item) ProductID() string {
	return i.productID
}

func (i Item) ProductName() string {
	return i.productName
}

func (i Item) Price() float64 {
	return i.price
}

func (i Item) Quantity() int {
	return i.quantity
}
