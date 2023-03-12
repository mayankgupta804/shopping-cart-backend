package serializer

type RemoveFromCartRequest struct {
	ItemID string `validate:"empty=false"`
}

type RemoveFromCartResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
