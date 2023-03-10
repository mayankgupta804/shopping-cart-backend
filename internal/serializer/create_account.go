package serializer

type CreateAccountRequest struct {
	Name     string `validate:"gte=3 & lte=25"`
	Email    string `validate:"empty=false | format=email"`
	Password string `validate:"empty=false & gte=7 & lte=25"`
	Role     string `validate:"one_of=admin,user"`
}

type CreateAccountResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
