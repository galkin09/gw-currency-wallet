package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"time"
)

// PSQL представляет собой обертку для работы с PostgreSQL.
type PSQL struct {
	pool    *pgxpool.Pool
	timeout time.Duration
	logger  *zap.Logger
}

func NewPSQL(logger *zap.Logger) *PSQL {
	return &PSQL{
		pool:   nil,
		logger: logger.With(zap.String("component", "postgres")),
	}
}

// Start функция для инициализации БД
func (p *PSQL) Start(ctx context.Context, url string, timeout time.Duration, migrationsPath string) error {
	const op = "postgres.Start"

	p.timeout = timeout

	ctxTimeout, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	p.logger.Info("Подключаемся к базе данных", zap.String("url", url))

	pool, err := pgxpool.New(ctxTimeout, url)
	if err != nil {
		p.logger.Error("Ошибка подключения к БД", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := pool.Ping(ctxTimeout); err != nil {
		p.logger.Error("Ошибка проверки подключения к БД", zap.Error(err), zap.String("op", op))
		return fmt.Errorf("%s: %w", op, err)
	}

	p.pool = pool

	if err := doMigrate(url, migrationsPath); err != nil {
		p.logger.Error("Не удалось применить миграции", zap.Error(err), zap.String("op", op))
		return fmt.Errorf("%s: %w", op, err)
	}

	p.logger.Info("Подключение и миграции успешно выполнены", zap.String("url", url))
	return nil
}

// doMigrate выполняет миграции
func doMigrate(dbURL, migrationsPath string) error {
	const op = "postgres.doMigrate"
	if migrationsPath == "" {
		return fmt.Errorf("%s: %s", op, "путь к миграциям пуст")
	}

	m, err := migrate.New("file://"+migrationsPath, dbURL)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop закрывает пул соединений с базой данных.
func (p *PSQL) Stop() {
	const op = "postgres.Stop"

	if p.pool == nil {
		p.logger.Warn("Пул соединений уже закрыт или не был инициализирован", zap.String("op", op))
		return
	}

	p.logger.Info("Закрытие пула соединений с базой данных", zap.String("op", op))
	p.pool.Close()
	p.logger.Info("Пул соединений успешно закрыт", zap.String("op", op))
}
