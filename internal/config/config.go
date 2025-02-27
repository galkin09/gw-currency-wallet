package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	Postgres PostgresConfig
	GRPC     GRPCConfig
}

// PostgresConfig содержит настройки для подключения к PostgreSQL.
type PostgresConfig struct {
	User           string
	Password       string
	Host           string
	Port           string
	Database       string
	ConnTimeout    int
	MigrationsPath string
}

// ConnectionURL генерирует строку подключения к PostgreSQL
func (c *PostgresConfig) ConnectionURL() (string, error) {
	if c.Host == "" || c.Port == "" || c.Database == "" || c.User == "" || c.Password == "" {
		return "", fmt.Errorf("некоторые параметры подключения отсутствуют")
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Database), nil
}

type GRPCConfig struct {
	Addr string
}

// LoadConfig загружает конфигурацию из файла .env.
func LoadConfig(envPath string) (*Config, error) {
	const op = "config.LoadConfig"

	// Загружаем переменные окружения из файла .env
	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("%s: ошибка загрузки .env файла: %w", op, err)
	}

	// Парсим настройки PostgreSQL
	connTimeout, err := strconv.Atoi(os.Getenv("POSTGRES_CONN_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("%s: неверное значение для POSTGRES_CONN_TIMEOUT: %w", op, err)
	}

	postgresConfig := PostgresConfig{
		User:           os.Getenv("POSTGRES_USER"),
		Password:       os.Getenv("POSTGRES_PASSWORD"),
		Host:           os.Getenv("POSTGRES_HOST"),
		Port:           os.Getenv("POSTGRES_PORT"),
		Database:       os.Getenv("POSTGRES_DBNAME"),
		ConnTimeout:    connTimeout,
		MigrationsPath: os.Getenv("POSTGRES_MIGRATIONS_PATH"),
	}

	// Парсим настройки gRPC
	grpcConfig := GRPCConfig{
		Addr: os.Getenv("GRPC_ADDR"),
	}

	return &Config{
		Postgres: postgresConfig,
		GRPC:     grpcConfig,
	}, nil
}
