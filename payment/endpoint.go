package payment

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	DoPayment   endpoint.Endpoint
	GetPayments endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		DoPayment:   makeDoPaymentEndpoint(s),
		GetPayments: makeGetPaymentsEndpoint(s),
	}
}

func makeDoPaymentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DoPaymentRequest)
		err := s.DoPayment(req.Account, req.ToAccount, req.Amount)
		return nil, err
	}
}

func makeGetPaymentsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetPaymentsRequest)
		total, payments, err := s.GetPayments(&req)
		if err != nil {
			return nil, err
		}

		return GetPaymentsResponse{
			Total: total,
			Data:  payments,
		}, nil
	}
}
