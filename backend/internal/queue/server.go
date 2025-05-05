package queue

import (
	queuepb "backend/generated/proto/queue"
	"backend/pkg/config"
	"backend/pkg/types"
	"context"

	"go.uber.org/zap"
)



type QueueServer struct {
	queuepb.UnimplementedQueueServiceServer
	
	logics *QueueLogics
	logger *zap.Logger
	cfg *config.Config
	ctx context.Context
}


func NewServer(logger *zap.Logger, cfg *config.Config, queueLogics *QueueLogics) *QueueServer{
	return &QueueServer{
		logics: queueLogics,
		logger: logger,
		cfg: cfg,
		ctx: context.Background(),
	}
}


func (s *QueueServer) StartConfig(ctx context.Context, req *queuepb.StartConfigRequest) (*queuepb.Empty, error) {
	task := types.InputForm{
		ConfigId: req.ConfigId,
		UserId: req.UserId,
		Prompt: req.Prompt,
		Data: req.Data,
	}
	s.logics.inputChannel<-task
	return nil, nil
}



