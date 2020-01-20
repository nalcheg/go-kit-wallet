package account

import (
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

func (s service) CreateAccount(id string, currency string, balance uint64) error {
	logger := log.With(s.logger, "method", "CreateAccount")

	acc := &Account{
		ID:       id,
		Balance:  balance,
		Currency: currency,
	}

	if err := s.repostory.CreateAccount(acc); err != nil {
		level.Error(logger).Log("err", err)
		return err
	}

	logger.Log("account created", id)

	return nil
}

func (s service) GetAccount(id string) (*Account, error) {
	logger := log.With(s.logger, "method", "GetAccount")

	email, err := s.repostory.GetAccount(id)

	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	logger.Log("get account", id)

	return email, nil
}
