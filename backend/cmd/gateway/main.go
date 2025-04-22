package main

import (
	_ "backend/docs"
	"backend/internal/gateway"
	"backend/pkg/config"
	"backend/pkg/logger"
	jwt "backend/pkg/security/JWT"
	"backend/pkg/validator"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// @title           Backend Gateway API
// @version         1.0
// @description     API-шлюз для взаимодействия с микросервисами пользователя.
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8080
// @BasePath  /

// @schemes http https
func main() {
	cfg, err := config.SetupConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	var logLevel zapcore.Level
	if cfg.APP.DEBUG {
		logLevel = zap.DebugLevel
	} else {
		logLevel = zap.InfoLevel
	}

	loggerInstanse, err := logger.SetupLogger(logger.Parameters{
		Level: logLevel,
	})
	if err != nil {
		log.Fatalf("Error setting up logger: %v", err)
	}

	validator.RegisterMyHandlers(loggerInstanse)

	app := gin.New()
	app.GET("/swagger/*any",ginSwagger.WrapHandler(swaggerFiles.Handler))

	app.Use(logger.MiddleWare(loggerInstanse, 30, time.Minute))
	clients, err := gateway.GetAllClients(loggerInstanse, cfg)
	if err != nil {
		loggerInstanse.Error("GRPC not started", zap.Error(err))
		return
	}
	authController := jwt.New(loggerInstanse, cfg)

	// TODO: setup routers
	gateway.RegisterAllRouters(app, loggerInstanse, clients, cfg, &authController)

	// TODO: run this server
	app.Run(fmt.Sprintf("%s:%d", cfg.APP.HOST, cfg.APP.PORT))

}
