package app

import (
	"context"
	pb "github.com/galkin09/proto-exchange/exchange"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"gw-currency-wallet/internal/auth"
	"gw-currency-wallet/internal/handlers"

	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "gw-currency-wallet/docs"
)

func NewRoutes(logger *zap.Logger, dbURL string, migrationsPath string, exchangeClient pb.ExchangeServiceClient, cache *cache.Cache) (*gin.Engine, error) {
	ctx := context.Background()

	h, err := handlers.NewHandler(ctx, logger, dbURL, migrationsPath, exchangeClient, cache)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	public := r.Group("/api/v1")
	{
		public.POST("/register", h.RegisterUser)
		public.POST("/login", h.LoginUser)
	}

	protected := r.Group("/api/v1").Use(auth.Auth())
	{
		protected.GET("/balance", h.GetBalance)
		protected.POST("/wallet/deposit", h.Deposit)
		protected.POST("/wallet/withdraw", h.Withdraw)
		protected.GET("/exchange/rates", h.ExchangeRates)
		protected.POST("/exchange", h.Exchange)
	}

	return r, nil
}
