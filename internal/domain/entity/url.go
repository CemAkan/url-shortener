package entity

import "gorm.io/gorm"

type URL struct {
	gorm.Model  `json:"-"` // ID, CreatedAt, UpdatedAt, DeletedAt
	Code        string     `gorm:"uniqueIndex;not null" json:"code"`
	OriginalURL string     `gorm:"not null"           json:"original_url"`
	UserID      uint       `gorm:"index"              json:"user_id"`
	TotalClicks int        `gorm:"default:0"          json:"total_clicks"`
}
