How to run concurrency payments test.
--

Running this test requires docker and docker-composer. Go version
should be 1.13.

All commands should be executed in the project directory by
the below order.

- `make run` - creates docker-compose with 5 instances of wallet service
  (after command finish, you need check docker-compose stack state with
  `docker-compose ps` if some service not UP - re-execution of the
  command (`make run`) is required)
- `make insert-accounts` - creates accounts for test

Next commands execute tests from `cmd/manual_service_test.go` (this test
file has custom build tag to avoid execution with unit tests).
- `make test-concurent-payments-1` - executes variant 1 test, this
  variant collects results in
  struct with mutex
- **AND/OR** `make test-concurent-payments-2` - executes variant 2 test,
  this variant collects results via channel messages

All this tests tries execute 1000 requests for transfer funds from
account A to account B and after that they check that this accounts
balances changed correctly.
