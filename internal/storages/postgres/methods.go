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

// RegisterUser регистрация нового пользователя
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

//func (p *PSQL) GetUser(ctx context.Context, username string) (storages.User, error) {
//	var user storages.User
//	query := "SELECT id, username, password, email FROM users WHERE username = $1"
//	err := p.pool.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
//	return user, err
//}

// CreateWallet создание нового кошелька, вызывается при создании пользователя
func (p *PSQL) CreateWallet(ctx context.Context, wallet storages.Wallet) error {
	if p.pool == nil {
		return errors.New("database pool is not initialized")
	}
	query := "INSERT INTO wallets (uuid, balanceRUB, balanceUSD, balanceEUR) VALUES ($1, $2, $3, $4)"
	_, err := p.pool.Exec(ctx, query, wallet.UUID, 0, 0, 0)
	return err
}

// GetWalletByUsername получение кошелька по имени, необходимо для получения баланса, депозита, снятия и обмена
func (p *PSQL) GetWalletByUsername(ctx context.Context, username string) (storages.Wallet, error) {
	var wallet storages.Wallet

	query := `SELECT w.id, w.uuid, w.balanceRUB, w.balanceUSD, w.balanceEUR 
			  FROM wallets w
			  JOIN users u ON w.uuid = u.wallet_id
		      WHERE u.username = $1`
	err := p.pool.QueryRow(ctx, query, username).Scan(&wallet.ID, &wallet.UUID, &wallet.Balance.RUB, &wallet.Balance.USD, &wallet.Balance.EUR)
	return wallet, err
}

// GetBalance узнать баланс в рублях, долларах, евро
func (p *PSQL) GetBalance(ctx context.Context, user storages.User) (storages.Wallet, error) {
	var wallet storages.Wallet

	p.logger.Info("Getting balance for user", zap.String("username", user.Username))

	if user.Username == "" {
		p.logger.Error("Username is empty")
		return storages.Wallet{}, fmt.Errorf("username is required")
	}

	query := `
        SELECT wallets.id, wallets.uuid, wallets.balanceRUB, wallets.balanceUSD, wallets.balanceEUR
        FROM wallets
        JOIN users ON wallets.uuid = users.wallet_id
        WHERE users.username = $1
    `

	err := p.pool.QueryRow(ctx, query, user.Username).Scan(
		&wallet.ID,
		&wallet.UUID,
		&wallet.Balance.RUB,
		&wallet.Balance.USD,
		&wallet.Balance.EUR,
	)
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

	p.logger.Info("Successfully got balance", zap.String("username", user.Username), zap.Any("wallet", wallet))

	return wallet, nil
}

// Deposit пополнение счёта
func (p *PSQL) Deposit(ctx context.Context, wallet storages.Wallet, currency string, amount float64) error {
	query := ""
	switch currency {
	case "RUB":
		query = `UPDATE wallets SET balanceRUB = balanceRUB + $1 WHERE uuid = $2`
	case "USD":
		query = `UPDATE wallets SET balanceUSD = balanceUSD + $1 WHERE uuid = $2`
	case "EUR":
		query = `UPDATE wallets SET balanceEUR = balanceEUR + $1 WHERE uuid = $2`
	default:
		return errors.New("unsupported currency")
	}
	_, err := p.pool.Exec(ctx, query, amount, wallet.UUID)
	if err != nil {
		return err
	}
	return err
}

// Withdraw снятие со счёта
func (p *PSQL) Withdraw(ctx context.Context, wallet storages.Wallet, currency string, amount float64) error {
	query := ""
	switch currency {
	case "RUB":
		query = `UPDATE wallets SET balanceRUB = balanceRUB - $1 WHERE uuid = $2`
	case "USD":
		query = `UPDATE wallets SET balanceUSD = balanceUSD - $1 WHERE uuid = $2`
	case "EUR":
		query = `UPDATE wallets SET balanceEUR = balanceEUR - $1 WHERE uuid = $2`
	default:
		return errors.New("unsupported currency")
	}
	_, err := p.pool.Exec(ctx, query, amount, wallet.UUID)
	if err != nil {
		return err
	}
	return err
}

// UpdateWalletBalance обновление счёта, необходимо для обновления после обмена валюты
func (p *PSQL) UpdateWalletBalance(ctx context.Context, wallet storages.Wallet) error {
	query := `
        UPDATE wallets
        SET balanceRUB = $1,
            balanceUSD = $2,
            balanceEUR = $3
        WHERE uuid = $4
    `

	_, err := p.pool.Exec(ctx, query, wallet.Balance.RUB, wallet.Balance.USD, wallet.Balance.EUR, wallet.UUID)
	return err
}
