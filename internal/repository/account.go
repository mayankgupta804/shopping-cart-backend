package repository

import (
	"fmt"
	"shopping-cart-backend/internal/domain"
	"shopping-cart-backend/pkg/database"
)

type accountRepository struct {
	db *database.PgClient
}

func NewAccountRepository(db *database.PgClient) accountRepository {
	return accountRepository{db: db}
}

func (acntRepo accountRepository) Create(acnt domain.Account) error {
	sql := fmt.Sprintf(`INSERT INTO accounts(name, password, email, role) VALUES('%s', '%s', '%s', '%s')`, acnt.Name, acnt.Password, acnt.Email, string(acnt.Role))
	if err := acntRepo.db.Execute(sql); err != nil {
		// format error properly
		return err
	}

	return nil
}

func (acntRepo accountRepository) Get(email string) (domain.Account, error) {
	acnt := domain.Account{}
	sql := fmt.Sprintf("SELECT * FROM accounts WHERE email='%s';", email)
	result := acntRepo.db.QueryRow(sql)
	if err := result.Scan(&acnt.ID, &acnt.Name, &acnt.Email, &acnt.Password, &acnt.Role, &acnt.Active); err != nil {
		// format error properly
		return acnt, err
	}
	return acnt, nil
}
