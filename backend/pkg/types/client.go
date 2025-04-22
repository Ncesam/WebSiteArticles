package types

import (
	authPb "backend/generated/proto/auth"
	userPb "backend/generated/proto/user"
	articlePb "backend/generated/proto/article"
	accountPb "backend/generated/proto/account"
	queuePb "backend/generated/proto/queue"
	grpcfactory "backend/pkg/grpcFactory"
)

type MapClients struct {
	User *grpcfactory.Client[userPb.UserServiceClient]
	Auth *grpcfactory.Client[authPb.AuthServiceClient]
	Article *grpcfactory.Client[articlePb.ArticleServiceClient]
	Account *grpcfactory.Client[accountPb.AccountServiceClient]
	Queue *grpcfactory.Client[queuePb.QueueServiceClient]
}
