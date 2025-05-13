package delivery

import (
	"github.com/CemAkan/url-shortener/internal/domain/response"
	"github.com/CemAkan/url-shortener/internal/health"
	"github.com/gofiber/fiber/v2"
)

// Health godoc
// @Summary Health Check
// @Description Returns current health status of DB and Redis
// @Tags Health
// @Success 200 {object} response.HealthStatusResponse
// @Router /health [get]
func Health(c *fiber.Ctx) error {
	dbStatus := "ok"
	if !health.GetDBStatus() {
		dbStatus = "error"
	}

	redisStatus := "ok"
	if !health.GetRedisStatus() {
		redisStatus = "error"
	}

	status := "healthy"
	if dbStatus != "ok" || redisStatus != "ok" {
		status = "degraded"
	}

	return c.JSON(response.HealthStatusResponse{
		Status:   status,
		Database: dbStatus,
		Redis:    redisStatus,
	})
}
