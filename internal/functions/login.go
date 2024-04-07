package functions

import (
	"fmt"
	"gorm.io/gorm"
	"clube/infraestructure/models"
	"golang.org/x/crypto/bcrypt"
)
func FindUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	var user models.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("User not found")
		}
		return nil, result.Error
	}
	return &user, nil
}


func VerifyPassword(password, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}