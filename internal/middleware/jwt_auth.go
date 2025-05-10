package middleware

import (
	"github.com/CemAkan/url-shortener/config"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")

		//Authorization header format check
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid Authorization header"})
		}

		//Token check
		tokenStr := strings.TrimPrefix(authHeader, " Bearer ")

		token, err := config.ValidateJWT(tokenStr)

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		return c.Next()
	}
}
