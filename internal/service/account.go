package service

import (
	"errors"
	"shopping-cart-backend/internal/domain"
	"shopping-cart-backend/internal/serializer"
	"strings"

	"gopkg.in/dealancer/validate.v2"
)

type AccountService interface {
	Create(req serializer.CreateAccountRequest) error
}

type accountService struct {
	acntRepo domain.AccountRepository
}

func NewAccountService(acntRepo domain.AccountRepository) *accountService {
	return &accountService{acntRepo: acntRepo}
}

func (acntServ *accountService) Create(req serializer.CreateAccountRequest) error {
	if err := validate.Validate(req); err != nil {
		return err
	}

	acnt := domain.Account{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     domain.Role(req.Role),
	}

	// any business logic
	// extra logging as necessary
	// storing of metrics for business
	a, err := acntServ.acntRepo.Get(acnt.Email)
	if err != nil {
		return err
	}

	if strings.TrimSpace(a.Email) == acnt.Email {
		return errors.New("account with the same email already exists")
	}

	return acntServ.acntRepo.Create(acnt)
}
