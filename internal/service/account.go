package service

import (
	"errors"
	"shopping-cart-backend/internal/domain"
	"shopping-cart-backend/internal/helpers/password"
	"shopping-cart-backend/internal/serializer"
)

var (
	ErrAccountAlreadyExists error = errors.New("account with this email exists")
	ErrAccountDoesNotExist  error = errors.New("account with this email does not exist")
)

type AccountService interface {
	Create(req serializer.CreateAccountRequest) error
	Suspend(req serializer.SuspendAccountRequest) error
}

type accountService struct {
	acntRepo domain.AccountRepository
}

func NewAccountService(acntRepo domain.AccountRepository) *accountService {
	return &accountService{acntRepo: acntRepo}
}

func (acntServ *accountService) Create(req serializer.CreateAccountRequest) error {
	if acntExists := acntServ.acntRepo.Exists(req.Email); acntExists {
		return ErrAccountAlreadyExists
	}

	passHash, err := password.GenerateHashedPassword([]byte(req.Password))
	if errors.Is(err, password.ErrPasswordTooLong) || errors.Is(err, password.ErrUnexpected) {
		// report errors for observability purposes maybe?
		return err
	}
	acnt := domain.Account{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(passHash),
		Role:     domain.Role(req.Role),
	}

	return acntServ.acntRepo.Create(acnt)
}

func (acntServ *accountService) Suspend(req serializer.SuspendAccountRequest) error {
	if acntExists := acntServ.acntRepo.Exists(req.Email); !acntExists {
		return ErrAccountDoesNotExist
	}

	return acntServ.acntRepo.Suspend(req.Email)
}
