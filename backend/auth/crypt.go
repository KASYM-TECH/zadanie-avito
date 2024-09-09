package auth

import (
	"avito/auth/config"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashedPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), config.PasswordCost)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func IsPasswordCorrect(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
