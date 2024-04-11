package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name       string
	Username   string `gorm:"unique"`
	Gender     string
	BirthDate  string
	PasswdHash string
	Email      string  `gorm:"unique"`
	Phone      string  `gorm:"unique"`
	Clubs      []*Club `gorm:"many2many:user_club;"`
	ClubOnwer  []*Club `gorm:"many2many:owner_club;"`
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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

func (u *User) Save(db *gorm.DB) error {
	return db.Create(u).Error
}

func (u *User) Update(db *gorm.DB) error {
	return db.Save(u).Error
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
	).Preload("Clubs", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "name", "created_at")
	}).First(&user, id).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}
