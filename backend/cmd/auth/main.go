package auth

import (
	authPb "backend/generated/proto/auth"
	authServer "backend/internal/auth"
	"backend/internal/database/postgres"
	"backend/pkg/config"
	"backend/pkg/logger"
	jwt "backend/pkg/security/JWT"
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

	authController := jwt.New(loggerInstanse, cfg)
	db := postgres.Connect(cfg, loggerInstanse)

	server := grpc.NewServer()
	authPb.RegisterAuthServiceServer(server, authServer.NewAuthServer(db, cfg, loggerInstanse, &authController))

	listener, err := net.Listen("tcp", cfg.AUTH_SERVICE.ADDRESS)
	if err != nil {
		loggerInstanse.Fatal("Failed to listen",
			zap.Error(err),
			zap.String("address", cfg.AUTH_SERVICE.ADDRESS),
		)
	}
	server.Serve(listener)

}
