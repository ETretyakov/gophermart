package exceptions

import "github.com/pkg/errors"

var ErrOrderAlreadyRegistered = errors.New("order already registered")
var ErrOrderAlreadyAccepted = errors.New("order already registered by user")
var ErrOrderNotFound = errors.New("order doesn't exist")
var ErrWrongOrderNumber = errors.New("wrong order number")
