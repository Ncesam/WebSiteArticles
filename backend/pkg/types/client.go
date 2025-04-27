package types

import (
	accountPb "backend/generated/proto/account"
	articlePb "backend/generated/proto/article"
	authPb "backend/generated/proto/auth"
	queuePb "backend/generated/proto/queue"
	themePb "backend/generated/proto/theme"
	promptPb "backend/generated/proto/prompt"
	generatorPb "backend/generated/proto/generator"
	partnerPb "backend/generated/proto/partner"
	userPb "backend/generated/proto/user"
	grpcfactory "backend/pkg/grpcFactory"
)

type MapClients struct {
	User      *grpcfactory.Client[userPb.UserServiceClient]
	Auth      *grpcfactory.Client[authPb.AuthServiceClient]
	Article   *grpcfactory.Client[articlePb.ArticleServiceClient]
	Account   *grpcfactory.Client[accountPb.AccountServiceClient]
	Queue     *grpcfactory.Client[queuePb.QueueServiceClient]
	Theme     *grpcfactory.Client[themePb.ThemeServiceClient]
	Prompt    *grpcfactory.Client[promptPb.PromptServiceClient]
	Generator *grpcfactory.Client[generatorPb.GeneratorServiceClient]
	Partner   *grpcfactory.Client[partnerPb.PartnerServiceClient]
}
