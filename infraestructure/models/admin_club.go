package models

import (
	"time"

	"gorm.io/gorm"
)

type AdminClub struct {
	gorm.Model
	UserID   uint
	ClubID   uint
	CreateAt time.Time
}
