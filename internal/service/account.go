package service

import (
	"errors"
	"shopping-cart-backend/internal/domain"
	"shopping-cart-backend/internal/serializer"

	"gopkg.in/dealancer/validate.v2"
)

var ErrAccountAlreadyExists error = errors.New("account with this email exists")

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
	if exists := acntServ.acntRepo.Exists(acnt.Email); exists {
		return ErrAccountAlreadyExists
	}

	return acntServ.acntRepo.Create(acnt)
}
