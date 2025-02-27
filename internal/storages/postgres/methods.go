package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"gw-currency-wallet/internal/storages"
)

func (p *PSQL) RegisterUser(ctx context.Context, user storages.User) error {
	UUID, _ := uuid.NewUUID()

	wallet := storages.Wallet{UUID: UUID.String()}

	if err := p.CreateWallet(ctx, wallet); err != nil {
		return err
	}

	query := "INSERT INTO users (username, password, email, wallet_id) VALUES ($1, $2, $3, $4)"
	_, err := p.pool.Exec(ctx, query, user.Username, user.Password, user.Email, UUID.String())
	return err
}

func (p *PSQL) GetUser(ctx context.Context, username string) (storages.User, error) {
	var user storages.User
	query := "SELECT id, username, password, email FROM users WHERE username = $1"
	err := p.pool.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	return user, err
}

func (p *PSQL) CreateWallet(ctx context.Context, wallet storages.Wallet) error {
	if p.pool == nil {
		return errors.New("database pool is not initialized")
	}
	query := "INSERT INTO wallets (uuid, balanceRUB, balanceUSD, balanceEUR) VALUES ($1, $2, $3, $4)"
	_, err := p.pool.Exec(ctx, query, wallet.UUID, 0, 0, 0)
	return err
}

func (p *PSQL) GetWalletByUsername(ctx context.Context, username string) (storages.Wallet, error) {
	var wallet storages.Wallet

	query := `SELECT w.id, w.uuid, w.balanceRUB, w.balanceUSD, w.balanceEUR 
			  FROM wallets w
			  JOIN users u ON w.uuid = u.wallet_id
		      WHERE u.username = $1`
	err := p.pool.QueryRow(ctx, query, username).Scan(&wallet.ID, &wallet.UUID, &wallet.Balance.RUB, &wallet.Balance.USD, &wallet.Balance.EUR)
	return wallet, err
}

func (p *PSQL) GetBalance(ctx context.Context, user storages.User) (storages.Wallet, error) {
	var wallet storages.Wallet

	// Логируем начало операции
	p.logger.Info("Getting balance for user", zap.String("username", user.Username))

	// Проверяем входные данные
	if user.Username == "" {
		p.logger.Error("Username is empty")
		return storages.Wallet{}, fmt.Errorf("username is required")
	}

	query := `
        SELECT wallets.balanceRUB, wallets.balanceUSD, wallets.balanceEUR
        FROM wallets
        JOIN users ON wallets.uuid = users.wallet_id
        WHERE users.username = $1
    `

	// Выполняем запрос
	err := p.pool.QueryRow(ctx, query, user.Username).Scan(&wallet.Balance.RUB, &wallet.Balance.USD, &wallet.Balance.EUR)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.logger.Info("User wallet not found", zap.String("username", user.Username))
			return storages.Wallet{}, fmt.Errorf("wallet not found for user: %s", user.Username)
		}
		if errors.Is(err, context.Canceled) {
			p.logger.Warn("Request canceled", zap.String("username", user.Username))
			return storages.Wallet{}, fmt.Errorf("request canceled")
		}
		p.logger.Error("Failed to get balance", zap.Error(err), zap.String("username", user.Username))
		return storages.Wallet{}, fmt.Errorf("failed to get balance: %w", err)
	}

	// Логируем успешное завершение
	p.logger.Info("Successfully got balance", zap.String("username", user.Username))

	return wallet, nil
}

func (p *PSQL) Deposit(ctx context.Context, wallet storages.Wallet) error {
	query := `UPDATE wallets SET balanceRUB = balanceRUB + $1,
                   				 balanceUSD = balanceUSD + $2,
								 balanceEUR = balanceEUR + $3 
               WHERE uuid = $4`
	_, err := p.pool.Exec(ctx, query, wallet.Balance.RUB, wallet.Balance.USD, wallet.Balance.EUR, wallet.UUID)
	return err
}

func (p *PSQL) Withdraw(ctx context.Context, wallet storages.Wallet) error {
	query := `UPDATE wallets SET balanceRUB = balanceRUB - $1,balanceUSD = balanceUSD - $2,
		balanceEUR = balanceEUR - $3 WHERE uuid = $4`
	_, err := p.pool.Exec(ctx, query, wallet.Balance.RUB, wallet.Balance.USD, wallet.Balance.EUR, wallet.UUID)
	return err
}

func (p *PSQL) Exchange(ctx context.Context, wallet storages.Wallet, exchanger storages.Exchanger) (storages.Wallet, error) {
	panic(nil)
}
