package api

import (
	"log"
	"log/slog"
	"os"
)

// Creates default logger and returns default logger if needed
// Returns simple log.Logger for server
func createDefaultLogger(level slog.Leveler) *log.Logger {
	opts := slog.HandlerOptions{
		Level: level,
	}
	jsonHandler := slog.NewJSONHandler(os.Stdout, &opts)
	logger := slog.New(jsonHandler)
	slog.SetDefault(logger)
	logLogger := slog.NewLogLogger(jsonHandler, level.Level())
	return logLogger
}
