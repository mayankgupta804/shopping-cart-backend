package repository

import (
	"fmt"
	"shopping-cart-backend/internal/domain"
	"shopping-cart-backend/pkg/database"
)

type itemRepository struct {
	db *database.PgClient
}

func NewItemRepository(db *database.PgClient) *itemRepository {
	return &itemRepository{db: db}
}

func (itemRepo *itemRepository) Add(item domain.Item) error {
	sql := fmt.Sprintf(`INSERT INTO items(name, sku) VALUES('%s', '%d')`, item.Name, item.SKU)
	if err := itemRepo.db.Execute(sql); err != nil {
		// format error properly
		return err
	}

	return nil
}
