package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title    string
	Content  string
	UserID   uint
	User     User `gorm:"foreignKey:UserID"`
	ImageID  uint
	Image    UserUploadPost `gorm:"foreignKey:ImageID"`
	Likes    int
	Comments []Comment `gorm:"many2many:comment_post"`
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

func PostGetByID(db *gorm.DB, id int) (*Post, error) {
	var post Post
	err := db.Preload("User",
		func(tx *gorm.DB) *gorm.DB {
			return tx.Select("ID", "Username").Joins("JOIN user_uploads ON users.id = user_uploads.user_id").Select("users.*, user_uploads.file_path")
		}).Preload("Comments", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("ID", "Content", "Created_at", "Updated_at", "Updated")
	}).Preload("Image", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("ID", "FilePath")
	}).First(&post, id).Error

	if err != nil {
		return nil, err
	}

	return &post, nil
}
