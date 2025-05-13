package delivery

import (
	"github.com/CemAkan/url-shortener/internal/app"
	"github.com/CemAkan/url-shortener/internal/domain/response"
	"github.com/gofiber/fiber/v2"
)

type VerificationHandler struct {
	userService app.UserService
	mailService app.MailService
}

// NewVerificationHandler generate a new VerificationHandler struct with given UserService and mailService inputs
func NewVerificationHandler(userService app.UserService, mailService app.MailService) *VerificationHandler {
	return &VerificationHandler{
		userService: userService,
		mailService: mailService,
	}
}

func (h *VerificationHandler) VerifyMailAddress(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	_, err := h.userService.GetByID(userID)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{Error: "User not found"})
	}

	if err := h.userService.SetTrueEmailConfirmation(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{Error: "Database error"})
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse{Message: "mail confirmation successfully"})
}
