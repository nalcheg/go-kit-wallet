// +build manual

package main

import (
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

func TestConcurrentPayments(t *testing.T) {
	db, err := dbconf.Connect(dbDSN)
	if err != nil {
		t.Fatal(err)
	}

	initialBalanceA, initialBalanceB := uint64(0), uint64(0)
	selectAccountBalanceQuery := `SELECT balance FROM accounts WHERE id = $1`
	if err := db.QueryRow(selectAccountBalanceQuery, "bank_usd").Scan(&initialBalanceA); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(selectAccountBalanceQuery, "alice").Scan(&initialBalanceB); err != nil {
		t.Fatal(err)
	}

	results := &result{
		TwoHundredth:  0,
		FiveHundredth: 0,
		Responses:     0,
	}

	for i := 0; i < requestsCount; i++ {
		go doPaymentRequest(results)
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(2 * time.Second)
	for {
		results.Lock()
		if results.Responses == requestsCount {
			results.Unlock()
			break
		}
		results.Unlock()
		time.Sleep(2 * time.Second)
	}

	results.Lock()
	defer results.Unlock()
	wantedBalanceDiff := uint64(results.TwoHundredth)

	finalBalanceA, finalBalanceB := uint64(0), uint64(0)
	if err := db.QueryRow(selectAccountBalanceQuery, sideA).Scan(&finalBalanceA); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(selectAccountBalanceQuery, sideB).Scan(&finalBalanceB); err != nil {
		t.Fatal(err)
	}

	if initialBalanceA-wantedBalanceDiff != finalBalanceA {
		t.Error("unexpected balance change for side A")
	}
	if initialBalanceB+wantedBalanceDiff != finalBalanceB {
		t.Error("unexpected balance change for side B")
	}

	t.Logf("successfull payments - %d", results.TwoHundredth)
	t.Logf("failed payments - %d", results.FiveHundredth)
	t.Logf("sideA balance diff - %d", initialBalanceA-finalBalanceA)
	t.Logf("sideB balance diff - %d", finalBalanceB-initialBalanceB)
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
