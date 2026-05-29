package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

var (
	FallbackLogger *Logger
)

func init() {
	FallbackLogger = NewLogger("fallback-logger")
}

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

func NewLogger(loggerName string, opts ...option) *Logger {
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

	return newLogger(logger)
}

type Logger struct {
	Logger *slog.Logger
}

func newLogger(logger *slog.Logger) *Logger {
	var l Logger
	l.Logger = logger

	return &l
}

func (l *Logger) With(args ...any) *Logger {
	l.Logger = l.Logger.With(args...)

	return l
}

func (l *Logger) Debugf(msg string, args ...any) {
	l.Logger.Debug(fmt.Sprintf(msg, args...))
}

func (l *Logger) Infof(msg string, args ...any) {
	l.Logger.Info(fmt.Sprintf(msg, args...))
}

func (l *Logger) Warnf(msg string, args ...any) {
	l.Logger.Warn(fmt.Sprintf(msg, args...))
}

func (l *Logger) Errorf(msg string, args ...any) {
	l.Logger.Error(fmt.Sprintf(msg, args...))
}
