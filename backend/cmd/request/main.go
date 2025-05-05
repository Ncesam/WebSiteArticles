package main

import (
	"log"
	"net"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	configPb "backend/generated/proto/config"
	requestPb "backend/generated/proto/request"
	"backend/internal/request"
	"backend/pkg/config"
	grpcfactory "backend/pkg/grpcFactory"
	"backend/pkg/logger"
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
	defer loggerInstanse.Sync()

	grpcServer := grpc.NewServer()

	configClient, err := grpcfactory.NewClient(cfg.CONFIG_SERVICE.ADDRESS, loggerInstanse, cfg, configPb.NewConfigServiceClient)
	if err != nil {
		loggerInstanse.Error("Can't connect to config service", zap.Error(err))
		return
	}
	requestServer := request.NewRequestServer(loggerInstanse, cfg, configClient)

	requestPb.RegisterRequestServiceServer(grpcServer, requestServer)

	listener, err := net.Listen("tcp", cfg.REQUEST_SERVICE.ADDRESS)
	if err != nil {
		loggerInstanse.Fatal("Failed to listen",
			zap.Error(err),
			zap.String("address", cfg.REQUEST_SERVICE.ADDRESS),
		)
	}
	grpcServer.Serve(listener)
}
