package types

import (
	authPb "backend/generated/proto/auth"
	queuePb "backend/generated/proto/queue"
	userPb "backend/generated/proto/user"
	grpcfactory "backend/pkg/grpcFactory"
)

type MapClients struct {
	User  *grpcfactory.Client[userPb.UserServiceClient]
	Auth  *grpcfactory.Client[authPb.AuthServiceClient]
	Queue *grpcfactory.Client[queuePb.QueueServiceClient]
}
