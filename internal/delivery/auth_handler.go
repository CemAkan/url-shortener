package delivery

import (
	"github.com/CemAkan/url-shortener/config"
	"github.com/CemAkan/url-shortener/internal/app"
	"github.com/CemAkan/url-shortener/internal/domain/request"
	"github.com/CemAkan/url-shortener/internal/domain/response"
	"github.com/gofiber/fiber/v2"
	"time"
)

type AuthHandler struct {
	userService app.UserService
	mailService app.MailService
}

// NewAuthHandler creates a new AuthHandler struct with given UserService input
func NewAuthHandler(userService app.UserService, mailService app.MailService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		mailService: mailService,
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

	user, err := h.userService.Register(req.Email, req.Password, req.Name, req.Surname)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{Error: err.Error()})
	}

	// get baseURL from fiber
	baseURL := c.BaseURL()

	//verify link generator
	verifyLink, err := h.mailService.VerifyLinkGenerator(user.ID, baseURL)

	if err != nil {
		h.mailService.GetMailLogger().Warnf("verify token generation for %s mail address failed: %v", user.Email, err.Error())
	}

	// email address verification mail sending
	if err := h.mailService.SendVerificationMail(user.Name, user.Email, verifyLink); err != nil {
		h.mailService.GetMailLogger().Warnf("send verification mail to %s mail address failed: %v", user.Email, err.Error())
	}

	var res response.UserResponse
	res.ID, res.Email = user.ID, user.Email

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

	user, err := h.userService.Login(req.Email, req.Password)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{Error: err.Error()})
	}

	// generate jwt token
	token, err := config.GenerateToken(user.ID, time.Duration(24*time.Hour), "auth")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{Error: "failed to generate token"})
	}

	var res response.LoginResponse
	res.ID, res.Email, res.Token = user.ID, user.Email, token

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

	user, err := h.userService.GetByID(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{Error: "User not found"})
	}

	var res response.UserResponse
	res.ID, res.Email = user.ID, user.Email

	// success return with user's data
	return c.JSON(res)
}
