package domain

type AccountRepository interface {
	Create(acnt Account) error
	Get(userName, password string) (Account, error)
}

const (
	AdminRole Role = "Admin"
	UserRole  Role = "User"
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
