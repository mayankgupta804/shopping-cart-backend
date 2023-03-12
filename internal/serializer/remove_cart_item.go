package serializer

type RemoveFromCartRequest struct {
	ItemID string
}

type RemoveFromCartResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
