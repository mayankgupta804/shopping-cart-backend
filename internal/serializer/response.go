package serializer

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
