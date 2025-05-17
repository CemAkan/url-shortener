package handler

import (
	"github.com/CemAkan/url-shortener/internal/domain/response"
	"github.com/CemAkan/url-shortener/internal/health"
	"github.com/gofiber/fiber/v2"
)

// Health godoc
// @Summary Health Check
// @Description Returns current health status of DB, Redis and Email services
// @Tags Health
// @Success 200 {object} response.HealthStatusResponse
// @Router /health [get]
func Health(c *fiber.Ctx) error {
	dbStatus := response.StatusHealthy
	if !health.GetDBStatus() {
		dbStatus = response.StatusUnhealthy
	}

	redisStatus := response.StatusHealthy
	if !health.GetRedisStatus() {
		redisStatus = response.StatusUnhealthy
	}

	emailStatus := response.StatusHealthy
	if !health.GetEmailStatus() {
		emailStatus = response.StatusUnhealthy
	}
	status := response.StatusHealthy
	if dbStatus != response.StatusHealthy || redisStatus != response.StatusHealthy || emailStatus != response.StatusHealthy {
		status = response.StatusDegraded
	}

	return c.JSON(response.HealthStatusResponse{
		Status:   status,
		Database: dbStatus,
		Redis:    redisStatus,
		Email:    emailStatus,
	})
}
