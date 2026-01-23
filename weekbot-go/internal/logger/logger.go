package logger

import (
	"log/slog"
	"os"
)

var defaultLogger *slog.Logger

func init() {
	// Use JSON handler for production, Text handler for development
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	
	// Check if we're in development mode
	if os.Getenv("LOG_FORMAT") == "text" || os.Getenv("ENV") == "development" {
		defaultLogger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	} else {
		defaultLogger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}
}

// GetLogger returns the default logger instance
func GetLogger() *slog.Logger {
	return defaultLogger
}

// Info logs an info message with optional key-value pairs
func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

// Error logs an error message with optional key-value pairs
func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

// Warn logs a warning message with optional key-value pairs
func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

// Debug logs a debug message with optional key-value pairs
func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

// With returns a logger with the given key-value pairs added to context
func With(args ...any) *slog.Logger {
	return defaultLogger.With(args...)
}
