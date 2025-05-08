package routers

import (
	configPb "backend/generated/proto/config"
	"backend/internal/helpers"
	"backend/pkg/config"
	"backend/pkg/errors"
	jwt "backend/pkg/security/JWT"
	"backend/pkg/types"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetConfigs(logger *zap.Logger, cfg *config.Config, authcontroller *jwt.AuthController, clients *types.MapClients) gin.HandlerFunc {
	return func (c *gin.Context) {
		claims, _ := helpers.CheckUser(logger, c, authcontroller)
		id, _ := helpers.ExtractUserIDFromClaims(logger, c, claims)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()


		configs, err := clients.Config.Service.GetConfigs(ctx, &configPb.GetConfigsRequest{
			UserId: id,
		})
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "Get Configs failed");
		}
		c.JSON(200, configs)
	}
}

func AddConfig(logger *zap.Logger, cfg *config.Config, authcontroller *jwt.AuthController, clients *types.MapClients) gin.HandlerFunc {
	return func (c *gin.Context) {
		helpers.CheckUser(logger, c, authcontroller)
		var config types.ConfigForm
		err := c.ShouldBindJSON(&config)
		if err != nil {
			logger.Error("Failed to parse Body", zap.Error(err))
			c.AbortWithStatusJSON(errors.ErrBadRequest.Code, gin.H{"message": errors.ErrBadRequest.Message})
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		addConfigRequest := &configPb.AddConfigRequest{
			Name: config.Name,
			UserId: config.UserId,
			Prompt: config.Prompt,
			Delay: config.Delay,
			Email: config.Email,
			Password: config.Password,

		}
		_, err = clients.Config.Service.AddConfig(ctx, addConfigRequest)
		if err != nil {
			helpers.HandleGrpcError(logger, c, err, "Get Configs failed");
		}
		c.JSON(200, gin.H{"message": "successfully"})
	}
}