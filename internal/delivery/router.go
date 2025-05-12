package delivery

import (
	"github.com/CemAkan/url-shortener/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authHandler *AuthHandler, urlHandler *URLHandler) {
	api := app.Group("/api")

	// public routes (no need jwt)

	//redirect
	app.Get("/:code", urlHandler.Redirect)

	//user
	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)

	// protected routes (jwt required)
	api.Get("/me", middleware.JWTAuth(), authHandler.Me)

	//url
	api.Post("/shorten", middleware.JWTAuth(), urlHandler.Shorten)
	api.Get("/my/urls", middleware.JWTAuth(), urlHandler.ListUserURLs)
	api.Get("/my/urls/:code", middleware.JWTAuth(), urlHandler.GetSingleURL)
	api.Delete("/my/urls/:code", middleware.JWTAuth(), urlHandler.DeleteURL)
	api.Patch("/my/urls/:code", middleware.JWTAuth(), urlHandler.UpdateURL)
}
