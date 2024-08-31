package models

import "gorm.io/gorm"

type Author struct {
	gorm.Model
	Name             string
	Resume           string
	Rate             int
	ProfilePictureID uint        `gorm:"default:null"`
	ProfilePicture   *UserUpload `gorm:"foreignKey:ProfilePictureID"`
	Books            []Book      `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
}
