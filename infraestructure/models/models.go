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
)

type User struct {
	gorm.Model
	Name string
	BirthDate string
	Passwd string
	Cep string
	Email string
	Phone string
}

func NewUser(name string, birthDate string, passwd string, cep string, email string, phone string) *User {

	return &User{
		Name: name,
		BirthDate: birthDate,
		Passwd: passwd,
		Cep: cep,
		Email: email,
		Phone: phone,
	}
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