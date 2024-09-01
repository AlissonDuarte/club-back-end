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
	Certified   bool
}

type UserBook struct {
	gorm.Model
	BookID uint   `gorm:"not null"`
	Book   *Book  `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE"`
	UserID uint   `gorm:"not null"`
	User   *User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Tags   []*Tag `gorm:"many2many:userbook_tags;constraint:OnDelete:CASCADE"`
	Rate   int    `gorm:"default:0"`
}

func NewBook(name string, resume string, release time.Time, coverID uint, authorID uint, certified bool, db *gorm.DB) *Book {
	return &Book{
		Name:        name,
		Resume:      resume,
		Release:     release,
		BookCoverID: coverID,
		AuthorID:    authorID,
		Certified:   certified,
	}
}

func NewUserBook(bookID uint, userID uint, db *gorm.DB) *UserBook {
	return &UserBook{
		BookID: bookID,
		UserID: userID,
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

func (ub *UserBook) Save(db *gorm.DB) (uint, error) {
	err := db.Create(ub).Error
	if err != nil {
		return 0, err
	}
	return ub.ID, nil
}

func (ub *UserBook) Update(db *gorm.DB) error {
	err := db.Save(ub).Error
	if err != nil {
		return err
	}
	return nil
}

func BookGetByID(db *gorm.DB, bookID uint) (*Book, error) {
	var book Book
	if err := db.Preload("Tags").Preload("BookCover", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "file_path")
	}).Preload("Author", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "name", "profile_picture_id").Preload("ProfilePicture", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "file_path")
		})
	}).First(&book, bookID).Error; err != nil {
		return nil, err
	}
	return &book, nil
}
