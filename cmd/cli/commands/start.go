package commands

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/cmd/cli/console"
	"github.com/lucax88x/wentsketchy/cmd/cli/runner"
	"github.com/lucax88x/wentsketchy/internal/jobs"
	"github.com/lucax88x/wentsketchy/internal/wentsketchy"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewStartCmd(
	ctx context.Context,
	logger *slog.Logger,
	viper *viper.Viper,
	console *console.Console,
) *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "start wentsketchy",
		RunE: func(_ *cobra.Command, args []string) error {
			return runner.RunCmdE(ctx, logger, viper, console, args, runStartCmd())
		},
	}

	startCmd.SetOut(console.Stdout)
	startCmd.SetErr(console.Stderr)

	return startCmd
}

func runStartCmd() runner.RunE {
	return func(
		ctx context.Context,
		_ *console.Console,
		_ []string,
		di *wentsketchy.Wentsketchy,
	) error {
		di.Logger.InfoContext(
			ctx,
			"start: starting fifo",
			slog.String("path", settings.FifoPath),
		)

		err := di.Fifo.Start(settings.FifoPath)

		if err != nil {
			return fmt.Errorf("start: could not start fifo %w", err)
		}

		di.Logger.InfoContext(ctx, "start: refresh aerospace tree")

		di.Aerospace.SingleFlightRefreshTree()

		di.Logger.InfoContext(ctx, "start: config init")

		err = di.Config.Init(ctx)

		if err != nil {
			return fmt.Errorf("start: could not config init %w", err)
		}

		var wg sync.WaitGroup
		wg.Add(2)

		var aggregateError error
		go func() {
			runServer(ctx, di, &wg)

			if err != nil {
				aggregateError = errors.Join(aggregateError, fmt.Errorf("server: error. %w", err))
			}
		}()

		go func() {
			runJobs(ctx, di, &wg)

			if err != nil {
				aggregateError = errors.Join(aggregateError, fmt.Errorf("jobs: error. %w", err))
			}
		}()

		wg.Wait()

		di.Logger.InfoContext(ctx, "start: shutdown complete")

		if aggregateError != nil {
			di.Logger.ErrorContext(ctx, "server: error", slog.Any("error", aggregateError))
		}

		return aggregateError
	}
}

func runServer(
	ctx context.Context,
	di *wentsketchy.Wentsketchy,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	di.Logger.InfoContext(ctx, "server: starting")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	cancelCtx, cancel := context.WithCancel(ctx)

	go func(ctx context.Context) {
		di.Server.Start(ctx)
	}(cancelCtx)

	<-quit

	cancel()

	di.Logger.InfoContext(ctx, "server: shutdown")
}

func runJobs(
	ctx context.Context,
	di *wentsketchy.Wentsketchy,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	di.Logger.InfoContext(ctx, "jobs: starting")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	tickerCtx, tickerCancel := context.WithCancel(ctx)

	go func(ctx context.Context) {
		// di.Config.Cfg
		jobs.RefreshAerospaceData(ctx, di, time.Minute)
	}(tickerCtx)

	<-quit

	tickerCancel()

	di.Logger.InfoContext(ctx, "jobs: shutdown")
}
