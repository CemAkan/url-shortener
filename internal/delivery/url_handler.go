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
