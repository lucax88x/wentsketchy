package testutils

import (
	"log/slog"
	"os"
)

func CreateTestLogger() *slog.Logger {
	// slog.Default()
	return slog.New(slog.NewTextHandler(os.Stderr, nil))
}
