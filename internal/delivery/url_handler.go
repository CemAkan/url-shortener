package delivery

import (
	"github.com/CemAkan/url-shortener/internal/app"
	"github.com/gofiber/fiber/v2"
)

type URLhandler struct {
	service app.URLService
}

func NewURLService(urlService app.URLService) *URLhandler {
	return &URLhandler{
		service: urlService,
	}
}

// Shorten handles request to create short url
func (h *URLhandler) Shorten(c *fiber.Ctx) error {
	var req struct {
		OriginalURL string  `json:"original_url"`
		CustomCode  *string `json:"custom_code,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	userID := c.Locals("user_id").(uint)

	url, err := h.service.Shorten(req.OriginalURL, userID, req.CustomCode)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":         url.Code,
		"original_url": req.OriginalURL,
		"short_url":    c.BaseURL() + "/" + url.Code,
	})
}

// ListUserURLs handles request to list user's all urls
func (h *URLhandler) ListUserURLs(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	urls, err := h.service.GetUserURLs(userID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(urls)
}
