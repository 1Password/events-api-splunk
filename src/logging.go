package main

import (
	"log/slog"
	"os"
	"strings"

	slogmulti "github.com/samber/slog-multi"
)

func setupLogging() {
	// Set up slog logger with INFO level by default, DEBUG if LOG_LEVEL=debug
	jsonLevel := slog.LevelInfo
	if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" {
		jsonLevel = slog.LevelDebug
	}
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: jsonLevel,
	})
	errorHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		// only log errors to stderr
		Level: slog.LevelError,
	})

	logger := slog.New(slogmulti.Fanout(jsonHandler, errorHandler))
	slog.SetDefault(logger)
}
