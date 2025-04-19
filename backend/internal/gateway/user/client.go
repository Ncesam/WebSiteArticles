package user

import (
	userPb "backend/internal/proto/user"
	"backend/pkg/config"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)




type UserClient struct {
	logger *zap.Logger
	grpcConn *grpc.ClientConn
	userService *userPb.UserServiceClient
	cxt context.Context
	cfg *config.Config
}
func New(logger *zap.Logger, cfg *config.Config) UserClient {
	logger.Info("Connection to the gRPC User server")
	conn, err := grpc.NewClient("user_service:50001", grpc.WithInsecure())
	if err != nil {
		logger.Error(fmt.Sprintf("Connection to User Service raise error: ",err))
	}
	userClient := userPb.NewUserServiceClient(conn)

	logger.Info("Connection to the User Service successful")

	return UserClient{
		logger: logger,
		cfg: cfg,
		grpcConn: conn,
		userService: &userClient,
		cxt: context.WithTimeout(context.Background(), 2 * time.Minute),
	}


}
