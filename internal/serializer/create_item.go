package serializer

type CreateItemRequest struct {
	Name string `validate:"empty=false & gte=2 & lte=50 & format=alnum_unicode" json:"name"`
	SKU  string `validate:"empty=false & gte=1 & format=number" json:"sku"`
}

type CreateItemResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
