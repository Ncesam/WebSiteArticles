package main

import (
	configPb "backend/generated/proto/config"
	"backend/internal/config"
	configPkg "backend/pkg/config"
	"backend/internal/database/mongo"
	"backend/pkg/logger"
	"log"
	"net"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)


func main() {
	cfg, err := configPkg.SetupConfig()
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
	defer loggerInstanse.Sync()

	grpcServer := grpc.NewServer();

	mongoDataBase, err := mongo.Connect(loggerInstanse, cfg)
	if err != nil {
		return
	}

	configServer := config.NewConfigServer(loggerInstanse, cfg, mongoDataBase)

	configPb.RegisterConfigServiceServer(grpcServer, configServer)

	listener, err := net.Listen("tcp", cfg.CONFIG_SERVICE.ADDRESS)
	if err != nil {
		loggerInstanse.Error("Fail to open connection on Config Service", zap.Error(err))
		return
	}
	grpcServer.Serve(listener)
}