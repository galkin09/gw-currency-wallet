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

// GRPCConfig содержит настройки для подключения к gRPC
type GRPCConfig struct {
	Addr string
}

// LoadConfig загружает конфигурацию из файла .env.
func LoadConfig(envPath string) (*Config, error) {
	const op = "config.LoadConfig"

	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("%s: ошибка загрузки .env файла: %w", op, err)
	}

	connTimeout, err := strconv.Atoi(os.Getenv("DB_CONN_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("%s: неверное значение для DB_CONN_TIMEOUT: %w", op, err)
	}

	postgresConfig := PostgresConfig{
		User:           os.Getenv("DB_USER"),
		Password:       os.Getenv("DB_PASSWORD"),
		Host:           os.Getenv("DB_HOST"),
		Port:           os.Getenv("DB_PORT"),
		Database:       os.Getenv("DB_NAME"),
		ConnTimeout:    connTimeout,
		MigrationsPath: os.Getenv("DB_MIGRATIONS_PATH"),
	}

	grpcConfig := GRPCConfig{
		Addr: os.Getenv("GRPC_ADDR"),
	}

	return &Config{
		Postgres: postgresConfig,
		GRPC:     grpcConfig,
	}, nil
}
