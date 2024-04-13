package functions

import (
	"regexp"
)

func EmailCheck(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func PhoneCheck(phone string) bool {
	phoneRegex := regexp.MustCompile(`^[0-9]{11}$`)
	return phoneRegex.MatchString(phone)
}

func ValidGender(gender string) bool {
	validGenders := []string{"m", "f", "o"}
	for _, opt := range validGenders {
		if opt == gender {
			return true
		}
	}
	return false
}
