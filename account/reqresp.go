package account

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type (
	CreateAccountRequest struct {
		ID       string
		Balance  uint64
		Currency string
	}

	GetAccountRequest struct {
		ID string
	}
	GetAccountResponse struct {
		ID       string `json:"id"`
		Balance  uint64 `json:"balance"`
		Currency string `json:"currency"`
	}
)

var UnexpextedInput = errors.New("unexpected input")

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func EmptyResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return nil
}

func DecodeCreateAccountReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateAccountRequest
	if err := r.ParseForm(); err != nil {
		return nil, UnexpextedInput
	}

	if len(r.FormValue("account")) == 0 ||
		len(r.FormValue("balance")) == 0 ||
		len(r.FormValue("currency")) == 0 {
		return nil, UnexpextedInput
	}

	balanceFloat, err := strconv.ParseFloat(r.FormValue("balance"), 64)
	if err != nil {
		return nil, UnexpextedInput
	}
	balance := uint64(balanceFloat * 100)

	req = CreateAccountRequest{
		ID:       r.FormValue("account"),
		Balance:  balance,
		Currency: r.FormValue("currency"),
	}

	return req, nil
}

func DecodeAccountReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req GetAccountRequest
	vars := mux.Vars(r)

	req = GetAccountRequest{
		ID: vars["id"],
	}

	return req, nil
}
