package logger

import (
	"gicicm/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var (
	logger *zap.Logger
	err    error
	once   sync.Once
)

// New returns a new instance of the logger.
// XXX: Only on the first call a new instance is received
// any subsequent calls return previously initialized instance.
func New(config config.Config) (*zap.Logger, error) {
	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(getZapLogLevel(config.LogLevel)),
		Encoding:         "json",
		OutputPaths:      []string{"stderr", "/tmp/log"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			TimeKey:     "timestamp",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
		},
		DisableStacktrace: true,
		DisableCaller:     false,
	}

	// ensures that there is only one instance of the logger at anytime.
	once.Do(func() {
		logger, err = zapConfig.Build()
	})

	return logger, err
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
// based on the logLevel input string.
// defaults to zapcore.InfoLevel
func getZapLogLevel(logLevel string) zapcore.Level {

	var level zapcore.Level

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
