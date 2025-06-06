package security

import (
	"golang.org/x/crypto/bcrypt"

	"backend/pkg/errors"
)

func CompareHashAndPassword(password string, hashPassword string) error {
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
