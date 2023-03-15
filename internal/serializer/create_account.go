package serializer

type CreateAccountRequest struct {
	Name     string `validate:"empty=false & gte=2 & lte=25 & format=alpha_unicode" json:"name"`
	Email    string `validate:"empty=false & format=email" json:"email"`
	Password string `validate:"empty=false & gte=7 & lte=25" json:"password"`
	Role     string `validate:"one_of=admin,user" json:"role"`
}

type CreateAccountResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
