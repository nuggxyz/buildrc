package logging

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

func NewJsonLogger() *zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	return &logger
}

func NewJsonStdErrLogger() *zerolog.Logger {
	logger := zerolog.New(os.Stderr).With().Timestamp().Caller().Logger()
	return &logger
}

func NewVerboseLoggerContext(context context.Context) context.Context {
	verboseLogger := NewVerboseLogger()
	return verboseLogger.WithContext(context)
}

func NewVerboseLoggerContextWithLevel(context context.Context, level zerolog.Level) context.Context {
	verboseLogger := NewVerboseLogger()
	lev := verboseLogger.Level(level)
	verboseLogger = &lev
	return verboseLogger.WithContext(context)
}

func NewJsonLoggerContext(context context.Context) context.Context {
	jsonLogger := NewJsonLogger()
	return jsonLogger.WithContext(context)
}

func NewJsonStdErrLoggerContext(context context.Context) context.Context {
	jsonLogger := NewJsonStdErrLogger()
	return jsonLogger.WithContext(context)
}
