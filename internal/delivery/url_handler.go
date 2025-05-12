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

// Shorten godoc
// @Summary Shorten a URL
// @Description Create a shortened URL with optional custom code
// @Tags URL
// @Accept json
// @Produce json
// @Param request body request.ShortenURLRequest true "URL to shorten"
// @Success 201 {object} response.URLResponse
// @Failure 400 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /shorten [post]
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

// ListUserURLs godoc
// @Summary Get user's URLs
// @Description Retrieves all shortened URLs for authenticated user
// @Tags URL
// @Produce json
// @Success 200 {array} model.URL
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /my/urls [get]
func (h *URLHandler) ListUserURLs(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	urls, err := h.service.GetUserURLs(userID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(urls)
}

// GetSingleURL godoc
// @Summary Get a single URL detail
// @Description Retrieves a single short URL details with daily click count
// @Tags URL
// @Produce json
// @Param code path string true "Short URL code"
// @Success 200 {object} response.DetailedURLResponse
// @Failure 404 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /my/urls/{code} [get]
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

// Redirect godoc
// @Summary Redirect short URL to original URL
// @Description Redirects to original URL based on short code
// @Tags URL
// @Param code path string true "Short URL code"
// @Success 302 {string} string "Redirects to original URL"
// @Failure 404 {object} response.ErrorResponse
// @Router /{code} [get]
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

// DeleteURL godoc
// @Summary Delete a shortened URL
// @Description Deletes a user's shortened URL by code
// @Tags URL
// @Produce json
// @Param code path string true "Short URL code"
// @Success 200 {object} response.SuccessResponse
// @Failure 403 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /my/urls/{code} [delete]
func (h *URLHandler) DeleteURL(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	code := c.Params("code")

	if err := h.service.DeleteUserURL(userID, code); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(response.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(response.SuccessResponse{Message: "user deleted successfully"})
}

// UpdateURL godoc
// @Summary Update a shortened URL
// @Description Updates original URL or custom code for a user's URL
// @Tags URL
// @Accept json
// @Produce json
// @Param code path string true "Short URL code"
// @Param request body request.UpdateURLRequest true "Updated URL info"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /my/urls/{code} [patch]
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
