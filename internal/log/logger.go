package log

import (
	"context"
	"io"
	"os"

	"github.com/rs/zerolog"
)

var (
	bootstrapLogger zerolog.Logger
	globalLogger    zerolog.Logger
	globalOutput    io.Writer

	withStack         bool
	withCaller        bool
	withConsoleWriter bool
)

func init() {
	bootstrapLogger = zerolog.Nop()
	globalLogger = bootstrapLogger
}

func InitDefault(level string) {
	globalOutput = os.Stderr

	runLogFile, _ := os.OpenFile(
		"logs/stdout.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	multi := zerolog.MultiLevelWriter(globalOutput, runLogFile)
	bootstrapLogger = zerolog.New(multi)

	applyOptions()

	ctx := bootstrapLogger.With().Timestamp().Caller().Stack()

	globalLogger = ctx.Logger()

	GlobalLevelFromString(level)
}

func Init(w io.Writer, level string, opts ...Option) {
	if w == nil {
		w = os.Stderr
	}

	globalOutput = w

	bootstrapLogger = zerolog.New(globalOutput)

	applyOptions(opts...)

	ctx := bootstrapLogger.With().Timestamp()

	if withCaller {
		ctx = ctx.Caller()
	}

	if withStack {
		ctx = ctx.Stack()
	}

	globalLogger = ctx.Logger()

	GlobalLevelFromString(level)
}

func consoleWriterSetup(out io.Writer) func(cw *zerolog.ConsoleWriter) {
	return func(cw *zerolog.ConsoleWriter) {
		cw.Out = out

		cw.PartsOrder = []string{
			TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.CallerFieldName,
			MessageFieldName,
		}
	}
}

func Debug(ctx context.Context, msg string, kv ...any) {
	globalLogger.Debug().Fields(kv).Msg(msg)
}

func Info(ctx context.Context, msg string, kv ...any) {
	globalLogger.Info().Fields(kv).Msg(msg)
}

func Warn(ctx context.Context, msg string, kv ...any) {
	globalLogger.Warn().Fields(kv).Msg(msg)
}

func Error(ctx context.Context, msg string, err error, kv ...any) {
	globalLogger.Error().Fields(kv).Err(err).Msg(msg)
}

func Fatal(ctx context.Context, msg string, err error, kv ...any) {
	globalLogger.Fatal().Fields(kv).Err(err).Msg(msg)
}
