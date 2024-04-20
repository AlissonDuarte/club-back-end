package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title    string
	Content  string
	UserID   uint
	User     User `gorm:"foreignKey:UserID"`
	Likes    int
	Comments []Comment
	Updated  bool
}

func NewPost(title string, content string, userID uint, db *gorm.DB) *Post {
	return &Post{
		Title:   title,
		Content: content,
		UserID:  userID,
	}
}

func (p *Post) Save(db *gorm.DB) (uint, error) {
	err := db.Create(p).Error

	if err != nil {
		return 0, err
	}

	return p.ID, nil
}

func (p *Post) Update(db *gorm.DB) error {
	err := db.Save(p).Error
	if err != nil {
		return err
	}

	return nil
}
