package grpcfactory

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"backend/pkg/config"
	"backend/pkg/errors"
)

type GRPCClientConstructor[T any] func(conn grpc.ClientConnInterface) T

type Client[T any] struct {
	Service    T
	Connection *grpc.ClientConn
	Ctx        context.Context
	Cancel     context.CancelFunc
	Logger     *zap.Logger
	Cfg        *config.Config
}


func NewClient[T any](address string, logger *zap.Logger, cfg *config.Config, constructor GRPCClientConstructor[T]) (*Client[T], error) {
	serviceConfig := `{
		"methodConfig": [{
		  "retryPolicy": {
			"maxAttempts": 4,
			"initialBackoff": "0.5s",
			"maxBackoff": "5s",
			"backoffMultiplier": 2.0,
			"retryableStatusCodes": ["UNAVAILABLE"]
		  }
		}]
	}`
	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultServiceConfig(serviceConfig),
		grpc.WithBackoffMaxDelay(5*time.Second),
	)
	if err != nil {
		logger.Error(errors.ErrServiceUnavailable.Message, zap.Error(errors.ErrServiceUnavailable.Err))
		return nil, errors.ErrServiceUnavailable
	}
	client := constructor(conn)

	logger.Info("Successfully connected to Service")

	return &Client[T]{
		Logger:     logger,
		Cfg:        cfg,
		Connection: conn,
		Service:    client,
	}, nil

}

func (c *Client[T]) Close() {
	c.Cancel()
	c.Connection.Close()
}
