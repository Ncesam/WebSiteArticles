package security

import (
	"backend/pkg/config"
	"backend/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func CompareHashAndPassword(password string, hashPassword string, cfg *config.Config) error {
	// Проверяем, не пустые ли пароли
	if password == "" || hashPassword == "" {
		return errors.ErrInvalidCredentials
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		return err
	}

	return nil
}

func GenerateHashedPassword(password string) (string, error) {
	if password == "" {
		return "", errors.ErrInvalidCredentials
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errors.ErrInternalServer
	}
	return string(hashedPassword), nil
}
