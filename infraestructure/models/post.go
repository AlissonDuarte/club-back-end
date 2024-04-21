package models

import (
	"errors"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title    string
	Content  string
	UserID   uint
	User     User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	ImageID  uint
	Image    UserUploadPost `gorm:"foreignKey:ImageID"`
	Likes    int
	Comments []Comment `gorm:"many2many:comment_post;constraint:OnDelete:CASCADE"`
	Updated  bool
}

func NewPost(title string, content string, userID uint, imageID uint, db *gorm.DB) *Post {
	return &Post{
		Title:   title,
		Content: content,
		UserID:  userID,
		ImageID: imageID,
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

func AddCommentToPost(db *gorm.DB, postID uint, commentID uint) error {
	var post Post
	if err := db.First(&post, postID).Error; err != nil {
		return err
	}

	var comment Comment
	if err := db.First(&comment, commentID).Error; err != nil {
		return err
	}

	if err := db.Model(&post).Association("Comments").Append(&comment); err != nil {
		return err
	}

	return nil
}

func PostGetByID(db *gorm.DB, id uint) (*Post, error) {
	var post Post
	if err := db.Preload("User", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id, name ,username", "profile_picture_id").Preload("ProfilePicture", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "file_path")
		})
	}).Preload("Image").Preload("Comments").Preload("Comments.User", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id, name, username, profile_picture_id").Preload("ProfilePicture", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "file_path")
		})
	}).
		First(&post, id).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func GetPostUploadByPostID(db *gorm.DB, postID uint, userID uint) (*UserUploadPost, error) {
	var post Post
	if err := db.First(&post, postID).Error; err != nil {
		return nil, err
	}

	if post.UserID != userID {
		return nil, errors.New("post does not belong to user")
	}

	var upload UserUploadPost
	if err := db.First(&upload, post.ImageID).Error; err != nil {
		return nil, err
	}

	return &upload, nil
}
