package payment

type Service interface {
	DoPayment(fromAccountID, toAccountID string, amount uint64) error
	GetPayments(req *GetPaymentsRequest) (uint64, []*Payment, error)
}
