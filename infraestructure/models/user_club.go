package models

import (
	"time"

	"gorm.io/gorm"
)

type UserClub struct {
	gorm.Model
	UserID   uint
	ClubID   uint
	CreateAt time.Time
}
