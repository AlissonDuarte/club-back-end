package models

import (
	"clube/internal/responses"

	"gorm.io/gorm"
)

// Grupo
type Club struct {
	gorm.Model
	Name        string
	Description string
	ImageID     uint
	Image       *UserUploadClub `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"`
	OwnerID     uint
	OwnerRefer  *User   `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"`
	Users       []*User `gorm:"many2many:user_club;constraint:OnDelete:CASCADE"`
}

func NewClub(name string, description string, userIds []int, owner uint, imageID uint, db *gorm.DB) *Club {
	users := []*User{}

	for _, userID := range userIds {
		user := &User{}
		db.First(&user, userID)

		if user.ID != 0 {
			users = append(users, user)
		}
	}

	return &Club{
		Name:        name,
		Description: description,
		Users:       users,
		OwnerID:     owner,
		ImageID:     imageID,
	}
}

func (c *Club) Save(db *gorm.DB) error {
	return db.Create(c).Error
}

func (c *Club) Update(db *gorm.DB) error {
	return db.Save(c).Error
}

func ClubGetById(db *gorm.DB, id int) (*Club, error) {
	var club Club
	err := db.Preload("OwnerRefer",
		func(tx *gorm.DB) *gorm.DB {
			return tx.Select("ID", "Username", "Gender", "BirthDate")
		}).Preload("Users", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("ID", "Username", "Gender", "BirthDate")
	}).First(&club, id).Error

	if err != nil {
		return nil, err
	}

	return &club, nil
}

func GetClubFeed(db *gorm.DB, clubID uint, offset, limit int) ([]responses.FeedResponse, error) {
	var posts []Post
	var responseData []responses.FeedResponse

	err := db.Preload("User", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "name", "username", "profile_picture_id").Preload("ProfilePicture", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "file_path")
		})
	}).Preload("Image").Where("club_id = ?", clubID).Order("id desc").Offset(offset).Limit(limit).Find(&posts).Error

	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		var commentCount int64
		err = db.Model(&Comment{}).Where("post_id = ?", post.ID).Count(&commentCount).Error
		if err != nil {
			return nil, err
		}

		responseData = append(responseData, responses.FeedResponse{
			ID:           post.ID,
			Title:        post.Title,
			Content:      post.Content,
			UserID:       post.UserID,
			ImageID:      post.ID,
			CommentCount: commentCount,
			UpdatedAt:    post.UpdatedAt,
			CreatedAt:    post.CreatedAt,

			User: responses.FeedUserData{
				ID:               post.User.ID,
				Name:             post.User.Name,
				Username:         post.User.Username,
				ProfilePictureID: post.User.ProfilePictureID,

				ProfilePicture: responses.FeedProfilePicture{
					ID:       post.User.ProfilePicture.ID,
					FilePath: post.User.ProfilePicture.FilePath,
				},
			},
		})
	}
	for i := range posts {
		var commentCount int64
		err = db.Model(&Comment{}).Where("post_id = ?", posts[i].ID).Count(&commentCount).Error
		if err != nil {
			return nil, err
		}
		posts[i].CommentCount = commentCount
	}

	return responseData, nil
}
func GetClubUploadByID(db *gorm.DB, clubID uint) (*UserUploadClub, error) {
	var club Club
	if err := db.First(&club, clubID).Error; err != nil {
		return nil, err
	}

	var upload UserUploadClub
	if err := db.First(&upload, club.ImageID).Error; err != nil {
		return nil, err
	}
	return &upload, nil

}
