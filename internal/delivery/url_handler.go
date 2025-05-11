package delivery

import (
	"github.com/CemAkan/url-shortener/internal/app"
	"github.com/CemAkan/url-shortener/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type URLHandler struct {
	service app.URLService
}

func NewURLHandler(urlService app.URLService) *URLHandler {
	return &URLHandler{
		service: urlService,
	}
}

// Shorten handles request to create short url
func (h *URLHandler) Shorten(c *fiber.Ctx) error {
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
func (h *URLHandler) ListUserURLs(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	urls, err := h.service.GetUserURLs(userID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(urls)
}

// GetSingleURL handles request to list a single url details with daily click info
func (h *URLHandler) GetSingleURL(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	code := c.Params("code")

	url, dailyClicks, err := h.service.GetSingleUrlRecord(code, userID)
	if err != nil || url == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}

	return c.Status(fiber.StatusFound).JSON(fiber.Map{
		"code":         url.Code,
		"original_url": url.OriginalURL,
		"total_clicks": url.TotalClicks,
		"daily_clicks": dailyClicks,
	})
}

// Redirect redirects short links to original links
func (h *URLHandler) Redirect(c *fiber.Ctx) error {
	code := c.Params("code")
	ctx := c.Context()

	//code resolving
	originalURL, err := h.service.ResolveRedirect(ctx, code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Short URL not found"})
	}

	//daily click increment
	go utils.TrackClick(ctx, code)

	return c.Redirect(originalURL, fiber.StatusFound)
}
