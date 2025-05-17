package response

type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusDegraded  Status = "degraded"
	StatusUnhealthy Status = "unhealthy"
)

type HealthStatusResponse struct {
	Status   Status `json:"status" example:"healthy"`
	Database Status `json:"database" example:"healthy"`
	Redis    Status `json:"redis" example:"healthy"`
	Email    Status `json:"email" example:"healthy"`
}
