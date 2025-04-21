package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func RegisterMyHandlers(logger *zap.Logger) {
	if validatorInstanse, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validatorInstanse.RegisterValidation("password_strength", password_strength)
		logger.Info("Validator successfully is ready")
		return
	}
	logger.Error("Validator is broken, please check it")
}
