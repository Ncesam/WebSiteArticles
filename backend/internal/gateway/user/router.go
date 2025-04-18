package user

import (
	"backend/pkg/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetUserHandler(logger *zap.Logger, cfg *config.Config, client) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug("Request for GetUser")

	}

}
