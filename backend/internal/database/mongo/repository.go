package mongo

import (
	"backend/pkg/config"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
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
	return &MongoDatabase{
		logger: logger,
		cfg:    cfg,
		client: client,
		db:     db,
		ctx:    context.Background(),
	}, nil
}

func (m *MongoDatabase) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.client.Disconnect(ctx); err != nil {
		m.logger.Error("Failed to disconnect from MongoDB", zap.Error(err))
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	m.logger.Info("MongoDB connection closed successfully")
	return nil
}

func (m *MongoDatabase) AddConfig(config *Config) (error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := m.db.Collection("configs")

	result, err := collection.InsertOne(ctx, config)
	if err != nil {
		m.logger.Error("Failed to insert config",
			zap.Error(err),
			zap.Any("config", config))
		return err
	}
	if oid, ok := result.InsertedID.(int64); ok {
		config.Id = oid
	}

	m.logger.Debug("Config successfully added",
		zap.Any("config", config))
	return nil
}

func (m *MongoDatabase) GetConfigs() ([]Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.db.Collection("configs")

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

	return configs, nil
}
func (m *MongoDatabase) GetConfigByID(id primitive.ObjectID) (*Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.db.Collection("configs")

	var config Config
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&config)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		m.logger.Error("Failed to get config",
			zap.Error(err),
			zap.String("id", id.Hex()))
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	return &config, nil
}
func (m *MongoDatabase) GetConfigWithFilter(filter interface{}) ([]Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.db.Collection("configs")

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		m.logger.Error("Failed to find configs with filter",
			zap.Error(err),
			zap.Any("filter", filter))
		return nil, fmt.Errorf("failed to find configs with filter: %w", err)
	}
	defer cursor.Close(ctx)

	var configs []Config
	if err = cursor.All(ctx, &configs); err != nil {
		m.logger.Error("Failed to decode configs", zap.Error(err))
		return nil, fmt.Errorf("failed to decode configs: %w", err)
	}

	return configs, nil
}

func (m *MongoDatabase) DeleteConfig (filter interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := m.db.Collection("configs")

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		m.logger.Error("Failed to delete config", zap.Error(err))
		return err
	}
	return nil;
}