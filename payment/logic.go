package payment

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type service struct {
	repostory Repository
	logger    log.Logger
}

func NewService(rep Repository, logger log.Logger) Service {
	return &service{
		repostory: rep,
		logger:    logger,
	}
}

func (s service) DoPayment(fromAccountID, toAccountID string, amount uint64) error {
	logger := log.With(s.logger, "method", "DoPayment")

	if err := s.repostory.DoPayment(fromAccountID, toAccountID, amount); err != nil {
		level.Error(logger).Log("err", err)
		return err
	}

	logger.Log("payment created", fmt.Sprintf("from %s to %s amount %d", fromAccountID, toAccountID, amount))

	return nil
}

func (s service) GetPayments(req *GetPaymentsRequest) (uint64, []*Payment, error) {
	logger := log.With(s.logger, "method", "GetPayments")

	payments, err := s.repostory.GetPayments(req)
	if err != nil {
		level.Error(logger).Log("err", err)
		return 0, nil, err
	}

	logger.Log("get payments")

	return payments.Total, payments.Payments, nil
}
