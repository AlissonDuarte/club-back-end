package models

import (
	"gorm.io/gorm"
)

type UserUpload struct {
	gorm.Model
	UserID      int    `gorm:"not null"`
	FilePath    string `gorm:"not null"`
	FileSize    int64  `gorm:"not null"`
	ContentType string `gorm:"not null"`
}
