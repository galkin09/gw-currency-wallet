package storages

type User struct {
	ID       int
	Username string
	Password string
	Email    string
	WalletID int
}

type Currency struct {
	RUB float64
	USD float64
	EUR float64
}
type Wallet struct {
	ID      int
	UUID    string
	Balance Currency
}

type Deposit struct {
	Amount   float64
	Currency string
}

type Withdraw struct {
	Amount   float64
	Currency string
}

type Exchanger struct {
	FromCurrency string
	ToCurrency   string
	Amount       float64
	Rate         Currency
}

type Rates Currency
