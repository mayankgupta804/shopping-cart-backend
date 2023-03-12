package serializer

type CreateItemRequest struct {
	Name string `validate:"gte=2 & lte=25"`
	SKU  string `validate:"empty=false"`
}

type CreateItemResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
