package domain

type Items map[string]*Item

type Item struct {
	StoreID     string
	StoreName   string
	ProductID   string
	ProductName string
	Price       float64
	Quantity    int
}
