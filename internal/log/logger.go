package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

type options struct {
	defaultLogLevel slog.Level
	writer          io.Writer
}

type option func(*options)

func WithWriter(w io.Writer) option {
	return func(o *options) {
		o.writer = w
	}
}

func NewLogger(loggerName string, opts ...option) *slog.Logger {
	options := options{
		defaultLogLevel: slog.LevelInfo,
		writer:          os.Stdout,
	}
	for _, opt := range opts {
		opt(&options)
	}

	lvl := &slog.LevelVar{}
	lvl.Set(options.defaultLogLevel) // Default

	if env := os.Getenv("GO_LOG_LEVEL"); env != "" {
		if err := lvl.UnmarshalText([]byte(env)); err != nil {
			msg := fmt.Sprintf("invalid log level in GO_LOG_LEVEL, defaulting to %q", options.defaultLogLevel)
			slog.Warn(msg, slog.Any("error", err), slog.String("value", env))
		}
	}

	logger := slog.New(slog.NewJSONHandler(options.writer, &slog.HandlerOptions{
		Level:     lvl,
		AddSource: true,
	}))
	logger = logger.With(slog.String("app", "card-game"), slog.String("logger_name", loggerName))

	return logger
}
