package queue

import (
	queuepb "backend/generated/proto/queue"
	"backend/pkg/config"
	"backend/pkg/types"
	"context"
	"fmt"

	"go.uber.org/zap"
)

type QueueServer struct {
	queuepb.UnimplementedQueueServiceServer

	logics *QueueLogics
	logger *zap.Logger
	cfg    *config.Config
	ctx    context.Context
}

func NewServer(logger *zap.Logger, cfg *config.Config, queueLogics *QueueLogics) *QueueServer {
	logger.Info("Creating a new QueueServer instance")
	return &QueueServer{
		logics: queueLogics,
		logger: logger,
		cfg:    cfg,
		ctx:    context.Background(),
	}
}

func (s *QueueServer) StartConfig(ctx context.Context, req *queuepb.StartConfigRequest) (*queuepb.Empty, error) {
	s.logger.Debug("Received StartConfig request", zap.Int64("config_id", req.ConfigId), zap.Int64("user_id", req.UserId))

	task := types.InputForm{
		ConfigId: req.ConfigId,
		UserId:   req.UserId,
		Prompt:   req.Prompt,
		Data:     req.Data,
	}

	select {
	case s.logics.inputChannel <- task:
		s.logger.Info("Task successfully added to the input channel", zap.Int64("config_id", req.ConfigId))
	default:
		s.logger.Warn("Failed to add task to input channel (channel might be full)", zap.Int64("config_id", req.ConfigId))
		return nil, fmt.Errorf("failed to add task to input channel")
	}

	return &queuepb.Empty{}, nil
}
