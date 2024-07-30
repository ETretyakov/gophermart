package log

import "net/http"

type HTTPLogger interface {
	Info(r *http.Request, msg string, args ...any)
	Warn(r *http.Request, msg string, args ...any)
	Debug(r *http.Request, msg string, args ...any)
	Error(r *http.Request, msg string, err error, args ...any)
}
