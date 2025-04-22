package gateway

import (
	accountPb "backend/generated/proto/account"
	articlePb "backend/generated/proto/article"
	authPb "backend/generated/proto/auth"
	userPb "backend/generated/proto/user"
	queuePb "backend/generated/proto/queue"
	"backend/pkg/config"
	grpcfactory "backend/pkg/grpcFactory"
	"backend/pkg/types"

	"go.uber.org/zap"
)

func GetAllClients(logger *zap.Logger, cfg *config.Config) (*types.MapClients, error) {
	user, err := grpcfactory.NewClient[userPb.UserServiceClient]("user_service:50001", logger, cfg, userPb.NewUserServiceClient)
	if err != nil {
		logger.Error("User Service isn't ready", zap.Error(err))
		return nil, err
	}
	auth, err := grpcfactory.NewClient[authPb.AuthServiceClient]("auth_service:50002", logger, cfg, authPb.NewAuthServiceClient)
	if err != nil {
		logger.Error("Auth Service isn't ready", zap.Error(err))
		return nil, err
	}
	article, err := grpcfactory.NewClient[articlePb.ArticleServiceClient]("article_service:50003", logger, cfg, articlePb.NewArticleServiceClient)
	if err != nil {
		logger.Error("Article Service isn't ready", zap.Error(err))
		return nil, err
	}
	account, err := grpcfactory.NewClient[accountPb.AccountServiceClient]("account_service:50004", logger, cfg, accountPb.NewAccountServiceClient)
	if err != nil {
		logger.Error("Account Service isn't ready", zap.Error(err))
		return nil, err
	}
	queue, err := grpcfactory.NewClient[queuePb.QueueServiceClient]("queue_service:50005", logger, cfg, queuePb.NewQueueServiceClient)
	if err != nil {
		logger.Error("Queue Service isn't ready", zap.Error(err))
		return nil, err
	}
	return &types.MapClients{
		User:    user,
		Auth:    auth,
		Article: article,
		Account: account,
		Queue: queue,
	}, nil
}
