package exceptions

import "github.com/pkg/errors"

var ErrBalanceNotFound = errors.New("balance hasn't been found")
var ErrBalanceIsNegative = errors.New("balance is negative")
