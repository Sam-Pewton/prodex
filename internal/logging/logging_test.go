package logging

import (
	"log/slog"
	"testing"
)

func TestSetupLogger(*testing.T) {
	StartLogger(slog.LevelDebug)
}

func TestDebug(*testing.T) {
	logger.Debug("")
}

func TestInfo(*testing.T) {
	logger.Info("")
}

func TestError(*testing.T) {
	logger.Error("")
}

func TestWarn(*testing.T) {
	logger.Warn("")
}

// Test that the debug level can be retrieved
func TestGetDebugLevel(t *testing.T) {
	logLevel, err := GetLogLevel("debug")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if logLevel != slog.LevelDebug {
		t.Fatalf("incorrect log level retrieved")
	}

	logLevel, err = GetLogLevel("DEBUG")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if logLevel != slog.LevelDebug {
		t.Fatalf("incorrect log level retrieved")
	}

	logLevel, err = GetLogLevel("DeBuG")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if logLevel != slog.LevelDebug {
		t.Fatalf("incorrect log level retrieved")
	}
}

// Test that the info level can be retrieved
func TestGetInfoLevel(t *testing.T) {
	logLevel, err := GetLogLevel("info")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if logLevel != slog.LevelInfo {
		t.Fatalf("incorrect log level retrieved")
	}

	logLevel, err = GetLogLevel("INFO")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if logLevel != slog.LevelInfo {
		t.Fatalf("incorrect log level retrieved")
	}

	logLevel, err = GetLogLevel("InFo")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if logLevel != slog.LevelInfo {
		t.Fatalf("incorrect log level retrieved")
	}
}

// Test that the error level can be retrieved
func TestGetErrorLevel(t *testing.T) {
	logLevel, err := GetLogLevel("error")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if logLevel != slog.LevelError {
		t.Fatalf("incorrect log level retrieved")
	}

	logLevel, err = GetLogLevel("ERROR")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if logLevel != slog.LevelError {
		t.Fatalf("incorrect log level retrieved")
	}

	logLevel, err = GetLogLevel("ErRoR")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if logLevel != slog.LevelError {
		t.Fatalf("incorrect log level retrieved")
	}
}

// Test that the warn level can be retrieved
func TestGetWarningLevel(t *testing.T) {
	logLevel, err := GetLogLevel("warn")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if logLevel != slog.LevelWarn {
		t.Fatalf("incorrect log level retrieved")
	}

	logLevel, err = GetLogLevel("WARN")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if logLevel != slog.LevelWarn {
		t.Fatalf("incorrect log level retrieved")
	}

	logLevel, err = GetLogLevel("WaRn")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if logLevel != slog.LevelWarn {
		t.Fatalf("incorrect log level retrieved")
	}
}

// Test that there is an error when using an unknown log level
func TestGetUnknownLevel(t *testing.T) {
	_, err := GetLogLevel("aaa")
	if err == nil {
		t.Fatalf("%s", err)
	}
}
