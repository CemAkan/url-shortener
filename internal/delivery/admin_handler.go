package delivery

import (
	"github.com/CemAkan/url-shortener/internal/app"
	"github.com/CemAkan/url-shortener/internal/domain"
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

// ListUsers lists all users nad their links
func (h *AdminHandler) ListUsers(c *fiber.Ctx) error {
	users, err := h.userService.ListAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	type res struct {
		User domain.User
		Urls []domain.URL
	}

	var resp res
	var resps []res

	for _, user := range users {
		urls, _ := h.urlService.GetUserURLs(user.ID)

		resp.User = user
		resp.Urls = urls

		resps = append(resps, resp)
	}

	return c.JSON(resps)
}
