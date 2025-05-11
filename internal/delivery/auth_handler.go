package delivery

import (
	"github.com/CemAkan/url-shortener/config"
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

// Login handles to request for user logins
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	user, err := h.service.Login(req.Username, req.Password)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// generate jwt token
	token, err := config.GenerateJWT(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"token":    token,
		"id":       user.ID,
		"username": user.Username,
	})
}

// Me returns authenticated user's data
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	// getting userId which comes from middleware
	id := c.Locals("user_id").(uint)

	//user existence check

	user, err := h.service.GetByID(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User not found"})
	}

	// success return with user's data
	return c.JSON(fiber.Map{
		"id":       user.ID,
		"username": user.Username,
	})
}
