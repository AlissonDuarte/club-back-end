package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID  uint
	User    User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	PostID  uint
	Post    Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	Content string
	Updated bool
}

func NeWComment(userID uint, postID uint, content string) *Comment {
	return &Comment{
		UserID:  userID,
		PostID:  postID,
		Content: content,
	}
}

func (c *Comment) Save(db *gorm.DB) (uint, error) {

	err := db.Create(c).Error
	if err != nil {
		return 0, err
	}

	return c.ID, nil

}

func GetPostComment(db *gorm.DB, postID int) (*Post, error) {
	var post Post
	err := db.Preload("Comments", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("ID", "Content")
	}).First(&post, postID).Error

	if err != nil {
		return nil, err
	}
	return &post, err
}

func GetCommentByID(db *gorm.DB, id uint, userID uint, postID uint) (*Comment, error) {
	var comment Comment
	err := db.Where("id = ? AND user_id = ? AND post_id = ?", id, userID, postID).First(&comment).Error

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (p *Comment) Update(db *gorm.DB) error {
	err := db.Save(p).Error

	if err != nil {
		return err
	}

	return nil
}
