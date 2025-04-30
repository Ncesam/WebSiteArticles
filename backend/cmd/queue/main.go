package main

import (
	queuepb "backend/generated/proto/queue"
	"backend/internal/queue"
	"backend/pkg/config"
	grpcfactory "backend/pkg/grpcFactory"
	"backend/pkg/logger"
	"log"
	"net"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
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
	logicController := queueLogics.New(logger, cfg, queue.RabbitClient)
	rabbitClient, err:= queue.NewQueueClient()
	if err != nil {
		loggerInstanse.Error("Client not starting", zap.Error(err))
		panic(err)
	}
	requestClient, err := grpcfactory.NewClient[requestPb.RequestServiceClient]("request_service:50001", client.logger, client.cfg, requestPb.NewRequestServiceClient)
	if err != nil {
		client.logger.Error("User Service isn't ready", zap.Error(err))
	}
	queuepb.RegisterQueueServiceServer(grpcServer, queueServer)

	listener, err := net.Listen("tcp", cfg.QUEUE_SERVICE.ADDRESS)
	if err != nil {
		loggerInstanse.Fatal("Failed to listen",
			zap.Error(err),
			zap.String("address", cfg.QUEUE_SERVICE.ADDRESS),
		)
	}
	grpcServer.Serve(listener)
}