package logging

import (
	"log/slog"
	"os"
	"testing"
)

// Logger instance
var Logger *slog.Logger

// InitLogger initializes the logger
func InitLogger() {
	if Logger != nil {
		return
	}
	// Check if running in test mode
	if isTestMode() {
		Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		}))
		slog.SetDefault(Logger)
		Logger.Warn("Logging disabled in test mode")
		return
	}

	// Check logs directory exists
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0755)
		if err != nil {
			panic("Failed to create logs directory: " + err.Error())
		}
	}

	// Define log file path
	logFilePath := "logs/app.log"

	// Clear previous logs (overwrite with an empty file)
	err := os.WriteFile(logFilePath, []byte{}, 0644)
	if err != nil {
		panic("Failed to clear log file: " + err.Error())
	}

	// Open log file (Create if not exists)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}

	// Setup JSON Logger
	jsonHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	Logger = slog.New(jsonHandler)
	slog.SetDefault(Logger)

	Logger.Info("Logger initialized")
}

// Function to detect if we are running in test mode
func isTestMode() bool {
	return testing.Testing()
}
