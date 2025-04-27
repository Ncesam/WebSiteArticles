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
	registerArticleRouters(router, logger, clients, cfg, authController)
	registerAccountRouters(router, logger, clients, cfg, authController)
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

func registerArticleRouters(router *gin.Engine, logger *zap.Logger, clients *types.MapClients, cfg *config.Config, authController *jwt.AuthController) {
	articleGroup := router.Group("/articles")
	articleGroup.POST("/", routers.AddArticleHandler(logger, cfg, clients, authController))
	articleGroup.GET("/", routers.ListArticlesHandler(logger, cfg, clients, authController))
	logger.Info("Article routes registered")
}

func registerAccountRouters(router *gin.Engine, logger *zap.Logger, clients *types.MapClients, cfg *config.Config, authController *jwt.AuthController) {
	accountGroup := router.Group("/account")
	accountGroup.POST("/", routers.AddAccountHandler(logger, cfg, clients, authController))
	accountGroup.GET("/", routers.ListAccountsHandler(logger, cfg, clients, authController))
	accountGroup.DELETE("/", routers.DeleteAccountHandler(logger, cfg, clients, authController))
	logger.Info("Account routes registered")
}

func registerQueueRouters(router *gin.Engine, logger *zap.Logger, clients *types.MapClients, cfg *config.Config, authController *jwt.AuthController) {
	queueGroup := router.Group("/queue")
	queueGroup.POST("/", routers.AddQueueHandler(logger, cfg, clients, authController))
	queueGroup.GET("/", routers.ListQueueHandler(logger, cfg, clients, authController))
	queueGroup.GET("/view", routers.ViewQueueHandler(logger, cfg, clients, authController))
	queueGroup.DELETE("/", routers.DeleteQueueHandler(logger, cfg, clients, authController))
	logger.Info("Queue routes registered")
}


