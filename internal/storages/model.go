package storages

type User struct {
	ID       int    `json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	//WalletID int    //TODO: нужно ли это?
}

type Currency struct {
	RUB float64 `json:"rub"`
	USD float64 `json:"usd"`
	EUR float64 `json:"eur"`
}
type Wallet struct {
	ID      int      `json:"id"`
	UUID    string   `json:"uuid"`
	Balance Currency `json:"balance"`
}

type Deposit struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type Withdraw struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type Exchanger struct {
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Amount       float64 `json:"amount"`
}

type Rates Currency
