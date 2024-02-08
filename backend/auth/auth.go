package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func GetPassword(password string) (string, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return "", err
	}

	return hashedPassword, nil
}
