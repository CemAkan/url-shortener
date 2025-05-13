package model

import "gorm.io/gorm"

type User struct {
	gorm.Model      `json:"-"` // id, created_at, updated_at, deleted_at
	Name            string     `gorm:"size:50;not null"       json:"name"`
	Surname         string     `gorm:"size:50;not null"       json:"surname"`
	Email           string     `gorm:"size:120;uniqueIndex;not null" json:"email"`
	Password        string     `gorm:"not null"                json:"-"`
	IsAdmin         bool       `gorm:"default:false"          json:"is_admin"`
	IsMailConfirmed bool       `gorm:"default:false"          json:"is_verified"`
}
