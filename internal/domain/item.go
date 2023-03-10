package domain

type ItemRepository interface {
	Add(item Item) error
	List() ([]Item, error)
}

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	SKU  int64  `json:"sku"`
}
