package exceptions

import "github.com/pkg/errors"

var ErrLoginAlreadyTaken = errors.New("login already taken")
var ErrUserNotFound = errors.New("user hasn't been found")
var ErrNotAuthorised = errors.New("user hasn't been authorised")
