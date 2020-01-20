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
- `make test-concurent-payments` - executes cmd/manual_service_test.go
  (this test has custom build tag to avoid execution with unit tests)

  this test tries execute 1000 requests for transfer funds from account
  A to account B and after that it checks that their balances changed
  correctly