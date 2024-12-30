package srvenv

import (
	"context"
	"todo-api/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewDatabaseConnection(cfg *config.Config, logger *zap.Logger) (*pgxpool.Pool, error) {
	logger.Info("Connecting to postgres database")

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.PostgresDSN) // database connection
	if err != nil {
		logger.Error("Failed to create connection pool", zap.Error(err))
		return nil, err
	}

	// ping database
	if err := pool.Ping(ctx); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		return nil, err
	}

	logger.Info("Successfully connected to databawse")
	return pool, err
}
