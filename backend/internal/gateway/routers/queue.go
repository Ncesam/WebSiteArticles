package routers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	queuePb "backend/generated/proto/queue"
	"backend/internal/helpers"
	"backend/pkg/config"
	"backend/pkg/errors"
	jwt "backend/pkg/security/JWT"
	"backend/pkg/types"
)

// StartConfig godoc
// @Summary     Запуск конфигурации генерации
// @Description Принимает JSON с ConfigId, Prompt и Data и инициирует процесс генерации текста
// @Tags        config
// @Accept      json
// @Produce     json
// @Param       Authorization header string true "Bearer токен"
// @Param       body body types.InputForm true "Данные для генерации"
// @Success     200 {object} map[string]string "Successfully"
// @Failure     400 {object} map[string]string "Bad Request"
// @Failure     401 {object} map[string]string "Invalid Credentials"
// @Router      /queue/ [post]
func StartConfig(logger *zap.Logger, cfg *config.Config, clients *types.MapClients, authController *jwt.AuthController) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug("Handling StartConfig request")

		var body types.InputForm
		if err := c.ShouldBindJSON(&body); err != nil {
			logger.Error("Invalid request body", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errors.ErrBadRequest.Message})
			return
		}
		logger.Debug("Parsed input body", zap.String("config_id", body.ConfigId))

		claims, ok := helpers.CheckUser(logger, c, authController)
		if !ok {
			logger.Error("User authentication failed", zap.String("error", errors.ErrInvalidCredentials.Message))
			c.AbortWithStatusJSON(errors.ErrInvalidCredentials.Code, gin.H{"message": errors.ErrInvalidCredentials.Message})
			return
		}

		userId, ok := helpers.ExtractUserIDFromClaims(logger, c, claims)
		if !ok {
			logger.Error("Failed to extract user ID from claims", zap.String("error", errors.ErrInvalidCredentials.Message))
			c.AbortWithStatusJSON(errors.ErrInvalidCredentials.Code, gin.H{"message": errors.ErrInvalidCredentials.Message})
			return
		}
		logger.Debug("Authenticated user", zap.Int64("user_id", userId))

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		logger.Debug("Sending gRPC StartConfig request", zap.String("config_id", body.ConfigId), zap.Int64("user_id", userId))
		_, err := clients.Queue.Service.StartConfig(ctx, &queuePb.StartConfigRequest{
			ConfigId: body.ConfigId,
			UserId:   userId,
			Prompt:   body.Prompt,
			Data:     body.Data,
		})
		if err != nil {
			logger.Error("StartConfig gRPC call failed", zap.Error(err))
			helpers.HandleGrpcError(logger, c, err, "Fail to start config")
			return
		}

		logger.Debug("Config started successfully")
		c.JSON(http.StatusOK, gin.H{"message": "Successfully"})
	}
}
