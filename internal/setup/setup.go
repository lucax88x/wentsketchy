package setup

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/lucax88x/wentsketchy/cmd/cli/console"
	"github.com/spf13/viper"
)

type ExecutionResult = int

const (
	Ok    ExecutionResult = 0
	NotOk ExecutionResult = -1
)

func initViper() (*viper.Viper, error) {
	viperInstance := viper.New()

	return viperInstance, nil
}

type ProgramExecutor func(ctx context.Context, logger *slog.Logger) error

type ExecutorBuilder func(
	viper *viper.Viper,
	console *console.Console,
) ProgramExecutor

func Run(buildExecutor ExecutorBuilder) ExecutionResult {
	start := time.Now()

	logger := slog.New(tint.NewHandler(
		os.Stderr,
		&tint.Options{Level: slog.LevelDebug},
	))

	defer func() {
		elapsed := time.Since(start)
		logger.Info("cli: took", slog.Duration("elapsed", elapsed))
	}()

	viper, err := initViper()

	if err != nil {
		logger.Error("main: could not setup configuration", slog.Any("err", err))
		return NotOk
	}

	console := &console.Console{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	ctx := context.Background()
	err = buildExecutor(viper, console)(ctx, logger)

	if err != nil {
		logger.Error("main: failed to execute program", slog.Any("err", err))
		return NotOk
	}

	logger.Debug("main: completed", slog.Int("status_code", Ok))

	return Ok
}
