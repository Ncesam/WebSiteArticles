package validator

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func RegisterMyHandlers(logger *zap.Logger) {
	validator := validator.New()
	err := validator.RegisterValidation("password_strength", password_strength)
	if err != nil {
		logger.Error("Validator is broken, please check it")
		return
	}
	logger.Info("Validator successfully is ready")
}
