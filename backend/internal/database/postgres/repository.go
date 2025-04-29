package postgres

import (
	"backend/pkg/config"
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDatabase struct {
	logger *zap.Logger
	cfg    *config.Config
	db     *gorm.DB
}

func Connect(cfg *config.Config, logger *zap.Logger) *PostgresDatabase {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.POSTGRES.HOST, cfg.POSTGRES.USERNAME, cfg.POSTGRES.PASSWORD, cfg.POSTGRES.DATABASE, cfg.POSTGRES.PORT)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Postgres not connected", zap.Error(err))
		return nil
	}
	sqlDb, err := db.DB()
	if err != nil {
		logger.Error("Postgres Conn pool not initialize")
		return nil
	}

	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(time.Hour)
	logger.Info("Postgres Started")
	return &PostgresDatabase{
		logger: logger,
		cfg:    cfg,
		db:     db,
	}
}
func (db *PostgresDatabase) Migrate() error {
	err := db.db.AutoMigrate(&User{})
	if err != nil {
		db.logger.Error("Failed to migrate database", zap.Error(err))
		return err
	}
	db.logger.Info("Database migration completed successfully")
	return nil
}

func (db *PostgresDatabase) AddUser(user *User) error {
	result := db.db.Create(user)
	if result.Error != nil {
		db.logger.Error("Failed to create user", zap.Error(result.Error))
		return result.Error
	}
	db.logger.Debug("User created")
	return nil
}
func (db *PostgresDatabase) GetUser(filters interface{}) (*User, error) {
	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Логируем попытку запроса
	db.logger.Debug("Attempting to get user", zap.Any("filters", filters))

	// Выполняем запрос с контекстом
	result := db.db.WithContext(ctx).Where(filters).Take(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			db.logger.Debug("User not found", zap.Any("filters", filters))
			return nil, fmt.Errorf("user not found: %w", result.Error)
		}

		db.logger.Error("Failed to get user",
			zap.Error(result.Error),
			zap.Any("filters", filters),
			zap.String("query", db.db.ToSQL(func(tx *gorm.DB) *gorm.DB {
				return tx.Where(filters).Take(&user)
			})),
		)
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	// Логируем успешное выполнение (с ограничением данных)
	db.logger.Debug("User found",
		zap.Int64("id", user.ID),
		zap.String("email", user.Email),
		zap.Time("created_at", user.CreatedAt),
	)

	return &user, nil
}
