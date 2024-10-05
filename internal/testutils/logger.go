package testutils

import (
	"log/slog"
	"os"
)

func CreateTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, nil))
}
