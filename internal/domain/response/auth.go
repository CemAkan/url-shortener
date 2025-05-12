package response

type LoginResponse struct {
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"asko"`
	Token    string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type UserResponse struct {
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"asko"`
}
