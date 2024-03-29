package serializer

type AddToCartRequest struct {
	ItemID string `validate:"gte=1 & format=alnum_unicode" json:"item_id"`
}

type AddToCartResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
