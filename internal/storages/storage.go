package storages

import "context"

type Storage interface {
	//User methods
	RegisterUser(ctx context.Context, user User) error
	//GetUser(ctx context.Context, id string) (User, error)

	//Wallet methods
	CreateWallet(ctx context.Context, wallet Wallet) error
	GetWalletByUsername(ctx context.Context, username string) (Wallet, error)
	GetBalance(ctx context.Context, user User) (Wallet, error)

	//Deposit/Withdraw methods
	Deposit(ctx context.Context, wallet Wallet, currency string, amount float64) error
	Withdraw(ctx context.Context, wallet Wallet, currency string, amount float64) error

	//Exchange method
	//Exchange(ctx context.Context, wallet Wallet, exchanger Exchanger) (Wallet, error)
	UpdateWalletBalance(ctx context.Context, wallet Wallet) error
}
