package domain

type AccountRepository interface {
	Create(acnt Account) error
	Get(email string) (Account, error)
	Exists(email string) bool
	Suspend(email string) error
}

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

type Role string

type Account struct {
	ID       string
	Name     string
	Email    string
	Password string
	Active   bool
	Role     Role
}
