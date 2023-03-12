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

func (itemRepo *itemRepository) List() ([]domain.Item, error) {
	items := make([]domain.Item, 0)
	sql := "SELECT name, sku FROM items;"
	result, err := itemRepo.db.QueryRows(sql)
	if err != nil {
		return nil, err
	}
	for result.Next() {
		item := domain.Item{}
		if err := result.Scan(&item.Name, &item.SKU); err != nil {
			fmt.Printf("error encountered: %v", err)
			continue
		}
		items = append(items, item)
	}

	return items, nil
}

func (itemRepo *itemRepository) Get(itemID string) (domain.Item, error) {
	item := domain.Item{}
	sql := fmt.Sprintf("SELECT * FROM items WHERE id='%s';", itemID)
	result := itemRepo.db.QueryRow(sql)
	if err := result.Scan(&item.ID, &item.Name, &item.SKU); err != nil {
		// format error properly
		return item, err
	}
	return item, nil
}
