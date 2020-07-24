package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	err    error
	once   sync.Once
)

const (
	logLevelEnvVar = "LOG_LEVEL"
)

// New returns a new instance of the logger.
// XXX: Only on the first call a new instance is received
// any subsequent calls return previously initialized instance.
func Log() *zap.Logger {
	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(getZapLogLevel()),
		Encoding:         "json",
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			TimeKey:     "timestamp",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
		},
	}

	// ensures that there is only one instance of the logger at anytime.
	once.Do(func() {
		logger, err = zapConfig.Build()
		if err != nil {
			panic("Logger could not be initialized. Error: " + err.Error())
		}
	})

	return logger
}

/*
// DebugLevel logs are typically voluminous, and are usually disabled in
// production.

// InfoLevel is the default logging priority.

// WarnLevel logs are more important than Info, but don't need individual
// human review.

// ErrorLevel logs are high-priority. If an application is running smoothly,
// it shouldn't generate any error-level logs.

// DPanicLevel logs are particularly important errors. In development the
// logger panics after writing the message.

// PanicLevel logs a message, then panics.

// FatalLevel logs a message, then calls os.Exit(1).

*/

// getZapLogLevel returns a corresponding zaplog level
// based on the logLevel from the value set in the envvar LOG_LEVEL.
// defaults to zapcore.InfoLevel
func getZapLogLevel() zapcore.Level {

	var level zapcore.Level
	logLevel := os.Getenv(logLevelEnvVar)

	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}

	return level
}
