package payment

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type (
	DoPaymentRequest struct {
		Account   string
		ToAccount string
		Amount    uint64
	}

	GetPaymentsRequest struct {
		Page      *uint64
		PerPage   *uint64
		AccountID *string
		Direction *Direction
	}
	GetPaymentsResponse struct {
		Total uint64     `json:"total"`
		Data  []*Payment `json:"data"`
	}
)

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func DecodeDoPaymentReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req DoPaymentRequest

	if len(r.FormValue("account")) == 0 ||
		len(r.FormValue("to_account")) == 0 ||
		len(r.FormValue("amount")) == 0 {
		return nil, UnexpextedInput
	}

	amountFloat, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		return nil, UnexpextedInput
	}
	amount := uint64(amountFloat * 100)

	req = DoPaymentRequest{
		Account:   r.FormValue("account"),
		ToAccount: r.FormValue("to_account"),
		Amount:    amount,
	}

	return req, nil
}

func DecodeGetPaymentsReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req GetPaymentsRequest
	if err := r.ParseForm(); err != nil {
		return nil, UnexpextedInput
	}

	req = GetPaymentsRequest{}
	for key, value := range r.Form {
		switch key {
		case "page":
			page, err := strconv.ParseUint(value[0], 10, 64)
			if err != nil {
				return nil, err
			}
			req.Page = &page
		case "perpage":
			perpage, err := strconv.ParseUint(value[0], 10, 64)
			if err != nil {
				return nil, err
			}
			req.PerPage = &perpage
		case "account_id":
			req.AccountID = &value[0]
		case "direction":
			var d Direction
			if value[0] == "incoming" {
				d = Incoming
			} else if value[0] == "outgoing" {
				d = Outgoing
			} else {
				break
			}
			req.Direction = &d
		}
	}

	return req, nil
}
