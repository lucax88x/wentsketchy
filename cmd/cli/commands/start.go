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

	"github.com/lucax88x/wentsketchy/cmd/cli/console"
	"github.com/lucax88x/wentsketchy/cmd/cli/runner"
	"github.com/lucax88x/wentsketchy/internal/jobs"
	"github.com/lucax88x/wentsketchy/internal/wentsketchy"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type startOptions struct {
	fifo string
}

func NewStartCmd(
	ctx context.Context,
	logger *slog.Logger,
	viper *viper.Viper,
	console *console.Console,
) *cobra.Command {
	options := &startOptions{}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "start wentsketchy",
		RunE: func(_ *cobra.Command, args []string) error {
			return runner.RunCmdE(ctx, logger, viper, console, args, runStartCmd(options))
		},
	}

	startCmd.SetOut(console.Stdout)
	startCmd.SetErr(console.Stderr)

	configureStartCmdFlags(viper, startCmd, options)

	return startCmd
}

func configureStartCmdFlags(_ *viper.Viper, startCmd *cobra.Command, options *startOptions) {
	startCmd.Flags().StringVar(
		&options.fifo,
		"fifo",
		"/tmp/wentsketchy",
		"Path to fifo file you want to pipe commands to",
	)
}

func runStartCmd(options *startOptions) runner.RunE {
	return func(
		ctx context.Context,
		_ *console.Console,
		_ []string,
		di *wentsketchy.Wentsketchy,
	) error {
		di.Logger.InfoContext(ctx, "start: starting fifo", slog.String("path", options.fifo))

		err := di.Fifo.Start(options.fifo)

		if err != nil {
			return fmt.Errorf("fifo: could not init %w", err)
		}

		di.Logger.InfoContext(ctx, "start: get aerospace tree")

		di.RefreshAerospaceData()

		var wg sync.WaitGroup
		wg.Add(2)
		ch := make(chan error, 2)

		go func() {
			runServer(ctx, di, options.fifo, &wg)

			if err != nil {
				ch <- err

				di.Logger.ErrorContext(ctx, "server: error", slog.Any("error", err))
			}
		}()

		go func() {
			runJobs(ctx, di, &wg)

			if err != nil {
				ch <- err

				di.Logger.ErrorContext(ctx, "jobs: error", slog.Any("error", err))
			}
		}()

		defer close(ch)

		wg.Wait()

		di.Logger.InfoContext(ctx, "start: shutdown complete")

		var aggregateError error
		for message := range ch {
			if message != nil {
				aggregateError = errors.Join(aggregateError, message)
			}
		}

		return aggregateError
	}
}

func runServer(
	ctx context.Context,
	di *wentsketchy.Wentsketchy,
	path string,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	di.Logger.InfoContext(ctx, "server: starting")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	cancelCtx, cancel := context.WithCancel(ctx)

	go func(ctx context.Context) {
		di.Server.Start(ctx, path)
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
		jobs.RefreshAerospaceData(ctx, di, time.Minute)
	}(tickerCtx)

	<-quit

	tickerCancel()

	di.Logger.InfoContext(ctx, "jobs: shutdown")
}
