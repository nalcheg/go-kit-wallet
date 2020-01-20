package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"

	"github.com/nalcheg/go-kit-wallet/account"
	"github.com/nalcheg/go-kit-wallet/dbconf"
	"github.com/nalcheg/go-kit-wallet/payment"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "wallet",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "wallet service started")
	defer level.Info(logger).Log("msg", "wallet service ended")

	dbDSN := os.Getenv("DB_DSN")
	if len(dbDSN) == 0 {
		level.Error(logger).Log("exit", errors.New("db connection config environment variable missed"))
		os.Exit(-1)
	}

	dbc, err := dbconf.Connect(dbDSN)
	if err != nil {
		level.Error(logger).Log("exit", err)
		os.Exit(-1)
	}
	if err := dbconf.Migrate(dbc); err != nil {
		level.Error(logger).Log("exit", err)
		os.Exit(-1)
	}

	var accountSrv account.Service
	{
		repository := account.NewRepo(dbc, logger)
		accountSrv = account.NewService(repository, logger)
	}
	accountEndpoints := account.MakeEndpoints(accountSrv)

	var paymentSrv payment.Service
	{
		repository := payment.NewRepo(dbc, logger)
		paymentSrv = payment.NewService(repository, logger)
	}
	paymentEndpoints := payment.MakeEndpoints(paymentSrv)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		level.Error(logger).Log("exit", errors.New("listen address config environment variable missed"))
		os.Exit(-1)
	}

	go func() {
		level.Info(logger).Log("msg", "wallet service listen on "+listenAddr)
		r := mux.NewRouter()
		r.Use(commonMiddleware)

		r = account.AddEndpoints(r, accountEndpoints)
		r = payment.AddEndpoints(r, paymentEndpoints)

		errs <- http.ListenAndServe(listenAddr, r)
	}()

	level.Error(logger).Log("exit", <-errs)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
