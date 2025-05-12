package delivery

import (
	"github.com/CemAkan/url-shortener/internal/app"
	"github.com/CemAkan/url-shortener/internal/domain/request"
	"github.com/CemAkan/url-shortener/internal/domain/response"
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
	var req request.ShortenURLRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{Error: "invalid request"})
	}

	userID := c.Locals("user_id").(uint)

	url, err := h.service.Shorten(req.OriginalURL, userID, req.CustomCode)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{Error: err.Error()})
	}

	var res response.URLResponse
	res.Code, res.OriginalURL, res.ShortURL = url.Code, url.OriginalURL, c.BaseURL()+"/"+url.Code

	return c.Status(fiber.StatusCreated).JSON(res)
}

// ListUserURLs handles request to list user's all urls
func (h *URLHandler) ListUserURLs(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	urls, err := h.service.GetUserURLs(userID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(urls)
}

// GetSingleURL handles request to list a single url details with daily click info
func (h *URLHandler) GetSingleURL(c *fiber.Ctx) error {
	code := c.Params("code")

	url, dailyClicks, err := h.service.GetSingleUrlRecord(code)
	if err != nil || url == nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{Error: "not found"})
	}

	var res response.DetailedURLResponse
	res.OriginalURL, res.Code, res.TotalClicks, res.DailyClicks = url.OriginalURL, url.Code, url.TotalClicks, dailyClicks

	return c.Status(fiber.StatusFound).JSON(res)
}

// Redirect redirects short links to original links
func (h *URLHandler) Redirect(c *fiber.Ctx) error {
	code := c.Params("code")
	ctx := c.Context()

	//code resolving
	originalURL, err := h.service.ResolveRedirect(ctx, code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{Error: "Short url not found"})
	}

	//daily click increment
	go utils.TrackClick(ctx, code)

	return c.Redirect(originalURL, fiber.StatusFound)
}

// DeleteURL handle URL delete requests
func (h *URLHandler) DeleteURL(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	code := c.Params("code")

	if err := h.service.DeleteUserURL(userID, code); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(response.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(response.SuccessResponse{Message: "user deleted successfully"})
}

// UpdateURL handle URL update requests
func (h *URLHandler) UpdateURL(c *fiber.Ctx) error {
	var req request.UpdateURLRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{Error: "invalid request"})
	}

	userID := c.Locals("user_id").(uint)
	code := c.Params("code")

	if err := h.service.UpdateUserURL(userID, code, req.NewOriginalURL, req.NewCustomCode); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(response.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(response.SuccessResponse{Message: "user updated successfully"})

}
