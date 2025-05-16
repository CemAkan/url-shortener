package response

type HealthStatusResponse struct {
	Status   string `json:"status" example:"healthy"`
	Database string `json:"database" example:"ok"`
	Redis    string `json:"redis" example:"ok"`
	Email    string `json:"email" example:"ok"`
}
