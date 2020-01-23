package account

import (
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func AddEndpoints(r *mux.Router, endpoints Endpoints) *mux.Router {
	r.Methods("POST").Path("/account").Handler(httptransport.NewServer(
		endpoints.CreateAccount,
		DecodeCreateAccountReq,
		EmptyResponse,
	))

	r.Methods("GET").Path("/account/{id}").Handler(httptransport.NewServer(
		endpoints.GetAccount,
		DecodeAccountReq,
		EncodeResponse,
	))

	return r
}
