package service

import (
	"shopping-cart-backend/internal/domain"
	"shopping-cart-backend/internal/serializer"
	"strconv"
)

type ItemService interface {
	Add(req serializer.CreateItemRequest) error
	List() ([]serializer.GetItemResponse, error)
}

type itemService struct {
	itemRepo domain.ItemRepository
}

func NewItemService(itemRepo domain.ItemRepository) *itemService {
	return &itemService{itemRepo: itemRepo}
}

func (itemServ *itemService) Add(req serializer.CreateItemRequest) error {
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

func (itemServ *itemService) List() ([]serializer.GetItemResponse, error) {
	items, err := itemServ.itemRepo.List()
	if err != nil {
		return nil, err
	}
	response := make([]serializer.GetItemResponse, 0)
	for _, item := range items {
		sku := strconv.Itoa(int(item.SKU))
		itemResp := serializer.GetItemResponse{
			Name: item.Name,
			SKU:  sku,
		}
		response = append(response, itemResp)
	}
	return response, nil
}
