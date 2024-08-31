package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name        string `gorm:"default:''"`
	Resume      string `gorm:"default:''"`
	Release     time.Time
	Rate        int `gorm:"default:0"`
	BookCoverID uint
	BookCover   *UserUpload `gorm:"foreignKey:BookCoverID;constraint:OnDelete:CASCADE"`
	AuthorID    uint
	Author      *Author `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
	Readers     []*User `gorm:"many2many:book_readers;constraint:OnDelete:CASCADE"`
	Tags        []*Tag  `gorm:"many2many:book_tags;constraint:OnDelete:CASCADE"`
}
