package gateway

import (
	accountPb "backend/generated/proto/account"
	articlePb "backend/generated/proto/article"
	authPb "backend/generated/proto/auth"
	userPb "backend/generated/proto/user"
	queuePb "backend/generated/proto/queue"
	themePb "backend/generated/proto/theme"
	promptPb "backend/generated/proto/prompt"
	generatorPb "backend/generated/proto/generator"
	partnerPb "backend/generated/proto/partner"

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
	theme, err := grpcfactory.NewClient[themePb.ThemeServiceClient]("theme_service:50006", logger, cfg, themePb.NewThemeServiceClient)
	if err != nil {
		logger.Error("Theme Service isn't ready", zap.Error(err))
		return nil, err
	}
	prompt, err := grpcfactory.NewClient[promptPb.PromptServiceClient]("prompt_service:50007", logger, cfg, promptPb.NewPromptServiceClient)
	if err != nil {
		logger.Error("Prompt Service isn't ready", zap.Error(err))
		return nil, err
	}
	generator, err := grpcfactory.NewClient[generatorPb.GeneratorServiceClient]("generator_service:50008", logger, cfg, generatorPb.NewGeneratorServiceClient)
	if err != nil {
		logger.Error("Generator Service isn't ready", zap.Error(err))
		return nil, err
	}
	partner, err := grpcfactory.NewClient[partnerPb.PartnerServiceClient]("partner_service:50009", logger, cfg, partnerPb.NewPartnerServiceClient)
	if err != nil {
		logger.Error("Partner Service isn't ready", zap.Error(err))
		return nil, err
	}

	return &types.MapClients{
		User:      user,
		Auth:      auth,
		Article:   article,
		Account:   account,
		Queue:     queue,
		Theme:     theme,
		Prompt:    prompt,
		Generator: generator,
		Partner:   partner,
	}, nil
}
