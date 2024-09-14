package config

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
)

func (cfg *Config) Update(
	ctx context.Context,
	args *args.Args,
) error {
	cfg.logger.InfoContext(
		ctx,
		"reacting to event",
		slog.Any("name", args.Name),
		slog.Any("event", args.Event),
	)

	var batches = make([][]string, 0)

	batches, err := cfg.items.Calendar.Update(batches, args)

	if err != nil {
		return fmt.Errorf("update calendar %w", err)
	}

	batches, err = cfg.items.FrontApp.Update(batches, args)

	if err != nil {
		return fmt.Errorf("update calendar %w", err)
	}

	err = cfg.sketchybar.Run(ctx, flatten(batches...))

	if err != nil {
		return fmt.Errorf("apply to sketchybar %w", err)
	}

	return nil
}
