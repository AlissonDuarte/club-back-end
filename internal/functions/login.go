package functions

import (
	"clube/infraestructure/models"
	"fmt"
	"time"

	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func FindUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	var user models.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

func VerifyPassword(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func LoginTotp(user *models.User) error {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Clube",
		AccountName: user.Username,
	})

	if err != nil {
		return fmt.Errorf("error to generate totp")
	}

	fmt.Println("SECRET KEY", key.Secret())
	now := time.Now()
	totpCode, err := totp.GenerateCode(key.Secret(), now)

	if err != nil {
		return fmt.Errorf("cannot access key")
	}

	valid := totp.Validate(totpCode, key.Secret())

	if valid {
		return nil
	} else {
		return fmt.Errorf("invalid code")
	}

}
