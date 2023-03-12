package serializer

type AddToCartRequest struct {
	ItemID string `json:"item_id"`
}

type AddToCartResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
