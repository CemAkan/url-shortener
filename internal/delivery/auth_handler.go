package delivery

import (
	"github.com/CemAkan/url-shortener/config"
	"github.com/CemAkan/url-shortener/internal/app"
	"github.com/CemAkan/url-shortener/internal/domain/request"
	"github.com/CemAkan/url-shortener/internal/domain/response"
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

// Register godoc
// @Summary User Registration
// @Description Creates a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.AuthRequest true "User Credentials"
// @Success 201 {object} response.UserResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req request.AuthRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{Error: "invalid request body"})
	}

	user, err := h.service.Register(req.Username, req.Password)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{Error: err.Error()})
	}

	var res response.UserResponse
	res.ID, res.Username = user.ID, user.Username

	return c.Status(fiber.StatusCreated).JSON(res)
}

// Login godoc
// @Summary User Login
// @Description Authenticates a user and returns JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.AuthRequest true "User Credentials"
// @Success 201 {object} response.LoginResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req request.AuthRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{Error: "invalid request body"})
	}

	user, err := h.service.Login(req.Username, req.Password)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{Error: err.Error()})
	}

	// generate jwt token
	token, err := config.GenerateJWT(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{Error: "failed to generate token"})
	}

	var res response.LoginResponse
	res.ID, res.Username, res.Token = user.ID, user.Username, token

	return c.Status(fiber.StatusCreated).JSON(res)
}

// Me godoc
// @Summary Get current user's profile
// @Description Returns authenticated user's profile info
// @Tags Auth
// @Produce json
// @Success 200 {object} response.UserResponse
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /me [get]
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	// getting userId which comes from middleware
	id := c.Locals("user_id").(uint)

	//user existence check

	user, err := h.service.GetByID(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{Error: "User not found"})
	}

	var res response.UserResponse
	res.ID, res.Username = user.ID, user.Username

	// success return with user's data
	return c.JSON(res)
}
