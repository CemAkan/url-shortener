package model

import "gorm.io/gorm"

type URL struct {
	gorm.Model
	Code        string `gorm:"uniqueIndex;not null"`
	OriginalURL string `gorm:"not null"`
	UserID      uint   `gorm:"index"`
	TotalClicks int    `gorm:"default:0"`
}
