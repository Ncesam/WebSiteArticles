package gateway

import (
	"backend/internal/gateway/routers"
	"backend/pkg/config"
	jwt "backend/pkg/security/JWT"
	"backend/pkg/types"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterAllRouters(router *gin.Engine, logger *zap.Logger, clients *types.MapClients, cfg *config.Config, authController *jwt.AuthController) {
	registerUserRouters(router, logger, clients, cfg, authController)
	registerAuthRouters(router, logger, clients, cfg, authController)
	registerQueueRouters(router, logger, clients, cfg, authController)
}

func registerUserRouters(router *gin.Engine, logger *zap.Logger, clients *types.MapClients, cfg *config.Config, authController *jwt.AuthController) {
	userGroup := router.Group("/user")
	userGroup.POST("/", routers.AddUserHandler(logger, cfg, clients, authController))
	userGroup.GET("/", routers.GetUserHandler(logger, cfg, clients, authController))
	userGroup.PUT("/", routers.UpdateUserHandler(logger, cfg, clients, authController))
	userGroup.DELETE("/", routers.DeleteUserHandler(logger, cfg, clients, authController))
	logger.Info("User routes registered")
}
func registerAuthRouters(router *gin.Engine, logger *zap.Logger, clients *types.MapClients, cfg *config.Config, authController *jwt.AuthController) {
	authGroup := router.Group("/auth")
	authGroup.POST("/login", routers.Login(logger, cfg, clients))
	authGroup.POST("/register", routers.Register(logger, cfg, clients))
	authGroup.POST("/refresh", routers.Refresh(logger, cfg, clients, authController))
	authGroup.POST("/me", routers.Me(logger, cfg, clients, authController))
	authGroup.POST("/logout", routers.Logout(logger, cfg))
	logger.Info("Auth routes registered")
}

func registerQueueRouters(router *gin.Engine, logger *zap.Logger, clients *types.MapClients, cfg *config.Config, authController *jwt.AuthController) {
	queueGroup := router.Group("/queue")
	queueGroup.POST("/", routers.StartConfig(logger, cfg, clients, authController))
	logger.Info("Queue routes registered")
}
