package queue

import (
	"backend/generated/proto/config"
	"backend/generated/proto/request"
	"backend/pkg/config"
	grpcfactory "backend/pkg/grpcFactory"
	"context"
	"time"

	"go.uber.org/zap"
)


type QueueLogics struct {
	
	logger *zap.Logger
	cfg *config.Config
	ctx context.Context
	configClient *grpcfactory.Client[configPb.ConfigServiceClient]
	requestClient *grpcfactory.Client[request.RequestServiceClient]
}


func NewLogics(logger *zap.Logger, cfg *config.Config, configClient *grpcfactory.Client[configPb.ConfigServiceClient], requestClient *grpcfactory.Client[request.RequestServiceClient]) *QueueLogics {
	logger.Info("Create a QueueLogics")
	ctx, _ := context.WithTimeout(context.Background(), time.Minute * 10)
	return &QueueLogics{
		logger: logger,
		cfg: cfg,
		ctx: ctx,
		configClient: configClient,
		requestClient: requestClient,
	}
}

func (logics *QueueLogics) start( )