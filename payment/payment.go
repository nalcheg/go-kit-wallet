package payment

type Payment struct {
	ID          string  `json:"id"`
	Amount      float64 `json:"amount"`
	Account     string  `json:"account"`
	FromAccount *string `json:"from_account,omitempty"`
	ToAccount   *string `json:"to_account,omitempty"`
	Direction   string  `json:"direction"`
}

type Repository interface {
	DoPayment(fromAccountID, toAccountID string, amount uint64) error
	GetPayments(req *GetPaymentsRequest) (*PaymentsResponse, error)
}
