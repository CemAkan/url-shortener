package request

// AuthRequest  ➜  Kayıt / giriş istek gövdesi
type AuthRequest struct {
	Name     string `json:"name"     example:"Cem"`
	Surname  string `json:"surname"  example:"Akan"`
	Email    string `json:"email"    validate:"required,email" example:"asko@kusko.com"`
	Password string `json:"password" validate:"required,min=8" example:"supersecret"`
}
