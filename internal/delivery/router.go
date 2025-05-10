package delivery

import (
	"github.com/CemAkan/url-shortener/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authHandler *AuthHandler) {
	api := app.Group("/api")

	// public routes (no need jwt)
	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)

	// protected routes (jwt required)
	api.Get("/me", middleware.JWTAuth(), authHandler.Me)
}
