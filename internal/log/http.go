package log

import (
	"fmt"
	"net/http"
)

type HTTPLoggerImpl struct {
	handlers string
}

func NewHTTPLogger(handlers string) *HTTPLoggerImpl {
	return &HTTPLoggerImpl{handlers: handlers}
}

func (l *HTTPLoggerImpl) extractTags(r *http.Request) []any {
	return []any{
		"handlers", l.handlers,
		"address", r.RemoteAddr,
		"method", r.Method,
		"url", r.URL,
	}
}

func (l *HTTPLoggerImpl) Info(r *http.Request, msg string, args ...any) {
	Info(
		r.Context(),
		fmt.Sprintf(msg, args...),
		l.extractTags(r)...,
	)
}

func (l *HTTPLoggerImpl) Warn(r *http.Request, msg string, args ...any) {
	Warn(
		r.Context(),
		fmt.Sprintf(msg, args...),
		l.extractTags(r)...,
	)
}

func (l *HTTPLoggerImpl) Debug(r *http.Request, msg string, args ...any) {
	Debug(
		r.Context(),
		fmt.Sprintf(msg, args...),
		l.extractTags(r)...,
	)
}

func (l *HTTPLoggerImpl) Error(r *http.Request, msg string, err error, args ...any) {
	Error(
		r.Context(),
		fmt.Sprintf(msg, args...),
		err,
		l.extractTags(r)...,
	)
}
