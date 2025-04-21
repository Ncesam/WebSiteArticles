package types

import (
	authPb "backend/internal/proto/auth"
	userPb "backend/internal/proto/user"
	grpcfactory "backend/pkg/grpcFactory"
)

type MapClients struct {
	User *grpcfactory.Client[userPb.UserServiceClient]
	Auth *grpcfactory.Client[authPb.AuthServiceClient]
}
