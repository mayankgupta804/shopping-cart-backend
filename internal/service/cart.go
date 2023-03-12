package service

import (
	"errors"
	"shopping-cart-backend/internal/domain"
	"shopping-cart-backend/internal/serializer"
)

type CartService interface {
	Add(accountID string, req serializer.AddToCartRequest) error
	Remove(accountID string, req serializer.RemoveFromCartRequest) error
}

type cartService struct {
	cartRepo domain.CartRepository
	itemRepo domain.ItemRepository
}

func NewCartService(cartRepo domain.CartRepository, itemRepo domain.ItemRepository) *cartService {
	return &cartService{cartRepo: cartRepo, itemRepo: itemRepo}
}

func (cartServ *cartService) Add(accountID string, req serializer.AddToCartRequest) error {
	item, err := cartServ.itemRepo.Get(req.ItemID)
	if err != nil {
		return err
	}
	if !item.IsInStock() {
		return errors.New("item cannot be added to cart: no stock")
	}
	return cartServ.cartRepo.Add(accountID, item)
}

func (cartServ *cartService) Remove(accountID string, req serializer.RemoveFromCartRequest) error {
	cartItem := domain.Item{
		ID: req.ItemID,
	}

	return cartServ.cartRepo.Remove(accountID, cartItem)
}
