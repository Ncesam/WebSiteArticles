package gateway

import (
	authPb "backend/internal/proto/auth"
	userPb "backend/internal/proto/user"
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
	return &types.MapClients{
		User: user,
		Auth: auth,
	}, nil
}
