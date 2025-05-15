package config

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	configPb "backend/generated/proto/config"
	"backend/internal/database/mongo"
	"backend/pkg/config"
)

type ConfigServer struct {
	configPb.UnimplementedConfigServiceServer
	logger *zap.Logger
	cfg    *config.Config
	ctx    context.Context
	db     *mongo.MongoDatabase
}

func NewConfigServer(logger *zap.Logger, cfg *config.Config, db *mongo.MongoDatabase) *ConfigServer {
	logger.Info("Created Config Server")
	return &ConfigServer{
		logger: logger,
		cfg:    cfg,
		ctx:    context.Background(),
		db:     db,
	}
}

func (s *ConfigServer) AddConfig(ctx context.Context, req *configPb.AddConfigRequest) (*configPb.Config, error) {
	s.logger.Debug("Got config data", zap.Int64("user id", req.UserId))
	doc := &mongo.Config{
		UserId:   req.UserId,
		Name:     req.Name,
		Prompt:   req.Prompt,
		Email:    req.Email,
		Password: req.Password,
		Delay:    req.Delay,
	}

	err := s.db.AddConfig(doc)
	if err != nil {
		s.logger.Error("Failed to insert config", zap.Error(err))
		return nil, err
	}
	s.logger.Debug("Config added", zap.Any("config id", doc.Id))
	result := &configPb.Config{
		Id:       doc.Id.Hex(),
		UserId:   doc.UserId,
		Name:     doc.Name,
		Prompt:   doc.Prompt,
		Email:    doc.Email,
		Password: doc.Password,
		Delay:    doc.Delay,
	}
	return result, nil
}

func (s *ConfigServer) GetConfigs(ctx context.Context, req *configPb.GetConfigsRequest) (*configPb.Configs, error) {
	s.logger.Debug("Got user data", zap.Int64("user id", req.UserId))
	configs, err := s.db.GetConfigWithFilter(bson.M{"userId": req.UserId})
	if err != nil {
		s.logger.Error("Failed to find configs", zap.Error(err))
		return nil, err
	}
	s.logger.Debug("Got user configs", zap.Int64("user id", req.UserId))
	var result []*configPb.Config
	for _, config := range configs {
		config := &configPb.Config{
			Id:       config.Id.Hex(),
			UserId:   config.UserId,
			Name:     config.Name,
			Prompt:   config.Prompt,
			Email:    config.Email,
			Password: config.Password,
			Delay:    config.Delay,
		}

		result = append(result, config)
	}

	return &configPb.Configs{Configs: result}, nil
}

func (s *ConfigServer) DeleteConfig(ctx context.Context, req *configPb.DeleteConfigRequest) (*configPb.Empty, error) {
	s.logger.Debug("Got config data", zap.String("config id", req.Id), zap.Int64("user id", req.UserId))
	err := s.db.DeleteConfig(bson.M{"id": req.Id, "userId": req.UserId})
	if err != nil {
		s.logger.Error("Failed to delete config", zap.Error(err))
		return nil, err
	}
	s.logger.Debug("Config is deleted", zap.Int64("user id", req.UserId))
	return &configPb.Empty{}, nil
}

func (s *ConfigServer) GetConfig(ctx context.Context, req *configPb.GetConfigRequest) (*configPb.Config, error) {
	s.logger.Debug("Got config data", zap.Any("config id", req.ConfigId))
	config, err := s.db.GetConfigByID(req.ConfigId)
	if err != nil {
		s.logger.Error("Failed to find configs", zap.Error(err))
		return nil, err
	}
	s.logger.Debug("Got config", zap.Any("config id", config.Id))
	result := &configPb.Config{
		Id:              config.Id.Hex(),
		UserId:          config.UserId,
		Name:            config.Name,
		Prompt:          config.Prompt,
		Email:           config.Email,
		Password:        config.Password,
		Delay:           config.Delay,
		RefreshTokenDTF: config.RefreshTokenDTF,
		RefreshTokenVC:  config.RefreshTokenVC,
	}

	return result, nil
}
