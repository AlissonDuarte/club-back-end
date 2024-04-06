package models


import (
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model 
	Name      		string 
	Username 		string `gorm:"unique"`
	Gender 			string
	BirthDate 		string 
	PasswdHash    	string 
	Cep       		string 
	Email     		string `gorm:"unique"`
	Phone     		string `gorm:"unique"` 
}


func NewUser(name string, username string, gender string, birthDate string, passwd string, cep string, email string, phone string) *User {

	return &User{
		Name: name,
		Username: username,
		Gender: gender,
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
