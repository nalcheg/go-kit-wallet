package payment

import (
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func AddEndpoints(r *mux.Router, endpoints Endpoints) *mux.Router {
	r.Methods("POST").Path("/payment").Handler(httptransport.NewServer(
		endpoints.DoPayment,
		DecodeDoPaymentReq,
		EmptyResponse,
	))

	r.Methods("GET").Path("/payment").Handler(httptransport.NewServer(
		endpoints.GetPayments,
		DecodeGetPaymentsReq,
		EncodeResponse,
	))

	return r
}
