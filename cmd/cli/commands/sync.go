package commands

import (
	"context"
	"errors"
	"log/slog"

	"github.com/lucax88x/wentsketchy/cmd/cli/console"
	"github.com/lucax88x/wentsketchy/cmd/cli/runner"
	"github.com/lucax88x/wentsketchy/internal/wentsketchy"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewSyncCmd(
	ctx context.Context,
	logger *slog.Logger,
	viper *viper.Viper,
	console *console.Console,
) *cobra.Command {
	syncCmd := &cobra.Command{
		Use:   "todo",
		Short: "TODO",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runner.RunCmdE(ctx, logger, viper, console, args, runSyncCmd)
		},
	}

	syncCmd.SetOut(console.Stdout)
	syncCmd.SetErr(console.Stderr)

	return syncCmd
}

func runSyncCmd(ctx context.Context, _ *console.Console, args []string, di *wentsketchy.Wentsketchy) error {
	di.Logger.ErrorContext(ctx, "todo")

	return errors.New("TODO")
}
