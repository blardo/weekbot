package logger

import (
	"log/slog"
	"os"
	"weekbot-go/internal/services"

	"github.com/joho/godotenv"
)

var defaultLogger *slog.Logger

func init() {
	// Load .env file early for logger initialization
	// This ensures environment variables are available even if config hasn't been loaded yet
	_ = godotenv.Load()
	_ = godotenv.Load(".env")
	
	// Use JSON handler for production, Text handler for development
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	
	// Get config to check log format and environment
	// Note: This will initialize config if not already done
	config := services.GetConfig()
	
	// Check if we're in development mode or text format requested
	useTextFormat := config.LogFormat == "text" || config.Env == "development"
	if useTextFormat {
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
