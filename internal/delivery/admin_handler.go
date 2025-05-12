package delivery

import (
	"github.com/CemAkan/url-shortener/internal/app"
	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	userService app.UserService
	urlService  app.URLService
}

// NewAdminHandler constructor
func NewAdminHandler(userService app.UserService, urlService app.URLService) *AdminHandler {
	return &AdminHandler{
		userService: userService,
		urlService:  urlService,
	}
}

// ListUsers lists all users
func (h *AdminHandler) ListUsers(c *fiber.Ctx) error {
	users, err := h.userService.ListAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(users)
}
