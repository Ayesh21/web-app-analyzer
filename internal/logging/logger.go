package logging

import (
	"log/slog"
	"os"
)

// Logger instance
var Logger *slog.Logger

// InitLogger initializes the logger
func InitLogger() {
	if Logger != nil {
		return // Avoid re-initialization
	}

	// Create logs directory if not exists
	os.Mkdir("logs", 0755)

	// Clear previous logs
	logFilePath := "logs/app.log"
	os.WriteFile(logFilePath, []byte{}, 0644)

	// Open log file
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
