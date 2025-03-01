package handlers

import (
	"context"
	pb "github.com/galkin09/proto-exchange/exchange"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"gw-currency-wallet/internal/storages"
	"gw-currency-wallet/internal/storages/postgres"
	"time"
)

type Handler struct {
	storage storages.Storage
	exch    pb.ExchangeServiceClient
	logger  *zap.Logger
	cache   *cache.Cache
}

func NewHandler(ctx context.Context, logger *zap.Logger, dbURL string, migrationsPath string,
	exchangeClient pb.ExchangeServiceClient, c *cache.Cache) (*Handler, error) {
	psql := postgres.NewPSQL(logger)
	if err := psql.Start(ctx, dbURL, 10*time.Second, migrationsPath); err != nil {
		logger.Error("Failed to initialize PostgreSQL", zap.Error(err))
		return nil, err
	}

	return &Handler{
		storage: psql,
		exch:    exchangeClient,
		logger:  logger,
		cache:   c,
	}, nil
}
