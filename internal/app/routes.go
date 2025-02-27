package app

import (
	"context"
	pb "github.com/galkin09/proto-exchange/exchange"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"gw-currency-wallet/internal/auth"
	"gw-currency-wallet/internal/handlers"
)

func NewRoutes(logger *zap.Logger, dbURL string, migrationsPath string, exchangeClient pb.ExchangeServiceClient, cache *cache.Cache) (*gin.Engine, error) {
	ctx := context.Background()

	h, err := handlers.NewHandler(ctx, logger, dbURL, migrationsPath, exchangeClient, cache)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	//TODO: Check default

	r := gin.Default()
	r.POST("/api/v1/register", h.RegisterUser)
	//{
	//	"username": "string",
	//	"password": "string",
	//	"email": "string"
	//}
	r.POST("/api/v1/login", h.LoginUser)
	//{
	//	"username": "string",
	//	"password": "string"
	//}
	rGroup := r.Group("/api/v1").Use(auth.Auth())

	rGroup.GET("/balance", h.GetBalance)
	//{
	//	"balance":
	//	{
	//		"USD": "float",
	//		"RUB": "float",
	//		"EUR": "float"
	//	}
	//}
	rGroup.POST("/wallet/deposit", h.Deposit)
	//{
	//	"amount": 100.00,
	//	"currency": "USD" // (USD, RUB, EUR)
	//}
	rGroup.POST("/wallet/withdraw", h.Withdraw)
	//{
	//	"amount": 50.00,
	//	"currency": "USD" // USD, RUB, EUR)
	//}
	rGroup.GET("/exchange/rates", h.ExchangeRates)
	//{
	//	"rates":
	//	{
	//		"USD": "float",
	//		"RUB": "float",
	//		"EUR": "float"
	//	}
	//}
	//rGroup.POST("/exchange", h.Exchange)
	//{
	//	"message": "Exchange successful",
	//	"exchanged_amount": 85.00,
	//	"new_balance":
	//	{
	//		"USD": 0.00,
	//		"EUR": 85.00
	//	}
	//}

	return r, nil
}
