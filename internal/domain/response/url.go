package response

type URLResponse struct {
	Code        string `json:"code" example:"abc123"`
	OriginalURL string `json:"original_url" example:"https://google.com"`
	ShortURL    string `json:"short_url" example:"https://localhost/abc123"`
}

type DetailedURLResponse struct {
	Code        string `json:"code" example:"abc123"`
	OriginalURL string `json:"original_url" example:"https://google.com"`
	TotalClicks int    `json:"total_clicks" example:"42"`
	DailyClicks int    `json:"daily_clicks" example:"10"`
}
