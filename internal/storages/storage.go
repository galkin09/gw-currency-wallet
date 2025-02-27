package storages

import "context"

type Storage interface {
	RegisterUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, id string) (User, error)
	CreateWallet(ctx context.Context, wallet Wallet) error
	GetWalletByUsername(ctx context.Context, username string) (Wallet, error)
	GetBalance(ctx context.Context, user User) (Wallet, error)
	Deposit(ctx context.Context, wallet Wallet) error
	Withdraw(ctx context.Context, wallet Wallet) error
	Exchange(ctx context.Context, wallet Wallet, exchanger Exchanger) (Wallet, error)
}
