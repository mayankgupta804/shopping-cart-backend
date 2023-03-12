package serializer

type AddToCartRequest struct {
	ItemID string `validate:"empty=false"`
}

type AddToCartResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
