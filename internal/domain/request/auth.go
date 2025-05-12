package request

type AuthRequest struct {
	Username string `json:"username" example:"asko"`
	Password string `json:"password" example:"supersecret"`
}
