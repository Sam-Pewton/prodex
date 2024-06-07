package logging

import (
	"errors"
	"log/slog"
	"os"
	"strings"
)

var logger *slog.Logger
var logLevel = new(slog.LevelVar)

func StartLogger(level slog.Level) {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	logLevel.Set(level)
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func GetLogLevel(s string) (slog.Level, error) {
	s = strings.ToLower(s)
	switch s {
	case "info":
		return slog.LevelInfo, nil
	case "debug":
		return slog.LevelDebug, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelDebug, errors.New("unknown log level provided, defaulting to debug")
	}
}
