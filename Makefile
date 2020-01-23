SHELL:=/bin/bash

run: ## docker-compose up -d --scale wallet=5
	@docker-compose up -d --scale wallet=5
run-wallets: ## docker-compose up -d --scale wallet=5 --build wallet
	@docker-compose rm -sf wallet
	@docker-compose up -d --scale wallet=5 --build wallet
flap-postgres: ## rm and up postgres service
	docker-compose rm -sf postgres
	docker-compose up -d postgres
test-concurent-payments-1: ## run concurent payments test variant 1 (results in struct with mutex)
	@go test -race -tags manual -run ^TestConcurrentPayments$$ -v cmd/manual_service_test.go
test-concurent-payments-2: ## run concurent payments test variant 2 (results received via channel)
	@go test -race -tags manual -run ^TestConcurrentPaymentsWithChannels$$ -v cmd/manual_service_test.go
insert-accounts: ## insert sample accounts
	curl -d "account=bank_usd&currency=USD&balance=10050000000000" http://127.0.0.1:9030/account
	curl -d "account=bank_rub&currency=RUB&balance=10050000000000" http://127.0.0.1:9030/account
	curl -d "account=alice&currency=USD&balance=0" http://127.0.0.1:9030/account
	curl -d "account=bob&currency=RUB&balance=0" http://127.0.0.1:9030/account
thousand-payments: ## 1000 payments from bank_usd to alice
	@number=1 ; while [[ $$number -le 1000 ]] ; do \
		curl -d "account=bank_usd&to_account=alice&amount=0.01" http://127.0.0.1:9030/payment & \
		((number = number + 1)) ; \
	done
golangci-lint: ## Run https://golangci.com/
	@docker run -ti --rm -v `pwd`:/goapp golangci/build-runner golangci-lint -v run

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
