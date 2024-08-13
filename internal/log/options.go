package log

import (
	"time"

	"github.com/rs/zerolog"
)

const (
	DefaultCallerSkipFrameCount = 3
	DefaultTimeFieldFormat      = time.RFC3339Nano
	MessageFieldName            = "msg"
	TimestampFieldName          = "timestamp"
	StacktraceFieldName         = "stacktrace"
	VersionFieldName            = "version"
)

var DefaultOptions = []Option{
	WithDefaults(),
	WithDefaultTimeFormat(),
}

type Option func()

func WithDefaults() Option {
	return func() {
		zerolog.TimestampFieldName = TimestampFieldName
		zerolog.MessageFieldName = MessageFieldName
		zerolog.ErrorStackFieldName = StacktraceFieldName

		withStack = false
		withCaller = false
		withConsoleWriter = false
	}
}

func WithStack() Option {
	return func() {
		withStack = true
	}
}

func WithCaller() Option {
	return func() {
		withCaller = true
	}
}

func WithConsoleWriter() Option {
	return func() {
		withConsoleWriter = true
		bootstrapLogger = zerolog.New(zerolog.NewConsoleWriter(consoleWriterSetup(globalOutput)))
	}
}

func WithDefaultTimeFormat() Option {
	return func() {
		zerolog.TimeFieldFormat = DefaultTimeFieldFormat
	}
}

func applyOptions(opts ...Option) {
	for _, opt := range DefaultOptions {
		opt()
	}

	for _, opt := range opts {
		opt()
	}
}
