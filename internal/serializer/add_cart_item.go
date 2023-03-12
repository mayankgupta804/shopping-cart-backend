package serializer

type AddToCartRequest struct {
	ItemID string `json:"item_id"`
}

type AddToCartResponse struct {
	Response
}
