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
	return func(c *gin.Context) {
		logger.Debug("Handling GetConfigs request")
		claims, _ := helpers.CheckUser(logger, c, authcontroller)
		id, _ := helpers.ExtractUserIDFromClaims(logger, c, claims)
		logger.Debug("Extracted user ID", zap.Int64("user_id", id))

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		logger.Debug("Sending gRPC request to GetConfigs")
		configs, err := clients.Config.Service.GetConfigs(ctx, &configPb.GetConfigsRequest{
			UserId: id,
		})
		if err != nil {
			logger.Error("GetConfigs gRPC call failed", zap.Error(err))
			helpers.HandleGrpcError(logger, c, err, "Get Configs failed")
			return
		}
		logger.Debug("Successfully received configs", zap.Int("count", len(configs.Configs)))
		c.JSON(200, configs)
	}
}

func AddConfig(logger *zap.Logger, cfg *config.Config, authcontroller *jwt.AuthController, clients *types.MapClients) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug("Handling AddConfig request")
		helpers.CheckUser(logger, c, authcontroller)

		var config types.ConfigForm
		err := c.ShouldBindJSON(&config)
		if err != nil {
			logger.Error("Failed to parse body", zap.Error(err))
			c.AbortWithStatusJSON(errors.ErrBadRequest.Code, gin.H{"message": errors.ErrBadRequest.Message})
			return
		}
		logger.Debug("Parsed request body", zap.String("name", config.Name), zap.Int64("user_id", config.UserId))

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		addConfigRequest := &configPb.AddConfigRequest{
			Name:     config.Name,
			UserId:   config.UserId,
			Prompt:   config.Prompt,
			Delay:    config.Delay,
			Email:    config.Email,
			Password: config.Password,
		}

		logger.Debug("Sending gRPC request to AddConfig")
		_, err = clients.Config.Service.AddConfig(ctx, addConfigRequest)
		if err != nil {
			logger.Error("AddConfig gRPC call failed", zap.Error(err))
			helpers.HandleGrpcError(logger, c, err, "Add Config failed")
			return
		}
		logger.Debug("Config successfully added")
		c.JSON(200, gin.H{"message": "successfully"})
	}
}
