package grpcfactory

import (
	"backend/pkg/config"
	"backend/pkg/errors"
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
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
	conn, err := grpc.NewClient(address, grpc.WithInsecure())
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
