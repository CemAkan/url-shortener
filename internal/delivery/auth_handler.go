package delivery

import (
	"github.com/CemAkan/url-shortener/internal/app"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service app.UserService
}

// NewAuthHandler creates a new AuthHandler struct with given UserService input
func NewAuthHandler(userService app.UserService) *AuthHandler {
	return &AuthHandler{
		service: userService,
	}
}

// Register handles user create request
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	user, err := h.service.Register(req.Username, req.Password)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":       user.ID,
		"username": user.Username,
	})
}
