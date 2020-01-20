package account

type Service interface {
	CreateAccount(id string, currency string, balance uint64) error
	GetAccount(id string) (*Account, error)
}
