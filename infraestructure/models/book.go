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
	BookCoverID uint
	BookCover   *UserUpload `gorm:"foreignKey:BookCoverID;constraint:OnDelete:CASCADE"`
	AuthorID    uint
	Author      *Author `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
	Tags        []*Tag  `gorm:"many2many:book_tags;constraint:OnDelete:CASCADE"`
}

type UserBook struct {
	gorm.Model
	UserID uint   `gorm:"not null"`
	User   *User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Tags   []*Tag `gorm:"many2many:userbook_tags;constraint:OnDelete:CASCADE"`
	Rate   int    `gorm:"default:0"`
}

func NewBook(name string, resume string, release time.Time, coverID uint, authorID uint, db *gorm.DB) *Book {
	return &Book{
		Name:        name,
		Resume:      resume,
		Release:     release,
		BookCoverID: coverID,
		AuthorID:    authorID,
	}
}

func (b *Book) Save(db *gorm.DB) (uint, error) {
	err := db.Create(b).Error
	if err != nil {
		return 0, err
	}
	return b.ID, nil
}

func (b *Book) Update(db *gorm.DB) error {
	err := db.Save(b).Error
	if err != nil {
		return err
	}
	return nil
}
