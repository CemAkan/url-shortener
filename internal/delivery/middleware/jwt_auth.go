package middleware

import (
	"github.com/CemAkan/url-shortener/config"
	"github.com/CemAkan/url-shortener/internal/domain/response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

func JWTAuth(purpose string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")

		//Authorization header format check
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{Error: "Missing or invalid Authorization header"})
		}

		// Trimming
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		tokenStr = strings.TrimSpace(tokenStr)

		// Token check
		token, err := config.ResolveToken(tokenStr)

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{Error: "Invalid or expired token"})

		}

		// map data claim
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
				Error: "Invalid token claims",
			})
		}

		if claims["type"] != purpose {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{Error: "invalid token type"})
		}

		// userID claim

		userID, ok := claims["user_id"].(float64)

		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
				Error: "Invalid token user ID",
			})
		}

		c.Locals("user_id", uint(userID))

		return c.Next()

	}
}
