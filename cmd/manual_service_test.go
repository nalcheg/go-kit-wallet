// +build manual

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/nalcheg/go-kit-wallet/dbconf"
)

const (
	dbDSN         = "user=postgres password=postgres host=127.0.0.1 port=55432 dbname=wallet sslmode=disable"
	paymentUrl    = "http://127.0.0.1:9030/payment"
	sideA         = "bank_usd"
	sideB         = "alice"
	requestsCount = 1000
)

type result struct {
	sync.Mutex

	TwoHundredth  int
	FiveHundredth int
	Responses     int
}

var accountBalanceQuery = `SELECT balance FROM accounts WHERE id = $1`

func TestConcurrentPayments(t *testing.T) {
	db, err := dbconf.Connect(dbDSN)
	if err != nil {
		t.Fatal(err)
	}

	initialBalances, err := getABSidesBalances(db, accountBalanceQuery, sideA, sideB)
	if err != nil {
		t.Fatal(err)
	}

	res := &result{
		TwoHundredth:  0,
		FiveHundredth: 0,
		Responses:     0,
	}

	for i := 0; i < requestsCount; i++ {
		go doPaymentRequest(res)
	}

	for {
		res.Lock()
		if res.Responses == requestsCount {
			res.Unlock()
			break
		}
		res.Unlock()
		time.Sleep(2 * time.Second)
	}

	res.Lock()
	defer res.Unlock()
	wantedBalanceDiff := uint64(res.TwoHundredth)

	finalBalances, err := getABSidesBalances(db, accountBalanceQuery, sideA, sideB)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("successfull payments - %d", res.TwoHundredth)
	t.Logf("failed payments - %d", res.FiveHundredth)

	for k, b := range finalBalances {
		if k == 0 {
			if initialBalances[k]-wantedBalanceDiff != b {
				t.Error("unexpected balance change for side A")
			}
			t.Logf("sideA balance diff - %d", initialBalances[k]-finalBalances[k])
		} else if k == 1 {
			if initialBalances[k]+wantedBalanceDiff != b {
				t.Error("unexpected balance change for side B")
			}
			t.Logf("sideB balance diff - %d", finalBalances[k]-initialBalances[k])
		}
	}
}

func getABSidesBalances(db *sql.DB, query string, side ...string) ([]uint64, error) {
	var balances []uint64
	for _, s := range side {
		b := uint64(0)
		if err := db.QueryRow(query, s).Scan(&b); err != nil {
			return nil, err
		}
		balances = append(balances, b)
	}

	return balances, nil
}

func doPaymentRequest(results *result) {
	formData := url.Values{
		"account":    {sideA},
		"to_account": {sideB},
		"amount":     {"0.01"},
	}

	statusCode := 500
	resp, err := http.PostForm(paymentUrl, formData)
	if err == nil {
		statusCode = resp.StatusCode
	}

	results.Lock()
	defer results.Unlock()
	results.Responses++
	switch statusCode {
	case 200:
		results.TwoHundredth++
	case 500:
		results.FiveHundredth++
	default:
		panic(errors.New(fmt.Sprintf("unexpected response code %d", statusCode)))
	}
}

func TestConcurrentPaymentsWithChannels(t *testing.T) {
	db, err := dbconf.Connect(dbDSN)
	if err != nil {
		t.Fatal(err)
	}

	initialBalances, err := getABSidesBalances(db, accountBalanceQuery, sideA, sideB)
	if err != nil {
		t.Fatal(err)
	}

	rCh := make(chan int, requestsCount)

	for i := 0; i < requestsCount; i++ {
		go doPaymentRequestChannel(rCh)
	}

	res := result{}
	i := 0
	for respCode := range rCh {
		switch respCode {
		case 200:
			res.TwoHundredth++
		case 500:
			res.FiveHundredth++
		default:
			panic(errors.New(fmt.Sprintf("unexpected response code %d", respCode)))
		}

		i++
		if i == requestsCount {
			close(rCh)
		}
	}

	wantedBalanceDiff := uint64(res.TwoHundredth)

	finalBalances, err := getABSidesBalances(db, accountBalanceQuery, sideA, sideB)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("successfull payments - %d", res.TwoHundredth)
	t.Logf("failed payments - %d", res.FiveHundredth)

	for k, b := range finalBalances {
		if k == 0 {
			if initialBalances[k]-wantedBalanceDiff != b {
				t.Error("unexpected balance change for side A")
			}
			t.Logf("sideA balance diff - %d", initialBalances[k]-finalBalances[k])
		} else if k == 1 {
			if initialBalances[k]+wantedBalanceDiff != b {
				t.Error("unexpected balance change for side B")
			}
			t.Logf("sideB balance diff - %d", finalBalances[k]-initialBalances[k])
		}
	}
}

func doPaymentRequestChannel(ch chan int) {
	formData := url.Values{
		"account":    {sideA},
		"to_account": {sideB},
		"amount":     {"0.01"},
	}

	statusCode := 500
	resp, err := http.PostForm(paymentUrl, formData)
	if err == nil {
		statusCode = resp.StatusCode
	}

	ch <- statusCode
}
