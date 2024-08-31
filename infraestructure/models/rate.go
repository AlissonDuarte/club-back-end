package models

import "gorm.io/gorm"

type Rate struct {
	gorm.Model
	Rate    int
	Updated bool
	UserID  uint
	User    *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	BookID  uint
	Book    *Book `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE"`
}
