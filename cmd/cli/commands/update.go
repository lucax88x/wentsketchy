package commands

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lucax88x/wentsketchy/cmd/cli/config"
	"github.com/lucax88x/wentsketchy/cmd/cli/console"
	"github.com/lucax88x/wentsketchy/cmd/cli/runner"
	"github.com/lucax88x/wentsketchy/internal/wentsketchy"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewUpdateCmd(
	ctx context.Context,
	logger *slog.Logger,
	viper *viper.Viper,
	console *console.Console,
) *cobra.Command {
	syncCmd := &cobra.Command{
		Use:   "update",
		Short: "update triggers from sketchybar, to not be used directly",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runner.RunCmdE(ctx, logger, viper, console, args, runUpdateCmd)
		},
	}

	syncCmd.SetOut(console.Stdout)
	syncCmd.SetErr(console.Stderr)

	return syncCmd
}

func runUpdateCmd(ctx context.Context, _ *console.Console, args []string, di *wentsketchy.Wentsketchy) error {
	di.Logger.InfoContext(ctx, "got args for update, tbd", slog.Any("args", args))

	err := config.Update(ctx, di)

	if err != nil {
		return fmt.Errorf("something went wrong updating sketchybar %w", err)
	}

	return nil
}
