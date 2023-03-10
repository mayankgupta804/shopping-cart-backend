package service

import (
	"shopping-cart-backend/internal/domain"
	"shopping-cart-backend/internal/serializer"
	"strconv"

	"gopkg.in/dealancer/validate.v2"
)

type ItemService interface {
	Add(req serializer.CreateItemRequest) error
}

type itemService struct {
	itemRepo domain.ItemRepository
}

func NewItemService(itemRepo domain.ItemRepository) *itemService {
	return &itemService{itemRepo: itemRepo}
}

func (itemServ *itemService) Add(req serializer.CreateItemRequest) error {
	if err := validate.Validate(req); err != nil {
		return err
	}

	sku, err := strconv.Atoi(req.SKU)
	if err != nil {
		// format error properly
		return err
	}

	item := domain.Item{
		Name: req.Name,
		SKU:  int64(sku),
	}

	// any business logic
	// extra logging as necessary
	// storing of metrics for business
	return itemServ.itemRepo.Add(item)
}
