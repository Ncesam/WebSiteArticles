package gateway

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"backend/internal/gateway/routers"
	"backend/pkg/config"
	jwt "backend/pkg/security/JWT"
	"backend/pkg/types"
)

func RegisterAllRouters(router *gin.Engine, logger *zap.Logger, clients *types.MapClients, cfg *config.Config, authController *jwt.AuthController) {
	registerAuthRouters(router, logger, clients, cfg, authController)
	registerQueueRouters(router, logger, clients, cfg, authController)
	registerConfigRouters(router, logger, clients, cfg, authController)
}

func registerAuthRouters(router *gin.Engine, logger *zap.Logger, clients *types.MapClients, cfg *config.Config, authController *jwt.AuthController) {
	authGroup := router.Group("/api/auth")
	authGroup.POST("/login", routers.Login(logger, cfg, clients))
	authGroup.POST("/register", routers.Register(logger, cfg, clients))
	authGroup.POST("/refresh", routers.Refresh(logger, cfg, clients, authController))
	authGroup.PUT("/me", routers.Me(logger, cfg, clients, authController))
	authGroup.POST("/logout", routers.Logout(logger, cfg))
	logger.Info("Auth routes registered")
}
func registerConfigRouters(router *gin.Engine, logger *zap.Logger, clients *types.MapClients, cfg *config.Config, authController *jwt.AuthController) {
	configGroup := router.Group("/api/config")
	configGroup.GET("/", routers.GetConfigs(logger, cfg, authController, clients))
	configGroup.POST("/", routers.AddConfig(logger, cfg, authController, clients))
	logger.Info("Config routes registered")
}

func registerQueueRouters(router *gin.Engine, logger *zap.Logger, clients *types.MapClients, cfg *config.Config, authController *jwt.AuthController) {
	queueGroup := router.Group("/api/queue")
	queueGroup.POST("/", routers.StartConfig(logger, cfg, clients, authController))
	logger.Info("Queue routes registered")
}
