package gateway

import (
	authPb "backend/generated/proto/auth"
	queuePb "backend/generated/proto/queue"
	userPb "backend/generated/proto/user"

	"backend/pkg/config"
	grpcfactory "backend/pkg/grpcFactory"
	"backend/pkg/types"

	"go.uber.org/zap"
)

func GetAllClients(logger *zap.Logger, cfg *config.Config) (*types.MapClients, error) {
	user, err := grpcfactory.NewClient[userPb.UserServiceClient]("user_service:50001", logger, cfg, userPb.NewUserServiceClient)
	if err != nil {
		logger.Error("User Service isn't ready", zap.Error(err))
		return nil, err
	}
	auth, err := grpcfactory.NewClient[authPb.AuthServiceClient]("auth_service:50002", logger, cfg, authPb.NewAuthServiceClient)
	if err != nil {
		logger.Error("Auth Service isn't ready", zap.Error(err))
		return nil, err
	}
	queue, err := grpcfactory.NewClient[queuePb.QueueServiceClient]("queue_service:50005", logger, cfg, queuePb.NewQueueServiceClient)
	if err != nil {
		logger.Error("Queue Service isn't ready", zap.Error(err))
		return nil, err
	}

	return &types.MapClients{
		User:  user,
		Auth:  auth,
		Queue: queue,
	}, nil
}
