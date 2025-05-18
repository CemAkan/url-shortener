package utils

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GetClientIP extracts the real IP from headers or falls back to c.IP()
func GetClientIP(c *fiber.Ctx) string {
	xForwardedFor := c.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// May contain multiple comma-separated IPs
		parts := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(parts[0])
	}

	// If no proxy headers, fallback to real IP
	xRealIP := c.Get("X-Real-IP")
	if xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}

	// Fallback to remote IP
	return c.IP()
}
