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

// recuperar file path com base no userid
func GetUserUploadByUserID(db *gorm.DB, userID int) (UserUpload, error) {
	var userUpload UserUpload
	err := db.Where("user_id = ?", userID).Last(&userUpload).Error
	return userUpload, err
}

type UserUploadPost struct {
	gorm.Model
	UserID   uint   `gorm:"not null"`
	FilePath string `gorm:"not null"`
	FileSize int64  `gorm:"not null"`
}

func GetUserUploadPostByID(db *gorm.DB, id uint) (*UserUploadPost, error) {
	var upload UserUploadPost
	if err := db.First(&upload, id).Error; err != nil {
		return nil, err
	}
	return &upload, nil
}

type UserUploadClub struct {
	gorm.Model
	UserID   uint   `gorm:"not null"`
	FilePath string `gorm:"not null"`
	FileSize int64  `gorm:"not null"`
}
