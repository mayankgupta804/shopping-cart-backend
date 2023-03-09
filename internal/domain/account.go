package domain

const (
	AdminRole Role = "Admin"
	UserRole  Role = "User"
)

type Role string

type Account struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Role     Role   `json:"role"`
	Active   bool   `json:"active"`
}
