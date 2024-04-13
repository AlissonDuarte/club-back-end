package models

import (
	"errors"

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
