package request

type ShortenURLRequest struct {
	OriginalURL string  `json:"original_url" example:"https://google.com"`
	CustomCode  *string `json:"custom_code,omitempty" example:"custom123"`
}

type UpdateURLRequest struct {
	NewOriginalURL *string `json:"new_original_url,omitempty" example:"https://updated.com"`
	NewCustomCode  *string `json:"new_custom_code,omitempty" example:"newcode123"`
}
