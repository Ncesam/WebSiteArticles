package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"

	"backend/pkg/config"
)

type MongoDatabase struct {
	logger *zap.Logger
	cfg    *config.Config
	client *mongo.Client
	db     *mongo.Database
	ctx    context.Context
}

func Connect(logger *zap.Logger, cfg *config.Config) (*MongoDatabase, error) {
	dsn := fmt.Sprintf("mongodb://%s:%s@%s:%d",
		cfg.MONGO.USERNAME,
		cfg.MONGO.PASSWORD,
		cfg.MONGO.HOST,
		cfg.MONGO.PORT)

	clientOptions := options.Client().ApplyURI(dsn).SetConnectTimeout(5 * time.Second).SetServerSelectionTimeout(5 * time.Second)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		logger.Error("Failed to connect to MongoDB", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	db := client.Database("Configs")
	logger.Info("Mongo Connection started")
	return &MongoDatabase{
		logger: logger,
		cfg:    cfg,
		client: client,
		db:     db,
		ctx:    context.Background(),
	}, nil
}

func (m *MongoDatabase) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := m.client.Disconnect(ctx); err != nil {
		m.logger.Error("Failed to disconnect from MongoDB", zap.Error(err))
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	m.logger.Info("MongoDB connection closed successfully")
	return nil
}

func (m *MongoDatabase) AddConfig(config *Config) error {
	m.logger.Debug("Got config data", zap.String("config name", config.Name), zap.Int64("user id", config.UserId))
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	collection := m.db.Collection("configs")
	m.logger.Debug("Got collection configs")
	_, err := collection.InsertOne(ctx, config)
	if err != nil {
		m.logger.Error("Failed to insert config",
			zap.Error(err),
			zap.Any("config", config))
		return err
	}

	m.logger.Debug("Config successfully added", zap.Any("config id", config.Id))
	return nil
}

func (m *MongoDatabase) GetConfigs() ([]Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	collection := m.db.Collection("configs")
	m.logger.Debug("Got collection config")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		m.logger.Error("Failed to find configs", zap.Error(err))
		return nil, fmt.Errorf("failed to find configs: %w", err)
	}
	defer cursor.Close(ctx)

	var configs []Config
	if err = cursor.All(ctx, &configs); err != nil {
		m.logger.Error("Failed to decode configs", zap.Error(err))
		return nil, fmt.Errorf("failed to decode configs: %w", err)
	}
	m.logger.Debug("Got all configs")
	return configs, nil
}
func (m *MongoDatabase) GetConfigByID(idStr string) (*Config, error) {
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	collection := m.db.Collection("configs")
	var config Config
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&config)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	return &config, nil
}

func (m *MongoDatabase) GetConfigWithFilter(filter map[string]interface{}) ([]Config, error) {
	m.logger.Debug("Getting configs with filter", zap.Any("filter", filter))
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	collection := m.db.Collection("configs")
	m.logger.Debug("Got collection configs")
	request := bson.D{}
	for k, v := range filter{
		request = append(request, bson.E{Key: k, Value: v})
	}
	cursor, err := collection.Find(ctx, request)
	if err != nil {
		m.logger.Error("Failed to find configs with filter", zap.Error(err), zap.Any("filter", filter))
		return nil, fmt.Errorf("failed to find configs with filter: %w", err)
	}
	defer cursor.Close(ctx)

	var configs []Config
	err = cursor.All(ctx, &configs)
	if err != nil {
		m.logger.Error("Failed to decode configs", zap.Error(err), zap.Any("configs", configs))
		return nil, fmt.Errorf("failed to decode configs: %w", err)
	}

	m.logger.Debug("Configs fetched with filter", zap.Int("count", len(configs)))
	return configs, nil
}

func (m *MongoDatabase) DeleteConfig(filter interface{}) error {
	m.logger.Debug("Deleting config", zap.Any("filter", filter))
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	collection := m.db.Collection("configs")
	m.logger.Debug("Got collection configs")

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		m.logger.Error("Failed to delete config", zap.Error(err), zap.Any("filter", filter))
		return err
	}

	m.logger.Debug("Config deletion result", zap.Int64("deletedCount", result.DeletedCount))
	return nil
}
