package main

import (
	"backend/pkg/config"
	"backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

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
	app := gin.New()
	app.Use(logger.MiddleWare(loggerInstanse, 30, time.Minute))
	app.Run(":8080")

	// TODO: setup routers

	// TODO: run this server

}
