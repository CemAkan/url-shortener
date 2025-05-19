package utils

import (
	"net"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GetClientIP extracts the real IP from headers or falls back to c.IP()
func GetClientIP(c *fiber.Ctx) string {
	if xff := c.Get("X-Forwarded-For"); xff != "" {
		// May contain multiple comma-separated IPs
		parts := strings.Split(xff, ",")
		ip := strings.TrimSpace(parts[0])
		return stripPort(ip)
	}

	// If no proxy headers, fallback to real IP
	if xRealIP := c.Get("X-Real-IP"); xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}

	return stripPort(c.IP())
}

func stripPort(ip string) string {
	//port cut
	if strings.Contains(ip, ":") {
		if h, _, err := net.SplitHostPort(ip); err == nil {
			return h
		}
	}
	// IPv6 zone ident %xxx cut
	if i := strings.Index(ip, "%"); i != -1 {
		return ip[:i]
	}
	return ip
}
