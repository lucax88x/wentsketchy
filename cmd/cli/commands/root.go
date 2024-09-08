package commands

import (
	"context"
	"log/slog"

	"github.com/lucax88x/wentsketchy/cmd/cli/console"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd(
	ctx context.Context,
	logger *slog.Logger,
	viper *viper.Viper,
	console *console.Console,
) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "cli",
		SilenceUsage: true,
	}

	rootCmd.SetOut(console.Stdout)
	rootCmd.SetErr(console.Stderr)

	configureRootCmdFlags(viper, rootCmd)

	rootCmd.AddCommand(NewSyncCmd(ctx, logger, viper, console))

	return rootCmd
}

func configureRootCmdFlags(viper *viper.Viper, rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Display more verbose output in console output. (default: false)")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Display debugging output in the console. (default: false)")

	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}
