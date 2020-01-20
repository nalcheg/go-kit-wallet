package account

type Account struct {
	ID       string `json:"id"`
	Balance  uint64 `json:"balance"`
	Currency string `json:"currency"`
}

type Repository interface {
	CreateAccount(account *Account) error
	GetAccount(id string) (*Account, error)
}
