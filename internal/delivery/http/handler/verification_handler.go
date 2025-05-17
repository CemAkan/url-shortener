package handler

import (
	"github.com/CemAkan/url-shortener/internal/domain/request"
	"github.com/CemAkan/url-shortener/internal/domain/response"
	"github.com/CemAkan/url-shortener/internal/service"
	"github.com/gofiber/fiber/v2"
)

type VerificationHandler struct {
	userService service.UserService
}

// NewVerificationHandler generate a new VerificationHandler struct with given UserService and mailService inputs
func NewVerificationHandler(userService service.UserService) *VerificationHandler {
	return &VerificationHandler{
		userService: userService,
	}
}

// VerifyMailAddress godoc
// @Summary Verify user's email address
// @Description Validates email address through verification token
// @Tags Verification
// @Produce json
// @Param token path string true "Verification Token"
// @Success 200 {object} response.SuccessResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /verify/mail/{token} [get]
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

// ResetPassword godoc
// @Summary Reset user password with verification token
// @Description Sets new password after token verification
// @Tags Verification
// @Accept json
// @Produce json
// @Param request body request.NewPassword true "New Password"
// @Success 200 {object} response.SuccessResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /verify/password [post]
func (h *VerificationHandler) ResetPassword(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	if _, err := h.userService.GetByID(userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{Error: "User not found"})
	}

	var req request.NewPassword //new password

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{Error: "invalid request"})
	}

	if err := h.userService.PasswordUpdate(userID, req.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{Error: "password update fail"})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{Message: "password updated"})
}

// ResetPasswordTokenResolve godoc
// @Summary  Return verification token to use reset user password
// @Description Sets new password after token verification
// @Tags Verification
// @Accept json
// @Produce json
// @Param token path string true "Verification Token"
// @Success 200 {object} response.SuccessResponse
// @Router /verify/password/{token} [get]
func (h *VerificationHandler) ResetPasswordTokenResolve(c *fiber.Ctx) error {
	tokenStr := c.Params("token")
	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{Message: tokenStr})
}
