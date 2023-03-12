package repository

import (
	"fmt"
	"shopping-cart-backend/internal/domain"
	"shopping-cart-backend/pkg/database"
	"strconv"
)

type cartRepository struct {
	db *database.PgClient
}

func NewCartRepository(db *database.PgClient) *cartRepository {
	return &cartRepository{db: db}
}

func (cartRepo *cartRepository) Add(accountID string, item domain.Item) error {
	sql := fmt.Sprintf(`INSERT INTO carts(account_id, item_id, item_name) VALUES('%s', '%s', '%s')`, accountID, item.ID, item.Name)
	if err := cartRepo.db.Execute(sql); err != nil {
		// format error properly
		return err
	}
	return nil
}

func (cartRepo *cartRepository) Remove(accountID string, item domain.Item) error {
	acntID, _ := strconv.Atoi(accountID)
	itemID, _ := strconv.Atoi(item.ID)
	sql := fmt.Sprintf(`DELETE FROM carts WHERE account_id='%d' AND item_id='%d'`, acntID, itemID)
	if err := cartRepo.db.Execute(sql); err != nil {
		// format error properly
		return err
	}
	return nil
}
