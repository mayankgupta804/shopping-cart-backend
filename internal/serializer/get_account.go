package serializer

type GetAccountRequest struct {
	Email    string
	Password string
}

type GetAccountResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
