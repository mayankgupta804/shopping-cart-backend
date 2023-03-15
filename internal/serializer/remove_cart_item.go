package serializer

type RemoveFromCartRequest struct {
	ItemID string `validate:"gte=1 & format=alnum_unicode" json:"item_id"`
}

type RemoveFromCartResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
