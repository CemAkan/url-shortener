package response

type ErrorResponse struct {
	Error string `json:"error" example:"invalid credentials"`
}
