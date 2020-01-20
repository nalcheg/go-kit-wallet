package account

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateAccount endpoint.Endpoint
	GetAccount    endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateAccount: makeCreateAccountEndpoint(s),
		GetAccount:    makeGetAccountEndpoint(s),
	}
}

func makeCreateAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateAccountRequest)
		err := s.CreateAccount(req.ID, req.Currency, req.Balance)
		return nil, err
	}
}

func makeGetAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAccountRequest)
		acc, err := s.GetAccount(req.ID)
		if err != nil {
			return nil, err
		}

		return GetAccountResponse{
			ID:       acc.ID,
			Balance:  acc.Balance,
			Currency: acc.Currency,
		}, nil
	}
}
