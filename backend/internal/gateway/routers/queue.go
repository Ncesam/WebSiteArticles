package routers

import (
	queuePb "backend/generated/proto/queue"
	"backend/internal/helpers"
	"backend/pkg/config"
	"backend/pkg/errors"
	"backend/pkg/types"
	"context"
	"net/http"

	jwt "backend/pkg/security/JWT"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
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
	return func (c *gin.Context)  {
		var body types.InputForm;
		if err := c.ShouldBindJSON(&body); err != nil {
			logger.Error("Invalid Body", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errors.ErrBadRequest.Message})
			return
		}
		claims, ok := helpers.CheckUser(logger, c, authController)
		if !ok {
			logger.Error(errors.ErrInvalidCredentials.Message)
			c.AbortWithStatusJSON(errors.ErrInvalidCredentials.Code, gin.H{"message": errors.ErrInvalidCredentials.Message})
			return
		}
		userId, ok := helpers.ExtractUserIDFromClaims(logger, c, claims)
		if !ok {
			logger.Error(errors.ErrInvalidCredentials.Message)
			c.AbortWithStatusJSON(errors.ErrInvalidCredentials.Code, gin.H{"message": errors.ErrInvalidCredentials.Message})
			return
		}
		clients.Queue.Service.StartConfig(context.Background(), &queuePb.StartConfigRequest{
			ConfigId: body.ConfigId,
			UserId: userId,
			Prompt: body.Prompt,
			Data: body.Data,
		})
		c.JSON(http.StatusOK, gin.H{"message": "Successfully"})
	}
}
