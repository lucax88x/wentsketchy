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

func NewInitCmd(
	ctx context.Context,
	logger *slog.Logger,
	viper *viper.Viper,
	console *console.Console,
) *cobra.Command {
	syncCmd := &cobra.Command{
		Use:   "init",
		Short: "init sketchybar",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runner.RunCmdE(ctx, logger, viper, console, args, runInitCmd)
		},
	}

	syncCmd.SetOut(console.Stdout)
	syncCmd.SetErr(console.Stderr)

	return syncCmd
}

func runInitCmd(ctx context.Context, _ *console.Console, _ []string, di *wentsketchy.Wentsketchy) error {
	err := config.Init(ctx, di)

	if err != nil {
		return fmt.Errorf("something went wrong initializing sketchybar %w", err)
	}

	return nil
}
