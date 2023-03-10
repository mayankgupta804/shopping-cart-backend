package domain

type AccountRepository interface {
	Create(acnt Account) error
	Get(email string) (Account, error)
	Exists(email string) bool
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
