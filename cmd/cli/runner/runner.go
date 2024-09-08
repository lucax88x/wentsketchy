package runner

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lucax88x/wentsketchy/cmd/cli/console"
	"github.com/lucax88x/wentsketchy/internal/clock"
	"github.com/lucax88x/wentsketchy/internal/wentsketchy"

	"github.com/spf13/viper"
)

type RunE func(context.Context, *console.Console, []string, *wentsketchy.Wentsketchy) error

func RunCmdE(
	ctx context.Context,
	logger *slog.Logger,
	viper *viper.Viper,
	console *console.Console,
	args []string,
	run RunE,
) error {
	clock := clock.NewSystemCock()

	di, err := wentsketchy.NewWentsketchy(ctx, logger, clock)

	if err != nil {
		return fmt.Errorf("runner: could not init wentsketchy. %w", err)
	}

	return run(ctx, console, args, di)
}
