package logging

import (
	"log/slog"
	"os"
)

func New(level string) *slog.Logger {
	var l slog.Level
	if level == "debug" {
		l = slog.LevelDebug
	}
	if level == "error" {
		l = slog.LevelError
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: l}))
}
