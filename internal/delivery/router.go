package delivery

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App, authHandler *AuthHandler) {
	app.Group("/api")

	app.Post("/register", authHandler.Register)
	app.Post("/login", authHandler.Login)
}
