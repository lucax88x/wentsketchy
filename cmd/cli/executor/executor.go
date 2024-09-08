package executor

import (
	"context"
	"log/slog"

	"github.com/lucax88x/wentsketchy/cmd/cli/commands"
	"github.com/lucax88x/wentsketchy/cmd/cli/console"
	"github.com/lucax88x/wentsketchy/internal/setup"
	"github.com/spf13/viper"
)

func NewCliExecutor(viper *viper.Viper, console *console.Console) setup.ProgramExecutor {
	return func(ctx context.Context, logger *slog.Logger) error {
		rootCmd := commands.NewRootCmd(ctx, logger, viper, console)

		err := rootCmd.Execute()

		if err != nil {
			logger.ErrorContext(ctx, "cli: failed to execute command", slog.Any("error", err))
		}

		return err
	}
}
