package gateway

import (
	authPb "backend/generated/proto/auth"
	configPb "backend/generated/proto/config"
	queuePb "backend/generated/proto/queue"

	"backend/pkg/config"
	grpcfactory "backend/pkg/grpcFactory"
	"backend/pkg/types"

	"go.uber.org/zap"
)

func GetAllClients(logger *zap.Logger, cfg *config.Config) (*types.MapClients, error) {
	auth, err := grpcfactory.NewClient[authPb.AuthServiceClient](cfg.AUTH_SERVICE.ADDRESS, logger, cfg, authPb.NewAuthServiceClient)
	if err != nil {
		logger.Error("Auth Service isn't ready", zap.Error(err))
		return nil, err
	}
	queue, err := grpcfactory.NewClient[queuePb.QueueServiceClient](cfg.QUEUE_SERVICE.ADDRESS, logger, cfg, queuePb.NewQueueServiceClient)
	if err != nil {
		logger.Error("Queue Service isn't ready", zap.Error(err))
		return nil, err
	}
	config, err := grpcfactory.NewClient[configPb.ConfigServiceClient](cfg.CONFIG_SERVICE.ADDRESS, logger, cfg, configPb.NewConfigServiceClient)
	if err != nil {
		logger.Error("Config Service isn't ready", zap.Error(err))
		return nil, err
	}

	return &types.MapClients{
		Auth:  auth,
		Queue: queue,
		Config: config,
	}, nil
}
