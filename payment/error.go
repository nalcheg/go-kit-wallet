package payment

import "errors"

var NeedGreaterBalanceError = errors.New("need greater balance")
var NoCurrenciesExchange = errors.New("we do not exchange unequal currencies")

var UnexpextedInput = errors.New("unexpected input")
