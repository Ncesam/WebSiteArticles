package types

import (
	authPb "backend/generated/proto/auth"
	configPb "backend/generated/proto/config"
	queuePb "backend/generated/proto/queue"
	grpcfactory "backend/pkg/grpcFactory"
)

type MapClients struct {
	Auth   *grpcfactory.Client[authPb.AuthServiceClient]
	Queue  *grpcfactory.Client[queuePb.QueueServiceClient]
	Config *grpcfactory.Client[configPb.ConfigServiceClient]
}
