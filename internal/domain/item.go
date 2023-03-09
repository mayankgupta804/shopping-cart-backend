package domain

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	SKU  int64  `json:"sku"`
}
