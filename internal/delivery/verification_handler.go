package delivery

import (
	"github.com/CemAkan/url-shortener/internal/app"
	"github.com/CemAkan/url-shortener/internal/domain/response"
	"github.com/gofiber/fiber/v2"
)

type VerificationHandler struct {
	userService app.UserService
}

// NewVerificationHandler generate a new VerificationHandler struct with given UserService and mailService inputs
func NewVerificationHandler(userService app.UserService) *VerificationHandler {
	return &VerificationHandler{
		userService: userService,
	}
}

func (h *VerificationHandler) VerifyMailAddress(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	if _, err := h.userService.GetByID(userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{Error: "User not found"})
	}

	if err := h.userService.SetTrueEmailConfirmation(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{Error: "Database error"})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{Message: "mail confirmation successfully"})
}

func (h *VerificationHandler) ResetPassword(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	if _, err := h.userService.GetByID(userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{Error: "User not found"})
	}

	var req string //new password

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{Error: "invalid request"})
	}

	if err := h.userService.PasswordUpdate(userID, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{Error: "password update fail"})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{Message: "password updated"})
}
