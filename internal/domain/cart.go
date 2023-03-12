package domain

type CartRepository interface {
	Add(accountID string, item Item) error
	Remove(accountID string, item Item) error
}
