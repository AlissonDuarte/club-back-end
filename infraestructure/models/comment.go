package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	UserID  uint
	User    User `gorm:"foreignKey:UserID"`
	PostID  uint
	Post    Post `gorm:"foreignKey:PostID"`
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

func (c *Comment) Save(db *gorm.DB) error {

	err := db.Create(c).Error
	if err != nil {
		return err
	}

	return nil

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
