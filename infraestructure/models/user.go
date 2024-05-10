package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name             string
	Username         string `gorm:"unique"`
	Gender           string
	BirthDate        string
	PasswdHash       string
	Email            string      `gorm:"unique"`
	Phone            string      `gorm:"unique"`
	Bio              string      `gorm:"default:null"`
	Clubs            []*Club     `gorm:"many2many:user_club;constraint:OnDelete:CASCADE"`
	ClubOnwer        []*Club     `gorm:"many2many:owner_club;constraint:OnDelete:CASCADE"`
	Posts            []Post      `gorm:"constraint:OnDelete:CASCADE"` // Relacionamento adicionado com a opção OnDelete("CASCADE")
	ProfilePictureID uint        `gorm:"default:null"`
	ProfilePicture   *UserUpload `gorm:"foreignKey:ProfilePictureID"`
	Followers        []*User     `gorm:"many2many:user_followers; constraint:OnDelete:Cascade"`
	Following        []*User     `gorm:"many2many:user_following; constraint:OnDelete:Cascade"`
}

func NewUser(name string, username string, gender string, birthDate string, passwd string, email string, phone string) *User {

	return &User{
		Name:       name,
		Username:   username,
		Gender:     gender,
		BirthDate:  birthDate,
		PasswdHash: passwd,
		Email:      email,
		Phone:      phone,
	}
}

func GeneratePasswordHash(password string) (string, error) {

	if len(password) < 8 || len(password) > 72 {
		return "", errors.New("password must be between 8 and 72 characters")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if u.PasswdHash != "" {
		hashedPassword, err := GeneratePasswordHash(u.PasswdHash)
		if err != nil {
			return err
		}
		u.PasswdHash = hashedPassword
	}
	return nil
}

func (u *User) Save(db *gorm.DB) (uint, error) {

	var existingUser User
	err := db.Where("username = ?", u.Username).First(&existingUser).Error
	if err == nil {
		return 0, errors.New("this username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {

		return 0, err
	}

	err = db.Where("email = ?", u.Email).First(&existingUser).Error
	if err == nil {
		return 0, errors.New("this email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {

		return 0, err
	}

	err = db.Where("phone = ?", u.Phone).First(&existingUser).Error
	if err == nil {
		return 0, errors.New("this phone number already exists")

	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}

	err = db.Create(u).Error
	if err != nil {
		return 0, err
	}

	return u.ID, nil
}

func (u *User) ChangePassword(db *gorm.DB, newPassword string) error {
	var user User

	if err := db.Where("id = ?", u.ID).First(&user).Error; err != nil {
		return err
	}

	// Salvar as alterações no banco de dados
	if err := db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (u *User) Update(db *gorm.DB, newPassword string) error {
	var existingUser User
	err := db.Where("username = ? AND id != ?", u.Username, u.ID).First(&existingUser).Error
	if err == nil {
		return errors.New("this username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	err = db.Select("passwd_hash").First(&existingUser, u.ID).Error
	if err != nil {
		return err
	}

	err = db.Where("email = ? AND id != ?", u.Email, u.ID).First(&existingUser).Error
	if err == nil {
		return errors.New("this email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	err = db.Where("phone = ? AND id != ?", u.Phone, u.ID).First(&existingUser).Error
	if err == nil {
		return errors.New("this phone number already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Tratamento da senha
	if newPassword == "" {
		u.PasswdHash = existingUser.PasswdHash
	}

	// Atualizar os dados do usuário
	err = db.Save(u).Error
	if err != nil {
		return err
	}

	return nil
}

func UserGetById(db *gorm.DB, id int) (*User, error) {
	var user User
	err := db.Select("id",
		"email",
		"phone",
		"username",
		"name",
		"gender",
		"created_at",
		"birth_date",
		"bio",
		"profile_picture_id",
	).Preload("Clubs", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "name", "created_at")
	}).First(&user, id).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func Follow(db *gorm.DB, userID, followedID uint) error {
	var user User
	err := db.First(&user, userID).Error
	if err != nil {
		return err
	}

	var followed User
	err = db.First(&followed, followedID).Error
	if err != nil {
		return err
	}

	err = db.Model(&user).Association("Following").Append(&followed)
	if err != nil {
		return err
	}

	err = db.Model(&followed).Association("Followers").Append(&user)
	if err != nil {
		return err
	}

	return nil
}

// Unfollow remove a relação de seguidor entre dois usuários

func Unfollow(db *gorm.DB, userID, followedID uint) error {
	var user User
	err := db.First(&user, userID).Error
	if err != nil {
		return err
	}

	var followed User
	err = db.First(&followed, followedID).Error
	if err != nil {
		return err
	}

	err = db.Model(&user).Association("Following").Delete(&followed)
	if err != nil {
		return err
	}

	err = db.Model(&followed).Association("Followers").Delete(&user)
	if err != nil {
		return err
	}

	return nil
}

// GetFollowers retorna todos os seguidores de um usuário

func GetFollowers(db *gorm.DB, userID uint) ([]User, error) {
	var user User
	err := db.Preload("Followers", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "username", "profile_picture_id").Preload("ProfilePicture", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "file_path")
		})
	}).First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	followers := make([]User, len(user.Followers))
	for i, follower := range user.Followers {
		followers[i] = *follower
	}

	return followers, nil
}

// GetFollowing retorna todos os usuários que um usuário segue

func GetFollowing(db *gorm.DB, userID uint) ([]User, error) {
	var user User
	err := db.Preload("Following", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "username", "profile_picture_id").Preload("ProfilePicture", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "file_path")
		})
	}).First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	following := make([]User, len(user.Following))
	for i, followed := range user.Following {
		following[i] = *followed
	}

	return following, nil
}

// GetFeed retorna todos os posts dos usuários que um usuário segue
func GetFeed(db *gorm.DB, userID uint, offset, limit int) ([]Post, error) {
	var user User
	err := db.Preload("Following", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id")
	}).First(&user, userID).Error
	if err != nil {
		return nil, err
	}

	// Extrair os IDs dos usuários seguidos
	var followingIDs []uint
	for _, following := range user.Following {
		followingIDs = append(followingIDs, following.ID)
	}

	var posts []Post
	err = db.Preload("User", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "name", "username", "profile_picture_id").Preload("ProfilePicture", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "file_path")
		})
	}).Preload("Image").Where("user_id IN (?) AND (club_id IS NULL OR club_id = 0)", followingIDs).Offset(offset).Limit(limit).Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}
