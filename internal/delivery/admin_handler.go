package delivery

import (
	"github.com/CemAkan/url-shortener/internal/app"
	"github.com/CemAkan/url-shortener/internal/domain/response"
	"github.com/gofiber/fiber/v2"
	"strconv"
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

// ListUsers lists all users nad their links
func (h *AdminHandler) ListUsers(c *fiber.Ctx) error {
	users, err := h.userService.ListAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{Error: err.Error()})
	}

	var resps []response.UserURLsResponse

	for _, user := range users {
		urls, _ := h.urlService.GetUserURLs(user.ID)

		resps = append(resps, response.UserURLsResponse{
			User: user,
			Urls: urls,
		})
	}

	return c.JSON(resps)
}

// RemoveUser delete selected user record
func (h *AdminHandler) RemoveUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{Error: "invalid user id"})
	}

	userID := uint(id)

	if err := h.urlService.DeleteUserAllURLs(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{Error: err.Error()})
	}

	if err := h.userService.DeleteUser(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(response.SuccessResponse{Message: "user deleted successfully"})
}
