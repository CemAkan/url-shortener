package middleware

import (
	"os"
	"strings"

	"github.com/CemAkan/url-shortener/pkg/logger"
	"github.com/CemAkan/url-shortener/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func IPWhitelistMiddleware() fiber.Handler {
	whitelist := os.Getenv("IP_WHITELIST") // comma-separated IPs
	allowedIPs := strings.Split(whitelist, ",")

	return func(c *fiber.Ctx) error {
		ip := utils.GetClientIP(c)
		for _, allowed := range allowedIPs {
			if strings.TrimSpace(ip) == strings.TrimSpace(allowed) {
				return c.Next()
			}
		}

		logger.Log.Warnf("Unauthorized IP: %s attempted to access protected route", ip)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
}
