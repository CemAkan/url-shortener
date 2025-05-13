package request

type NewPassword struct {
	Password string `json:"password" example:"newsecurepassword"`
}
