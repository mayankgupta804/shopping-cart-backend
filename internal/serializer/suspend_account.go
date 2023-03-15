package serializer

type SuspendAccountRequest struct {
	Email string `validate:"empty=false & format=email" json:"email"`
}

type SuspendAccountResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
