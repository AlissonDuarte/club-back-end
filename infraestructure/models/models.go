package models

// defined structs of models
// gorm.Model is a struct that contains the fields ID, CreatedAt, UpdatedAt, DeletedAt
// ID is the primary key
// CreatedAt is the field that will be filled with the creation date of the record
// UpdatedAt is the field that will be filled with the date of the last update of the record
// DeletedAt is the field that will be filled with the date of the deletion of the record
// The fields CreatedAt, UpdatedAt and DeletedAt are filled automatically by GORM
// The fields CreatedAt, UpdatedAt and DeletedAt are of the type time.Time
// The field DeletedAt is filled with the date of the deletion of the record and is not deleted from the database | soft delete

import (
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model 
	Name      		string 
	BirthDate 		string 
	PasswdHash    	string 
	Cep       		string 
	Email     		string `gorm:"unique"`
	Phone     		string `gorm:"unique"` 
}


func NewUser(name string, birthDate string, passwd string, cep string, email string, phone string) *User {

	return &User{
		Name: name,
		BirthDate: birthDate,
		PasswdHash: passwd,
		Cep: cep,
		Email: email,
		Phone: phone,
	}
}

func GeneratePasswordHash(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

func VerifyPassword(password, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
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

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	if err != nil {
		return err
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
	err := db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}	
