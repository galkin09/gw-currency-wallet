package main

import (
	"context"
	pb "github.com/galkin09/proto-exchange/exchange"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gw-currency-wallet/internal/app"
	"gw-currency-wallet/internal/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Инициализация логгера
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Контекст для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Загрузка конфигурации
	cfg, err := config.LoadConfig("config.env")
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}
	logger.Info("Конфигурация загружена", zap.Any("cfg", cfg))

	// Формирование URL для подключения к базе данных
	dbURL, err := cfg.Postgres.ConnectionURL()
	if err != nil {
		logger.Fatal("Failed to generate DB connection URL", zap.Error(err))
	}

	// Инициализация gRPC-клиента
	conn, err := grpc.Dial(cfg.GRPC.Addr, grpc.WithInsecure())
	if err != nil {
		logger.Fatal("Failed to connect to gRPC server", zap.Error(err))
	}
	defer conn.Close()

	exchangeClient := pb.NewExchangeServiceClient(conn)

	// Инициализация кэша
	c := cache.New(5*time.Minute, 10*time.Minute) // TTL 5 минут, интервал очистки 10 минут

	// Настройка роутинга
	router, err := app.NewRoutes(logger, dbURL, cfg.Postgres.MigrationsPath, exchangeClient, c)
	if err != nil {
		logger.Fatal("Failed to create routes", zap.Error(err))
	}

	// Запуск сервера
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Ожидание сигналов для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Остановка сервера с таймаутом
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
