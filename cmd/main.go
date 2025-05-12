package main

import (
	"github.com/CemAkan/url-shortener/config"
	appModule "github.com/CemAkan/url-shortener/internal/app"
	"github.com/CemAkan/url-shortener/internal/delivery"
	"github.com/CemAkan/url-shortener/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func main() {

	config.LoadEnv()
	config.InitLogger()
	config.InitDB()
	config.InitRedis()

	app := fiber.New()

	// Dependency injection

	//user
	userRepo := repository.NewUserRepository()
	userService := appModule.NewUserService(userRepo)
	userHandler := delivery.NewAuthHandler(userService)

	//url
	urlRepo := repository.NewURLRepository()
	urlService := appModule.NewURLService(urlRepo)
	urlHandler := delivery.NewURLHandler(urlService)

	// Routes
	delivery.SetupRoutes(app, userHandler, urlHandler)

	port := config.GetEnv("PORT", "3000")

	err := app.Listen(":" + port)
	if err != nil {
		config.Log.Fatalf("Port can not listen: %v", err)
		return
	}
}
