package domain

type Cart struct {
	ID        string `json:"-"`
	AccountID string `json:"account_id"`
	Items     []Item
}
