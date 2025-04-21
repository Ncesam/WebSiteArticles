package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func password_strength(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	hasUpper := regexp.MustCompile(`[A-Z]`)
	hasLower := regexp.MustCompile(`[a-z]`)
	return hasLower.MatchString(value) && hasUpper.MatchString(value) && len(value) >= 8
}
