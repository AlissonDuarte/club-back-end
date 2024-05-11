package functions

import (
	"clube/infraestructure/models"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/pbkdf2"
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

func VerifyPassword(password, hashedPassword string) error {
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 2 {
		return errors.New("invalid hashed password format")
	}
	saltString := parts[0]
	hash := parts[1]

	salt, err := base64.StdEncoding.DecodeString(saltString)
	if err != nil {
		return err
	}
	hashedPasswordInput := pbkdf2.Key([]byte(password), salt, 10000, sha256.Size, sha256.New)
	hashInput := base64.StdEncoding.EncodeToString(hashedPasswordInput)
	fmt.Println("Hash da senha digitada", hashInput)
	fmt.Println("Hash da senha do banco", hash)

	if len(hash) != len(hashInput) {
		return errors.New("incorrect password, lenght diff")
	}

	var diff byte

	for i := 0; i < len(hash); i++ {
		diff |= hash[i] ^ hashInput[i]
	}

	if diff != 0 {
		return errors.New("incorrect password, diff exists")
	}

	if hash != hashInput {
		return errors.New("incorrect password")
	}

	return nil
}
