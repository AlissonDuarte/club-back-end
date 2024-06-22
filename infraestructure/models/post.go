package models

import (
	"errors"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title        string
	Content      string
	UserID       uint
	User         User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	ImageID      uint
	Image        UserUploadPost `gorm:"foreignKey:ImageID"`
	Likes        int
	ClubID       uint      `gorm:"default:null"`
	Club         *Club     `gorm:"foreignKey:ClubID;constraint:OnDelete:CASCADE"`
	Comments     []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	CommentCount int64     `gorm:"-"`
	Updated      bool
}

func NewPost(title string, content string, userID uint, imageID uint, db *gorm.DB) *Post {
	return &Post{
		Title:   title,
		Content: content,
		UserID:  userID,
		ImageID: imageID,
	}
}

func NewPostClub(title string, content string, userID uint, imageID uint, clubID uint, db *gorm.DB) *Post {
	return &Post{
		Title:   title,
		Content: content,
		UserID:  userID,
		ImageID: imageID,
		ClubID:  clubID,
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

	// Append the comment to the post's Comments
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

func PostClubGetByID(db *gorm.DB, id uint, clubID uint) (*Post, error) {
	var post Post
	if err := db.Preload("User", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id, name, username, profile_picture_id").Preload("ProfilePicture", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "file_path")
		})
	}).Preload("Image").Preload("Comments").Preload("Comments.User", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id, name, username, profile_picture_id").Preload("ProfilePicture", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "file_path")
		})
	}).
		Where("club_id = ?", clubID).First(&post, id).Error; err != nil {
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

func IsUserIDInClub(db *gorm.DB, userID uint, clubID uint) (bool, error) {
	var count int64
	result := db.Model(&User{}).
		Joins("JOIN user_club ON users.id = user_club.user_id").
		Where("user_club.club_id = ?", clubID).
		Where("users.id = ?", userID).
		Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
